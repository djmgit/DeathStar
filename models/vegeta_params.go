package models

import (
	"time"
)

// struct to hold arguments to vegeta
type VegetaAttackParams struct {
	Method string	`json:"httpMethod" yaml:"httpMethod"`
	Url string	`json:"url" yaml:"url"`
	Rate int	`json:"rate" yaml:"rate"`
	Duration time.Duration	`json:"duration" yaml:"duration"`
}
