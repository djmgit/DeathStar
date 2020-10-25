package models

type struct LambdaRequest {
	AttackName string `json:"attackName"`
	AttackDesc string `json:"attckDesc"`
	VegetaParams VegetaAttackParams `json:"vegetaAttackParams"`
}
