package main

import (
	"flag"
	"github.com/aws/aws-lambda-go/lambda"
	dsDeploy "github.com/djmgit/DeathStar/deathstardeploy"
	vegetaModels "github.com/djmgit/DeathStar/models"
	vegetaUtil "github.com/djmgit/DeathStar/vegeta_core"
	"github.com/rs/zerolog"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"os"
	"time"
)

// handler for lambda
func HandleLambdaEvent(event vegetaModels.LambdaRequest) (vegeta.Metrics, error) {

	vegetaAttacker := vegetaUtil.VegetaUtil {
		VegetaParams: event.VegetaParams,
	}

	_, metrics := vegetaAttacker.EngageVegeta()

	return metrics, nil
}

func setupLogging(loglevel string) zerolog.Logger {

	LOGLEVELS := map[string]zerolog.Level{
		"trace": zerolog.TraceLevel,
		"debug": zerolog.DebugLevel,
		"info": zerolog.InfoLevel,
		"warn": zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"fatal": zerolog.FatalLevel,
		"panic": zerolog.PanicLevel,
	}

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(LOGLEVELS[loglevel])

	deathLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return deathLogger
}

func main() {

	var isLocal bool
	var zipFilePath string
	var confPath string
	var deploy bool
	var loglevel string

	flag.BoolVar(&isLocal, "local", false, "Denotes and DeathStar will use local zip. Zip path must be profiled")
	flag.StringVar(&zipFilePath, "zip-file-path", "./func.zip", "Path to local Zip file containing handler")
	flag.StringVar(&confPath, "conf", "deathstar.conf", "Path to conf file")
	flag.BoolVar(&deploy, "deploy", false, "Deploy function to lambda and initiate attack")
	flag.StringVar(&loglevel, "loglevel", "debug", "Set loglevel")

	flag.Parse()

	deathlogger := setupLogging(loglevel)

	if deploy == true {
		dsDeployHandler := dsDeploy.DeathStarDeploy {
			ZipFilePath: zipFilePath,
			ConfPath: confPath,
			LocalZip: isLocal,
			DeathLogger: deathlogger,
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
