package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// https://stackoverflow.com/questions/16465705/how-to-handle-configuration-in-go
type groupConfig struct {
	BotName string
	Name    string
	Users   []int
}

func (g *groupConfig) addUser(id int) {
	g.Users = append(g.Users, id)
}

type botConfig struct {
	Token string
	Name  string
}

type configuration struct {
	Bots   []botConfig
	Groups []groupConfig
}

// readConfiguration reads the configuration from the default configuration path
// and returns the configuration object or an empty configuration object
// if no such file exists.
func readConfiguration() (configuration, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return configuration{}, err
	}

	configFilePath := filepath.Join(homeDir, ".telebotlog.conf")

	// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		if *verbose {
			fmt.Printf("config [%s] does not exist\n", configFilePath)
		}
		return configuration{}, nil
	}

	// https://stackoverflow.com/questions/16465705/how-to-handle-configuration-in-go
	file, err := os.Open(configFilePath)
	if err != nil {
		if *verbose {
			fmt.Printf("could not open config [%s]\n", configFilePath)
		}
		return configuration{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		if *verbose {
			fmt.Printf("could not decode config [%s]\n", configFilePath)
		}
		return configuration{}, err
	}

	return conf, nil
}

// writeConfiguration tries to write the configuration to the default configuration
// path
func writeConfiguration(conf configuration) error {
	homeDir, err := homedir.Dir()
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(homeDir, ".telebotlog.conf")

	f, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	return encoder.Encode(conf)
}
