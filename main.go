package main

import (
	"flag"
	vegetaUtil "github.com/djmgit/DeathStar/vegeta_core"
	vegetaModels "github.com/djmgit/DeathStar/models"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"github.com/aws/aws-lambda-go/lambda"
	dsDeploy "github.com/DeathStar/deathstardeploy"
)

// handler for lambda
func HandleLambdaEvent(event vegetaModels.LambdaRequest) (vegeta.Metrics, error) {

	vegetaAttacker := vegetaUtil.VegetaUtil {
		VegetaParams: event.VegetaParams,
	}

	_, metrics := vegetaAttacker.EngageVegeta()

	return metrics, nil
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

	if deploy == true {
		dsDeployHandler := dsDeploy.DeathStarDeploy {
			ZipFilePath: zipFilePath,
			ConfPath: confPath,
			LocalZip: isLocal,
		}

		err := dsDeployHandler.Start()

		if err != nil {
			os.Exit(2)
		}

		os.Exit(0)
	}

	// Fucntion is running in lambda, initialte handler
	lambda.Start(HandleLambdaEvent)

}
