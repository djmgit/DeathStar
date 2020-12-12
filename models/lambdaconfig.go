package models

type LambdaFunctionConfig struct {
	LambdaRole string `json:"lambdaRole" yaml:"lambdaRole"`
	LambdaMemorySize int64 `json:"lambdaMemorySize" yaml:"lambdaMemorySize"`
	LambdaTimeOut int64 `json:"lambdaTimeOut" yaml:"lambdaTimeOut"`
	LambdaRegion string `json:"lambdaRegion" yaml:"lambdaRegion"`
	LambdaName string `json:"lambdaName" yaml:"lambdaName"`
}
