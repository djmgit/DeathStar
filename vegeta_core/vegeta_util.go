package vegeta_core

import (
	"fmt"
	"time"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	models "github.com/djmgit/DeathStar/models"
)

type Vegeta_util struct {
	vegeta_params models.Vegeta_attack_params
}

// validate the params sent to vegeta
func (vegeta_util *Vegeta_util) validate_vegeta_params() bool {
	return true
}

func (vegeta_util *Vegeta_util) init_vegeta(vegeta_params models.Vegeta_attack_params) (error, bool) {
	vegeta_util.vegeta_params = vegeta_params

	// validate params
	if !vegeta_util.validate_vegeta_params() {
		return nil, false
	}

	return nil, true
}

func (vegeta_util *Vegeta_util) Engage_vegeta() (error, vegeta.Metrics) {

	// prepare vegeta for attack

	err, status := vegeta_util.init()

	if err != nil {
		// print errpr
		return err, nil
	}

	// create params
	rate := vegeta.Rate{Freq: vegeta_util.vageta_params.Rate, Per: time.Second}
	duration := vegeta_util.vegeta_params.Duration * time.second

	// create target
	targetter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL: vegeta_util.vegeta_params.Url
	})
	attacker := vegeta.NewAttacker()

	// get the result metrics

}
