// Package config - will take care of reading the JSON config
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type State struct {
	Conf *Config
}

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c Config) SetUser(user string) error {
	c.CurrentUserName = user
	check(write(c))
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func write(cfg Config) error {
	jsonToWrite, err := json.Marshal(cfg)
	check(err)
	configPath, err := getConfigFilePath()
	check(err)
	f, err := os.Create(configPath)
	check(err)
	_, err = f.Write(jsonToWrite)
	check(err)
	return nil
}

func getConfigFilePath() (string, error) {
	const configFileName = ".gatorconfig.json"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed fetching config file path: %s", err)
	}
	return fmt.Sprintf("%s/%s", homeDir, configFileName), nil
}

func readJSONFile[T any](jsonPath string) (T, error) {
	jsonFile, err := os.Open(jsonPath)
	var result T
	if err != nil {
		return result, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return result, err
	}

	if err = json.Unmarshal(byteValue, &result); err != nil {
		return result, err
	}
	return result, nil
}

func ReadConfig() (Config, error) {
	configPath, _ := getConfigFilePath()
	config, err := readJSONFile[Config](configPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to extract config: %s", err)
	}
	return config, nil
}
