package models

import vegeta "github.com/tsenart/vegeta/v12/lib"

type LambdaResponse struct {
	ResultMetrics *vegeta.Metrics `json:"AttackResponseMetrics" yaml:"AttackResponseMetrics"`
	AttackDetails LambdaRequest `json:"attackDetails" yaml:"attackDetails"`
}
