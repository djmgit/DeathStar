package lambdautil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/lambda"
    "fmt"
    "io/ioutil"
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
		fmt.Println(err.Error())
		return err
	}

	lambdaUtil.AWSSession = sess

	return nil
}

func (lambdaUtil *LambdaUtil) CreateFunction() error {

	// create the lambda function using the provided informations
	if lambdaUtil.AWSSession != nil {
		err := lambdaUtil.GetAWSSession()
		if err != nil {
			return err
		}
	}

	svc := lambda.New(lambdaUtil.AWSSession)

	lambdaFuncContents, err := ioutil.ReadFile(lambdaUtil.LambdaFuncName)
	if err != nil {
		fmt.Println("Could not read " + lambdaUtil.LambdaFuncName + ".zip")
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
	}

	_, err = svc.CreateFunction(createArgs)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
