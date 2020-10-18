package models

// struct to hold arguments to vegeta
type Vegeta_attack_params struct {
	Method string
	Url string
	Rate int
	Duration int
}
