package defluxio

import (
	"code.google.com/p/gcfg"
	"fmt"
)

type ServerConfigurationData struct {
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

func LoadServerConfiguration(configFile string) (*ServerConfigurationData, error) {
	retval := new(ServerConfigurationData)
	err := gcfg.ReadFileInto(retval, configFile)
	if err != nil {
		return nil, fmt.Errorf("Cannot read configuration file: " + err.Error())
	} else {
		return retval, nil
	}
}
