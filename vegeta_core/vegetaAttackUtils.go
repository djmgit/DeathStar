package vegeta_core

import (
	"github.com/djmgit/DeathStar/lambdautil"
	vegetaModels "github.com/djmgit/DeathStar/models"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type VegetaAttackUtils struct {
	LmUtil *lambdautil.LambdaUtil
}

func (vegUtil *VegetaAttackUtils) VegetaSeqAttack(attackConfigs []vegetaModels.LambdaRequest) (error, []vegetaModels.LambdaResponse) {

	var resultMetrics []*vegeta.Metrics
	var lambdaResponses []vegetaModels.LambdaResponse

	for _, attackConfig := range attackConfigs {

		// conduct attact for the given request
		err, attackResult := vegUtil.LmUtil.RunFunction(attackConfig)
		if err == nil {
			resultMetrics = append(resultMetrics, attackResult)
			lambdaResponses = append(lambdaResponses, vegetaModels.LambdaResponse{
				ResultMetrics: attackResult,
				AttackDetails: attackConfig,
			})
		}
	}

	return nil, lambdaResponses
}
