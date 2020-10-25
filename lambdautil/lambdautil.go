package lambdautil

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/lambda"
    "fmt"
    //"os"
)

type LambdaUtil struct {
	AWSRegion string `json:"awsRegion"`
	LambdaRole string `json:"lambdaRole"`
	LambdaFuncName string `json:"lambdaFuncName"`
	ZipFilePath string `json:"zipFilePath"`
}
