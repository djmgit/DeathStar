package vegeta_core

import (
	"fmt"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	models "github.com/djmgit/DeathStar/models"
)

type Vegeta_util struct {
	vegeta_params models.Vegeta_attack_params
}

func (vegeta_util *Vegeta_util) init_vegeta(vegeta_params models.Vegeta_attack_params) error {
	vegeta_util.vegeta_params = vegeta_params

	return nil
}

func (vegeta_util *Vegeta_util) Engage_vegeta() (error, vegeta.Metrics) {
}
