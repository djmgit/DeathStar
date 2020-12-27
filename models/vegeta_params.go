package models

import (
	"time"
)

// VegetaAttackParams holds fields for vegeta attack config.
// These data are required by vegeta to carry out the attack
// The meaning and purpose of the fields can be found in the vegeta documentaion

type VegetaAttackParams struct {
	Method string	`json:"httpMethod" yaml:"httpMethod"`
	Url string	`json:"url" yaml:"url"`
	Rate int	`json:"rate" yaml:"rate"`
	Duration time.Duration	`json:"duration" yaml:"duration"`
}
