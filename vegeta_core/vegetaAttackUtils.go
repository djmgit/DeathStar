package vegeta_core

import (
	"github.com/djmgit/DeathStar/lambdautil"
	vegetaModels "github.com/djmgit/DeathStar/models"
	"github.com/rs/zerolog"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type VegetaAttackUtils struct {
	LmUtil *lambdautil.LambdaUtil
}

func (vegUtil *VegetaAttackUtils) VegetaSeqAttack(attackConfigs []vegetaModels.LambdaRequest, deathLogger zerolog.Logger) (error, []vegetaModels.LambdaResponse) {

	var resultMetrics []*vegeta.Metrics
	var lambdaResponses []vegetaModels.LambdaResponse
	noFailures := true

	for _, attackConfig := range attackConfigs {

		// conduct attact for the given request
		err, attackResult := vegUtil.LmUtil.RunFunction(attackConfig)
		if err == nil {
			resultMetrics = append(resultMetrics, attackResult)
			lambdaResponses = append(lambdaResponses, vegetaModels.LambdaResponse{
				ResultMetrics: attackResult,
				AttackDetails: attackConfig,
			})
		} else {
			deathLogger.Error().Msg(err.Error())
			noFailures = false
		}
	}

	if !noFailures {
		deathLogger.Info().Msg("One or more attacks could not be carried out successfully, please check logs!")
	}

	return nil, lambdaResponses
}
