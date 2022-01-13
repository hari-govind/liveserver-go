package config

import "os"

type Config struct {
	// File and directory patters to ignore
	IgnorePatterns []string `yaml:"ignorePatterns"`
	// Root directory to serve and watch
	Root string `yaml:"root"`
	// Levels of subdirectories to watch
	Depth int `yaml:"depth"`
	// Web server listen address
	ListenAddress string `yaml:"listenAddress"`
	WebserverPort int    `yaml:"webserverPort"`
	WebsocketPort int    `yaml:"websocketPort"`
	// Time in ms to wait before reload after file change
	Wait int `yaml:"wait"`
}

func GetDefaultConfig() Config {
	defaultConfig := new(Config)
	defaultConfig.IgnorePatterns = []string{".*", "node_modules", "*~"}
	defaultConfig.Depth = 2 // Root directory and one level down
	defaultConfig.ListenAddress = "127.0.0.1"
	defaultConfig.WebserverPort = 8081
	defaultConfig.WebsocketPort = 8082
	defaultConfig.Wait = 100
	workingDir, err := os.Getwd()
	if err == nil {
		defaultConfig.Root = workingDir
	}
	return *defaultConfig
}
