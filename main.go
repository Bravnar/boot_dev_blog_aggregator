package main

import (
	"fmt"

	"github.com/Bravnar/gator/internal/config"
)

type State struct {
	Config *config.Config
}

func main() {
	configFile, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("error occured while reading: %s", err)
		return
	}
	configFile.SetUser("bravnar")
	configFile, err = config.ReadConfig()
	if err != nil {
		fmt.Printf("error occured while reading: %s", err)
	}
	fmt.Println(configFile)
}
