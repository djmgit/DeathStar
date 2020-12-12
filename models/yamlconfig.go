package models

type  YAMLConfig struct {
	Attacks []LambdaRequest `json:"attacks" yaml:"attacks"`
	LambdaConfig LambdaFunctionConfig `json:"lambdaConfig" yaml:"lambdaConfig"`
}
