package models

import (
	"time"
)

// struct to hold arguments to vegeta
type VegetaAttackParams struct {
	Method string
	Url string
	Rate int
	Duration time.Duration
}
