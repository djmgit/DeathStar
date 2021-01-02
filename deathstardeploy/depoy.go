
// deathstardeploy package contains neccassary functions required
// to orchestrate the attack by reading the conf file.
// This mostly includes:
//	- reading the conf file, figure out the lambda function configs and the alert configs.
//	- Create the lambda function
// 	- Invoke the function to carry out the attack
//	- Clean up the created lambda function
package deathstardeploy

import (
	"encoding/json"
	"fmt"
	"github.com/djmgit/DeathStar/lambdautil"
	vegetaModels "github.com/djmgit/DeathStar/models"
	vegetaUtil "github.com/djmgit/DeathStar/vegeta_core"
	packager "github.com/djmgit/DeathStar/functionpackager"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
	"path/filepath"
)

// DeathStarDeploy holds various deploy related attributes and methods.
type DeathStarDeploy struct {
	ZipFilePath string
	ConfPath string
	LocalZip bool
	yamlConfig *vegetaModels.YAMLConfig
	DeathLogger zerolog.Logger
}

// readConfYaml reads the provided yaml config file.
func (deathStarDeploy *DeathStarDeploy) readConfYaml() {

	yamlFile, err := ioutil.ReadFile(deathStarDeploy.ConfPath)
	if err != nil {
		deathStarDeploy.DeathLogger.Fatal().Err(err).Msg("Cannot read config")
	}

	var yamlConfig vegetaModels.YAMLConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		deathStarDeploy.DeathLogger.Fatal().Err(err).Msg("Cannot parse yaml config")
	}

	deathStarDeploy.yamlConfig = &yamlConfig
}

// Start function is expected to start with the flow and carry it out till the end
func (deathStarDeploy *DeathStarDeploy) Start() {

	deathStarDeploy.readConfYaml()

	// This is important! In order to create a lambda function we have to provide the
	// zip file containing the handler code. In our case, DeathStar itself is the handler.
	// Depending on the value of the deploy option passed via CLI, we decide whether to call the
	// lambda handler or not. So basically DeathStar needs to send the zip of itself.
	// In order to do so DeathStar provides two options:
	// - pass the option "local" with the CLI. This means you have a local zip of the DeathStar binary.
	//   How this zip was created is upto you. The binary (ELF) might have been built from source or downloaded
	//	 from github. The important part is the binary should be an linux elf or it might not run on AWS.
	// - If "local" option was not passed, then DeathStar will create a zip of its own binary and pass that
	//   zip file path while creating the lambda function. This option is easier and cleaner however can be used
	//   only when DeathStar deployment/attack is being done from a linux machine. This is because for non-linux
	//   machines like OSX, we will be using the OSX binary for running DeathStar. Hence when DeathStar will be creating
	//   its own package,that package will contain the OSX binary which will not work on AWS, yeah its sad :(
	if deathStarDeploy.LocalZip == true {

		// check zip-file-path is present or not
		if deathStarDeploy.ZipFilePath == "" {
			deathStarDeploy.DeathLogger.Fatal().Msg("Code zip file not provided")
		}
	} else {

		// Let DeathStar package itself and use the zip file created
		deathStarDeploy.ZipFilePath = packager.Package(deathStarDeploy.DeathLogger)
	}

	lambdaConfig := deathStarDeploy.yamlConfig.LambdaConfig

	// set defaults for config values
	LambdaRegion := "us-east-1"
	LambdaMemorySize := int64(128)
	LambdaTimeOut := int64(3)

	// get the lambda handler name, it should be the name of this binary
	lambdaFunctionHandler := filepath.Base(os.Args[0])

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
		LambdaFunctionHandler: lambdaFunctionHandler,
		LambdaFunctionRuntime: "go1.x",
		ZipFilePath: deathStarDeploy.ZipFilePath,
		LambdaMemorySize: LambdaMemorySize,
		LambdaTimeOut: LambdaTimeOut,
	}
	err := lambdaUtil.CreateFunction()
	if err != nil {
		deathStarDeploy.DeathLogger.Fatal().Err(err).Msg("Function creation failed")
	}

	deathStarDeploy.DeathLogger.Info().Msg("Function creation succeeded...")

	// initiate attack and display result
	vegAttackUtil := vegetaUtil.VegetaAttackUtils{
		LmUtil: &lambdaUtil,
	}

	deathStarDeploy.DeathLogger.Info().Msg("Running attack")
	lambdaResponses := vegAttackUtil.VegetaSeqAttack(deathStarDeploy.yamlConfig.Attacks, deathStarDeploy.DeathLogger)
	deathStarDeploy.DeathLogger.Info().Msg("Attack complete")

	data, _ := json.MarshalIndent(lambdaResponses, "", "	")
	fmt.Println(string(data))

	deathStarDeploy.DeathLogger.Info().Msg("Cleaning up function...")

	// Clean up the lambda function which we created above to carry out the attack
	err = lambdaUtil.DeleteFunction()
	if err != nil {
		deathStarDeploy.DeathLogger.Fatal().Err(err).Msg("Faced error while deleting function")
	}

	// If the function package zip file was created by us then we have to clean it up
	// at the end of the attack.
	if !deathStarDeploy.LocalZip {
		err = os.Remove(deathStarDeploy.ZipFilePath)
	}
	if err != nil {
		deathStarDeploy.DeathLogger.Error().Err(err).Msg("Failed to remove function package zip file")
	}
}
