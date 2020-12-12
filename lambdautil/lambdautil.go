package lambdautil

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	vegetaModels "github.com/djmgit/DeathStar/models"
	"github.com/rs/zerolog"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"io/ioutil"
)

type LambdaUtil struct {
	AWSRegion string `json:"awsRegion"`
	LambdaRole string `json:"lambdaRole"`
	LambdaFuncName string `json:"lambdaFuncName"`
	LambdaFunctionHandler string `json:"lambdaFunctionHandler"`
	LambdaFunctionRuntime string `json:"lambdaFunctionRuntime"`
	LambdaMemorySize int64 `json:"lambdaMemorySize" yaml:"lambdaMemorySize"`
	LambdaTimeOut int64 `json:"lambdaTimeOut" yaml:"lambdaTimeOut"`
	ZipFilePath string `json:"zipFilePath"`
	AWSAccessKeyID string `json:"awsAccessKeyID"`
	AWSSecretAccessKey string `json:"awsSecretAccessKey`
	DeathLogger zerolog.Logger
	AWSSession *session.Session
}

func (lambdaUtil *LambdaUtil) GetAWSSession() (error) {

	// create the aws session and set it as struct property
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(lambdaUtil.AWSRegion),
	})

	if err != nil {
		// shared config not set, fall back to provided creds
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(lambdaUtil.AWSRegion),
			Credentials: credentials.NewStaticCredentials(lambdaUtil.AWSAccessKeyID, lambdaUtil.AWSSecretAccessKey, ""),
		})
	}

	if err != nil {
		return err
	}

	lambdaUtil.AWSSession = sess

	return nil
}

func (lambdaUtil *LambdaUtil) CreateFunction() error {

	// create the lambda function using the provided informations
	if lambdaUtil.AWSSession == nil {
		err := lambdaUtil.GetAWSSession()
		if err != nil {
			return err
		}
	}

	svc := lambda.New(lambdaUtil.AWSSession)

	lambdaFuncContents, err := ioutil.ReadFile(lambdaUtil.ZipFilePath)
	if err != nil {
		return err
	}

	createCode := &lambda.FunctionCode{
		ZipFile:         lambdaFuncContents,
	}

	createArgs := &lambda.CreateFunctionInput{
		Code:         createCode,
		FunctionName: &lambdaUtil.LambdaFuncName,
		Handler:      &lambdaUtil.LambdaFunctionHandler,
		Role:         &lambdaUtil.LambdaRole,
		Runtime:      &lambdaUtil.LambdaFunctionRuntime,
		MemorySize:   &lambdaUtil.LambdaMemorySize,
		Timeout:	  &lambdaUtil.LambdaTimeOut,
	}

	_, err = svc.CreateFunction(createArgs)

	if err != nil {
		return err
	}

	return nil
}

func (lambdaUtil *LambdaUtil) RunFunction(lambdaRequest vegetaModels.LambdaRequest) (error, *vegeta.Metrics) {

	// create the lambda function using the provided informations
	if lambdaUtil.AWSSession == nil {
		err := lambdaUtil.GetAWSSession()
		if err != nil {
			return err, &vegeta.Metrics{}
		}
	}

	svc := lambda.New(lambdaUtil.AWSSession)

	payload, err := json.Marshal(lambdaRequest)
	if err != nil {
		return err, &vegeta.Metrics{}
	}

	result, err := svc.Invoke(&lambda.InvokeInput{FunctionName: aws.String(lambdaUtil.LambdaFuncName), Payload: payload})
	if err != nil {
		return err, &vegeta.Metrics{}
	}

	var response vegeta.Metrics

	err = json.Unmarshal(result.Payload, &response)
	if err != nil {
		return err, &vegeta.Metrics{}
	}

	return nil, &response
}

func (lambdaUtil *LambdaUtil) DeleteFunction() error {

	// create the lambda function using the provided informations
	if lambdaUtil.AWSSession == nil {
		err := lambdaUtil.GetAWSSession()
		if err != nil {
			return err
		}
	}

	svc := lambda.New(lambdaUtil.AWSSession)

	// create the deletion input
	deleteInput := &lambda.DeleteFunctionInput {
		FunctionName: aws.String(lambdaUtil.LambdaFuncName),
	}

	_, err := svc.DeleteFunction(deleteInput)

	return err
}
