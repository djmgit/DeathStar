package main

import (
	"fmt"
	vegetaUtil "github.com/djmgit/DeathStar/vegeta_core"
	vegetaModels "github.com/djmgit/DeathStar/models"
)

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
