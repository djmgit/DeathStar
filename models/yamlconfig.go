package models

type  YAMLConfig struct {
	Attacks []LambdaRequest `json:"attacks" yaml:"attacks"`
}
