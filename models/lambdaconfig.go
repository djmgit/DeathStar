package models

type LambdaFunctionConfig struct {
	LambdaRole string `json:"lambdaRole" yaml:"lambdaRole"`
	LambdaMemorySize int `json:"lambdaMemorySize" yaml:"lambdaMemorySize"`
	LambdaTimeOut int `json:"lambdaTimeOut" yaml:"lambdaTimeOut"`
	LambdaRegion string `json:"lambdaRegion" yaml:"lambdaRegion"`
	LambdaName string `json:"lambdaName" yaml:"lambdaName"`
}
