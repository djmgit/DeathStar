package main

import (
	"fmt"
	vegetaUtil "github.com/djmgit/DeathStar/vegeta_core"
	vegetaModels "github.com/djmgit/DeathStar/models"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// handler for lambda
func HandleLambdaEvent(event vegetaModels.LambdaRequest) (vegeta.Metrics, error) {

	vegetaAttacker := vegetaUtil.VegetaUtil {
		VegetaParams: event.VegetaAttackParams
	}

	_, metrics := vegetaAttacker.EngageVegeta()

	return metrics, nil
}

func main() {
	// locally testing vegeta

	params := vegetaModels.VegetaAttackParams {
		Method: "GET",
		Url: "https://google.com",
		Rate: 100,
		Duration: 10,
	}

	vegetaAttacker := vegetaUtil.VegetaUtil {
		VegetaParams: params,
	}

	_, metrics := vegetaAttacker.EngageVegeta()

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99);
}
