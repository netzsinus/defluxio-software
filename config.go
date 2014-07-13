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

type ProviderConfigurationData struct {
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
	Validation struct {
		SpikeThreshold float64
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

func LoadProviderConfiguration(configFile string) (cfg *ProviderConfigurationData, err error) {
	cfg = new(ProviderConfigurationData)
	err = gcfg.ReadFileInto(cfg, configFile)
	if err != nil {
		return nil, fmt.Errorf("Cannot read configuration file: " + err.Error())
	} else {
		return cfg, nil
	}
}
