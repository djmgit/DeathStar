
// models contains struct definitions for various structures used through out
// the DeathStar codebase
package models

// YAMLConfig holds the config attributes for the yaml config file
// required to configure DeathStar in order to carry out attacks.
type  YAMLConfig struct {

	// Attacks is the list of attack configs
	Attacks []LambdaRequest `json:"attacks" yaml:"attacks"`

	// LambdaConfig is the configs for the lambda function itself that
	// will be created in order to carry out the attack.
	LambdaConfig LambdaFunctionConfig `json:"lambdaConfig" yaml:"lambdaConfig"`
}
