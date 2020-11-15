package models

type LambdaRequest struct {
	AttackName string `json:"attackName"`
	AttackDesc string `json:"attackDesc"`
	VegetaParams VegetaAttackParams `json:"vegetaAttackParams"`
}
