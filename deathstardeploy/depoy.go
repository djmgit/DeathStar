package deathstardeploy

import (
	"gopkg.in/yaml.v2"
	vegetaModels "github.com/djmgit/DeathStar/models"
	"io/ioutil"
	"fmt"
	"errors"
)

type DeathStarDeploy struct {
	ZipFilePath string
	ConfPath string
	LocalZip bool
	yamlConfig *vegetaModels.YAMLConfig
}

// function to read config yaml
func (deathStarDeploy *DeathStarDeploy) readConfYaml() (error) {

	yamlFile, err := ioutil.ReadFile(deathStarDeploy.ConfPath)
	if err != nil {
		fmt.Println("Unable to read conf yaml")
		fmt.Println(err.Error())
		return err
	}

	var yamlConfig vegetaModels.YAMLConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Println("Error parsing the yaml config")
		fmt.Println(err.Error())
		return err
	}

	deathStarDeploy.yamlConfig = &yamlConfig
	return nil
}

func (deathStarDeploy *DeathStarDeploy) Start() error {

	err := deathStarDeploy.readConfYaml()
	if err != nil {
		return err
	}

	if deathStarDeploy.LocalZip == true {

		// check zip-file-path is present or not
		if deathStarDeploy.ZipFilePath == "" {
			err := errors.New("Zip file path not provided")
			fmt.Println(err.Error())
			return err
		}
	} else {
		// donwload zipfile and set zip file path in zipFilePath
	}

	return nil
}
