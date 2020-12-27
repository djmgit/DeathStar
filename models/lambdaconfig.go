package models

// LambdaFunctionConfig contains fields required to configutre the
// lambda function
type LambdaFunctionConfig struct {

	// LambdaRole is the role that will be used by the lambda function. This will
	// hild the role ARN
	LambdaRole string `json:"lambdaRole" yaml:"lambdaRole"`

	// LambdaMemorySize is the max memory size that the lambda function will be allowed
	// to attain.
	LambdaMemorySize int64 `json:"lambdaMemorySize" yaml:"lambdaMemorySize"`

	// LambdaTimeOut is the max time in seconds after which the lambda function will timeout
	LambdaTimeOut int64 `json:"lambdaTimeOut" yaml:"lambdaTimeOut"`

	// LambdaRegion is the region where the lambda function will be deployed
	LambdaRegion string `json:"lambdaRegion" yaml:"lambdaRegion"`
}
