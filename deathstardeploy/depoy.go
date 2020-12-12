package deathstardeploy

import (
	"fmt"
	"github.com/djmgit/DeathStar/lambdautil"
	vegetaModels "github.com/djmgit/DeathStar/models"
	vegetaUtil "github.com/djmgit/DeathStar/vegeta_core"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

	deathStarDeploy.DeathLogger.Info().Msg("Creating the lambda attack function...")
	lambdaUtil := lambdautil.LambdaUtil {
		AWSRegion: "us-east-1",
		LambdaRole: "arn:aws:iam::253708721073:role/service-role/func-test-1-role-nyalwdp2",
		LambdaFuncName: "func-test-2",
		LambdaFunctionHandler: "main",
		LambdaFunctionRuntime: "go1.x",
		ZipFilePath: deathStarDeploy.ZipFilePath,
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
	err, resultMetrics := vegAttackUtil.VegetaSeqAttack(deathStarDeploy.yamlConfig.Attacks)

	fmt.Println("Attack complete...")
	for _, result := range resultMetrics {
		fmt.Println(result.Latencies.Mean)
		fmt.Println("\n")
		fmt.Println(result.Duration)
		fmt.Println("\n")
		fmt.Println(result.Success)
	}

	deathStarDeploy.DeathLogger.Info().Msg("Cleaning up function...")
	err = lambdaUtil.DeleteFunction()
	if err != nil {
		deathStarDeploy.DeathLogger.Fatal().Err(err).Msg("Faced error while deleting function")
		return err
	}

	fmt.Println("Exiting DeathStar...")

	return nil
}
