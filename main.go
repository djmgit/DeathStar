package main

import (
	//"fmt"
	vegetaUtil "github.com/djmgit/DeathStar/vegeta_core"
	vegetaModels "github.com/djmgit/DeathStar/models"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"github.com/aws/aws-lambda-go/lambda"
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

	flag.BoolVar(&isLocal, "local", false, "Denotes and DeathStar will use local zip. Zip path must be profiled")
	flag.StringVar(&zipFilePath, "zip-file-path", "./func.zip", "Path to local Zip file containing handler")
	flag.StringVar(&confPath, "conf", "deathstar.conf", "Path to conf file")

	flag.Parse()

}
