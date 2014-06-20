package main

import (
	"code.google.com/p/gcfg"
	"log"
)

type ConfigurationData struct {
	API struct {
		Key   string
		Meter string
	}
	Network struct {
		Host string
	}
	Device struct {
		Path     string
		Baudrate int
	}
}

var Cfg ConfigurationData

func loadConfiguration(configFile string) {
	err := gcfg.ReadFileInto(&Cfg, configFile)
	if err != nil {
		log.Fatal("Cannot read configuration file: " + err.Error())
	}
}
