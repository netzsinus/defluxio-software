package main

import (
	"code.google.com/p/gcfg"
	"log"
)

type ConfigurationData struct {
	API struct {
		Keys []string
	}
	Network struct {
		Host string
	}
}

var Cfg ConfigurationData

func loadConfiguration(configFile string) {
	err := gcfg.ReadFileInto(&Cfg, configFile)
	if err != nil {
		log.Fatal("Cannot read configuration file: " + err.Error())
	}
}
