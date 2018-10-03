package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/ghodss/yaml"
	"github.com/seeruk/barbara/barbara"
)

// Config holds all application configuration.
type Config struct {
	Primary   barbara.WindowConfig `json:"primary"`
	Secondary barbara.WindowConfig `json:"secondary"`
}

// LoadConfig returns Barbara's configuration. It will either default to a directory under the
// user's home directory, or can be overridden via the environment.
func LoadConfig() (Config, error) {
	var config Config

	usr, err := user.Current()
	if err != nil {
		return config, err
	}

	// If the config path doesn't already exist, create it.
	confPathName := fmt.Sprintf("%s/.config/barbara", usr.HomeDir)
	if _, err := os.Stat(confPathName); os.IsNotExist(err) {
		err := os.MkdirAll(confPathName, os.ModePerm)
		if err != nil {
			return config, err
		}
	}

	confFileName := fmt.Sprintf("%s/config.yml", confPathName)
	confFile, err := os.OpenFile(confFileName, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return config, err
	}

	confBytes, err := ioutil.ReadAll(confFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(confBytes, &config)

	return config, err
}
