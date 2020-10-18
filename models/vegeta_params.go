package models

// struct to hold arguments to vegeta
type Vegeta_attack_params struct {
	Method string
	URL string
	rate int
	duration string
}
