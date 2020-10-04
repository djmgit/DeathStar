package vegeta_core

import (
	"fmt"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	models "github.com/djmgit/DeathStar/models"
)

type struct Vegeta_util {
	vegeta_params models.Vegeta_attack_params
}

func (vegeta_util *Vegeta_util) Engage_vegeta() (error, vegeta.Metrics) {
	return nil, ""
}
