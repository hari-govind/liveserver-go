package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

const CONFIG_FILENAME = ".liveserver.yaml"

var liveServerConfig *Config

func GetConfig() Config {
	if liveServerConfig != nil {
		return *liveServerConfig
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Cannot get working directory")
	}

	// Have default values initially, and override values as found in user config file
	liveServerConfig = new(Config)
	*liveServerConfig = GetDefaultConfig()

	// Search order: current directory, home directory

	configPath := path.Join(cwd, CONFIG_FILENAME)
	_, err = os.Stat(configPath)
	if err == nil {
		log.Println("Found config file in location ", configPath)
		readConfig(configPath, liveServerConfig)
		return *liveServerConfig
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		configPath = path.Join(homeDir, CONFIG_FILENAME)
		_, err = os.Stat(configPath)
		if err == nil {
			log.Println("Found config file in location ", configPath)
			readConfig(configPath, liveServerConfig)
			return *liveServerConfig
		}
	}

	log.Println("Did not find configuration file, using default configuration.")
	return *liveServerConfig
}

func readConfig(filePath string, liveServerConfig *Config) {
	configFileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Cannot read configuration file.")
		log.Print(err)
	}

	err = yaml.Unmarshal(configFileContents, liveServerConfig)
	if err != nil {
		log.Println("Error parsing configuration file.")
		log.Print(err)
	}

}
