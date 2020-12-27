package models

import vegeta "github.com/tsenart/vegeta/v12/lib"

// LambdaRequest holds fields tp represent the attack response.
type LambdaResponse struct {

	// ResultMetrics is same as the metrics struct returned by the vegeta llibrary
	// More about this can be found in vegeta doc.
	ResultMetrics *vegeta.Metrics `json:"AttackResponseMetrics" yaml:"AttackResponseMetrics"`

	// AttackDetails contains the details of a particualar attack. This is included here
	// so that user knows the above metrics are for which attack.
	AttackDetails LambdaRequest `json:"attackDetails" yaml:"attackDetails"`
}
