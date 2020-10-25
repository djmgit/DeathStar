package lambdautil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/lambda"
    "fmt"
    //"os"
)

type LambdaUtil struct {
	AWSRegion string `json:"awsRegion"`
	LambdaRole string `json:"lambdaRole"`
	LambdaFuncName string `json:"lambdaFuncName"`
	LambdaFunctionHandler string `json:"lambdaFunctionHandler"`
	LambdaFunctionRuntime string `json:"lambdaFunctionRuntime"`
	ZipFilePath string `json:"zipFilePath"`
	AWSAccessKeyID string `json:"awsAccessKeyID"`
	AWSSecretAccessKey string `json:"awsSecretAccessKey`
	awsSession *session.Session
}

func (lambdautil *LambdaUtil) getAWSSession() (error) {

	// create the aws session and set it as struct property
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(lambdautil.AWSRegion),
	})

	if err != nil {
		// shared config not set, fall back to provided creds
		sess, err = session.NewSession(&Aws.Config{
			Region: aws.String(lambdautil.AWSRegion),
			Credentials: credentials.NewStaticCredentials(lambdautil.AWSAccessKeyID, lambdautil.AWSSecretAccessKey),
		})
	}

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	lambdautil.awsSession = sess

	return nil
}

func (lambdaUtil *LambdaUtil) CreateFunction() error {

	// create the lambda function using the provided informations
	if lambdautil.awsSession != nil {
		err := lambda.getAWSSession()
		if err != nil {
			return err
		}
	}

	svc := lambda.New(lambdautil.awsSession)

	createCode := &lambda.FunctionCode{
		ZipFile:         contents,
	}

	createArgs := &lambda.CreateFunctionInput{
		Code:         createCode,
		FunctionName: &LambdaFuncName,
		Handler:      &LambdaFunctionHandler,
		Role:         &LambdaRole,
		Runtime:      &LambdaFunctionRuntime,
	}

	result, err := svc.CreateFunction(createArgs)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
