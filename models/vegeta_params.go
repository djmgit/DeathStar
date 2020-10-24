package models

import (
	"time"
)

// struct to hold arguments to vegeta
type VegetaAttackParams struct {
	Method string	`json:"httpMethod"`
	Url string	`json:"url"`
	Rate int	`json:"rate"`
	Duration time.Duration	`json:"duration"`
}
