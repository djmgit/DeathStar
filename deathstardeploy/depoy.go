package deathstardeploy

import (
	"encoding/json"
	"fmt"
	"github.com/djmgit/DeathStar/lambdautil"
	vegetaModels "github.com/djmgit/DeathStar/models"
	vegetaUtil "github.com/djmgit/DeathStar/vegeta_core"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type DeathStarDeploy struct {
	ZipFilePath string
	ConfPath string
	LocalZip bool
	yamlConfig *vegetaModels.YAMLConfig
	DeathLogger zerolog.Logger
}

// function to read config yaml
func (deathStarDeploy *DeathStarDeploy) readConfYaml() (error) {

	yamlFile, err := ioutil.ReadFile(deathStarDeploy.ConfPath)
	if err != nil {
		deathStarDeploy.DeathLogger.Fatal().Err(err).Msg("Cannot read config")
		return err
	}

	var yamlConfig vegetaModels.YAMLConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		deathStarDeploy.DeathLogger.Fatal().Err(err).Msg("Cannot parse yaml config")
		return err
	}

	deathStarDeploy.yamlConfig = &yamlConfig
	return nil
}

func (deathStarDeploy *DeathStarDeploy) Start() error {

	err := deathStarDeploy.readConfYaml()
	if err != nil {
		return err
	}

	if deathStarDeploy.LocalZip == true {

		// check zip-file-path is present or not
		if deathStarDeploy.ZipFilePath == "" {
			deathStarDeploy.DeathLogger.Fatal().Err(err).Msg("Code zip file not provided")
			return err
		}
	} else {
		// donwload zipfile and set zip file path in zipFilePath
	}

	lambdaConfig := deathStarDeploy.yamlConfig.LambdaConfig

	// set defaults for config values
	LambdaRegion := "us-east-1"
	LambdaMemorySize := int64(128)
	LambdaTimeOut := int64(3)

	// generate the function name
	currentTime := time.Now()
	formattedDateTime := currentTime.Format("01-02-2006T15-04-05")
	LambdaName := "death-star-lambda-" + formattedDateTime

	// get config values from the yaml
	if lambdaConfig.LambdaRegion != "" {
		LambdaRegion = lambdaConfig.LambdaRegion
	}

	if lambdaConfig.LambdaMemorySize != 0 {
		LambdaMemorySize = lambdaConfig.LambdaMemorySize
	}

	if lambdaConfig.LambdaTimeOut != 0 {
		LambdaTimeOut = lambdaConfig.LambdaTimeOut
	}

	deathStarDeploy.DeathLogger.Info().Msg("Creating the lambda attack function...")
	lambdaUtil := lambdautil.LambdaUtil {
		AWSRegion: LambdaRegion,
		LambdaRole: lambdaConfig.LambdaRole,
		LambdaFuncName: LambdaName,
		LambdaFunctionHandler: "main",
		LambdaFunctionRuntime: "go1.x",
		ZipFilePath: deathStarDeploy.ZipFilePath,
		LambdaMemorySize: LambdaMemorySize,
		LambdaTimeOut: LambdaTimeOut,
	}
	err = lambdaUtil.CreateFunction()
	if err != nil {
		deathStarDeploy.DeathLogger.Fatal().Err(err).Msg("Function creation failed")
		return err
	}

	deathStarDeploy.DeathLogger.Info().Msg("Function creation succeeded...")

	// initiate attack and display result
	vegAttackUtil := vegetaUtil.VegetaAttackUtils{
		LmUtil: &lambdaUtil,
	}

	deathStarDeploy.DeathLogger.Info().Msg("Running attack")
	err, lambdaResponses := vegAttackUtil.VegetaSeqAttack(deathStarDeploy.yamlConfig.Attacks)
	deathStarDeploy.DeathLogger.Info().Msg("Attack complete")

	data, _ := json.MarshalIndent(lambdaResponses, "", "	")
	fmt.Println(string(data))

	deathStarDeploy.DeathLogger.Info().Msg("Cleaning up function...")
	err = lambdaUtil.DeleteFunction()
	if err != nil {
		deathStarDeploy.DeathLogger.Fatal().Err(err).Msg("Faced error while deleting function")
		return err
	}

	fmt.Println("Exiting DeathStar...")

	return nil
}
