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

// HandleLambdaEvent is the main lambda function handler called
// by the aws lambda Go SDK's lambda.Start() function. We pass this handler
// as a argument to the lambda.Start() function
func HandleLambdaEvent(event vegetaModels.LambdaRequest) (vegeta.Metrics, error) {

	vegetaAttacker := vegetaUtil.VegetaUtil {
		VegetaParams: event.VegetaParams,
	}

	// Start vegeta attack
	_, metrics := vegetaAttacker.EngageVegeta()

	return metrics, nil
}

// setupLogging sets up and returns the logger used by other
// parts in the codebase
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

	// DeathStar can be invoked in two situations
	// - It needs to perform an attack, in that case we simply have to invoke
	//   lambda.Start() method and rest the handler will take care. This basically means
	//   we are inside the lambda and we are running on AWS.
	// - Second, when death star is being invoked as a CLI tool for launching an end to end attack, that is,
	//   reading the provided config, creating the lambda function, invoking the lambda function for carying out
	//   the attack, and lastly cleaning up the lambda function.
	// For carrying out these two conditions, we check the deploy variable. If its true, then we follow the end to end flow,
	// create function, carry out attack, and do cleanup. If deploy variable is false then we let the handler to its job.
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

		deathlogger.Info().Msg("Exiting DeathStar...")
		os.Exit(0)
	}

	// Fucntion is running in lambda, invoke handler
	lambda.Start(HandleLambdaEvent)
}
