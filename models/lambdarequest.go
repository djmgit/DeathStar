package models

type LambdaRequest struct {
	AttackName string `json:"attackName" yaml:"attackName"`
	AttackDesc string `json:"attackDesc" yaml:"attackDesc"`
	VegetaParams VegetaAttackParams `json:"vegetaAttackParams" yaml:"vegetaAttackParams"`
}
