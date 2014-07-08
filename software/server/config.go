package main

import (
	"code.google.com/p/gcfg"
	"log"
)

type ConfigurationData struct {
	API struct {
		Keys []string
	}
	Assets struct {
		ViewPath  string
		AssetPath string
	}
	Network struct {
		Host string
		Port int
	}
	InfluxDB struct {
		Enabled  bool
		Host     string
		Port     int
		Database string
		User     string
		Pass     string
	}
}

var Cfg ConfigurationData

func loadConfiguration(configFile string) {
	err := gcfg.ReadFileInto(&Cfg, configFile)
	if err != nil {
		log.Fatal("Cannot read configuration file: " + err.Error())
	}
}
