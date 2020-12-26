package models

import vegeta "github.com/tsenart/vegeta/v12/lib"

type LambdaResponse struct {
	ResultMetrics *vegeta.Metrics `json:"lambdaRole" yaml:"AttackResponseMetrics"`
	AttackDetails LambdaRequest `json:"attacks" yaml:"attackDetails"`
}
