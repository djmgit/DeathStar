package deathstardeploy

import (
	"gopkg.in/yaml.v2"
	vegetaModels "github.com/djmgit/DeathStar/models"
	"github.com/djmgit/DeathStar/lambdautil"
	vegetaUtil "github.com/djmgit/DeathStar/vegeta_core"
	"io/ioutil"
	"fmt"
	"errors"
)

type DeathStarDeploy struct {
	ZipFilePath string
	ConfPath string
	LocalZip bool
	yamlConfig *vegetaModels.YAMLConfig
}

// function to read config yaml
func (deathStarDeploy *DeathStarDeploy) readConfYaml() (error) {

	yamlFile, err := ioutil.ReadFile(deathStarDeploy.ConfPath)
	if err != nil {
		fmt.Println("Unable to read conf yaml")
		fmt.Println(err.Error())
		return err
	}

	var yamlConfig vegetaModels.YAMLConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Println("Error parsing the yaml config")
		fmt.Println(err.Error())
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
			err := errors.New("Zip file path not provided")
			fmt.Println(err.Error())
			return err
		}
	} else {
		// donwload zipfile and set zip file path in zipFilePath
	}

	fmt.Println("Creating the lambda attack function...")
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
		fmt.Println("Function creation failed...")
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Function creation succeeded...")

	// initiate attack and display result
	vegAttackUtil := vegetaUtil.VegetaAttackUtils{
		LmUtil: &lambdaUtil,
	}

	fmt.Println("Running attack...")
	err, resultMetrics := vegAttackUtil.VegetaSeqAttack(deathStarDeploy.yamlConfig.Attacks)

	fmt.Println("Attack complete...")
	for _, result := range resultMetrics {
		fmt.Println(result.Latencies.Mean)
		fmt.Println("\n")
		fmt.Println(result.Duration)
		fmt.Println("\n")
		fmt.Println(result.Success)
	}

	fmt.Println("Cleaning up function...")
	err = lambdaUtil.DeleteFunction()
	if err != nil {
		fmt.Println("Faced error while deleting function")
	}

	fmt.Println("Exiting DeathStar...")

	return nil
}
