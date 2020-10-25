package vegeta_core

import (
	//"fmt"
	"time"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	models "github.com/djmgit/DeathStar/models"
)

type VegetaUtil struct {
	VegetaParams models.VegetaAttackParams
}

// validate the params sent to vegeta
func (vegetaUtil *VegetaUtil) validateVegetaParams() bool {
	return true
}

func (vegetaUtil *VegetaUtil) initVegeta() (error, bool) {

	if !vegetaUtil.validateVegetaParams() {
		return nil, false
	}

	return nil, true
}

func (vegetaUtil *VegetaUtil) EngageVegeta() (error, vegeta.Metrics) {

	// prepare vegeta for attack

	err, _ := vegetaUtil.initVegeta()

	if err != nil {
		// print errpr
		return err, vegeta.Metrics{}
	}

	// create params
	rate := vegeta.Rate{Freq: vegetaUtil.VegetaParams.Rate, Per: time.Second}
	duration := vegetaUtil.VegetaParams.Duration * time.Second

	// create target
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: vegetaUtil.Method,
		URL: vegetaUtil.VegetaParams.Url,
	})
	attacker := vegeta.NewAttacker()

	// get the result metrics
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	return nil, metrics
}
