package main

import (
	"fmt"

	"github.com/Bravnar/boot_dev_blog_aggregator/internal/config"
)

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
