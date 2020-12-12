package models

type LambdaFunctionConfig struct {
	LambdaRole string `json:"lambdaRole" yaml:"lambdaRole"`
	LambdaMemorySize int `json:"lambdaMemorySize" yaml:"lambdaMemorySize"`
	LambdaTimeOut int `json:"lambdaTimeOut" yaml:"lambdaTimeOut"`
}
