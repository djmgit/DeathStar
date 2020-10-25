package models

type LambdaRequest struct {
	AttackName string `json:"attackName"`
	AttackDesc string `json:"attckDesc"`
	VegetaParams VegetaAttackParams `json:"vegetaAttackParams"`
}
