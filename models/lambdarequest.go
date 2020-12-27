package models

// LambdaRequest holds properties to represent a request to the deployed lambda function
type LambdaRequest struct {

	// AttackName if the name of the attack. This can be used as a unique identifier for the attack.
	AttackName string `json:"attackName" yaml:"attackName"`

	// AttackDesc is the description of the attack
	AttackDesc string `json:"attackDesc" yaml:"attackDesc"`

	// VegetaParams contains the data required by vegeta.
	VegetaParams VegetaAttackParams `json:"vegetaAttackParams" yaml:"vegetaAttackParams"`
}
