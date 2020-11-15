package main

import (
	"fmt"
	"os"
	"flag"
	vegetaUtil "github.com/djmgit/DeathStar/vegeta_core"
	vegetaModels "github.com/djmgit/DeathStar/models"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"github.com/djmgit/DeathStar/lambdautil"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

// handler for lambda
func HandleLambdaEvent(event vegetaModels.LambdaRequest) (vegeta.Metrics, error) {

	vegetaAttacker := vegetaUtil.VegetaUtil {
		VegetaParams: event.VegetaParams,
	}

	_, metrics := vegetaAttacker.EngageVegeta()

	return metrics, nil
}

// function to read config yaml
func readConfYaml(confPath string) (error, *vegetaModels.YAMLConfig) {

	yamlFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Println("Unable to read conf yaml")
		fmt.Println(err.Error())
		return err, nil
	}

	var yamlConfig vegetaModels.YAMLConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Println("Error parsing the yaml config")
		fmt.Println(err.Error())
		return err, nil
	}

	return nil, &yamlConfig
}

func main() {

	var isLocal bool
	var zipFilePath string
	var confPath string
	var deploy bool

	flag.BoolVar(&isLocal, "local", false, "Denotes and DeathStar will use local zip. Zip path must be profiled")
	flag.StringVar(&zipFilePath, "zip-file-path", "./func.zip", "Path to local Zip file containing handler")
	flag.StringVar(&confPath, "conf", "deathstar.conf", "Path to conf file")
	flag.BoolVar(&deploy, "deploy", false, "Deploy function to lambda and initiate attack")

	flag.Parse()

	// read the yaml config
	err, yamlConfig := readConfYaml(confPath)
	if err != nil {
		os.Exit(2)
	}

	if deploy == true {

		if isLocal == true {

			// check zip-file-path is present or not
			if zipFilePath == "" {
				fmt.Println("zipfile path not providing. Exiting DeathStar...")
				os.Exit(2)
			}
		} else {
			// donwload zipfile and set zip file path in zipFilePath
		}

		// create the lambda
		lambdaUtil := lambdautil.LambdaUtil {
			AWSRegion: "us-east-1",
			LambdaRole: "arn:aws:iam::253708721073:role/service-role/func-test-1-role-nyalwdp2",
			LambdaFuncName: "func-test-2",
			LambdaFunctionHandler: "main",
			LambdaFunctionRuntime: "go1.x",
			ZipFilePath: zipFilePath,
		}

		err := lambdaUtil.CreateFunction()

		if err != nil {
			fmt.Println(err.Error())
		}

		// initiate attack and display result
		vegAttackUtil := vegetaUtil.VegetaAttackUtils{
			LmUtil: &lambdaUtil,
		}
		err, resultMetrics := vegAttackUtil.VegetaSeqAttack(yamlConfig.Attacks)
		for _, result := range resultMetrics {
			fmt.Println("%v", *result)
		}

		// clean up function

		// exit DeathStar

		os.Exit(0)
	}

	// Fucntion is running in lambda, initialte handler
	lambda.Start(HandleLambdaEvent)

}
