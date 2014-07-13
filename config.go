package defluxio

import (
	"code.google.com/p/gcfg"
	"fmt"
)

type InfluxDBConfig struct {
	Enabled  bool
	Host     string
	Port     int
	Database string
	User     string
	Pass     string
}

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
	InfluxDB InfluxDBConfig
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

type ExporterConfigurationData struct {
	InfluxDB InfluxDBConfig
}

func LoadServerConfiguration(configFile string) (cfg *ServerConfigurationData, err error) {
	cfg = new(ServerConfigurationData)
	err = gcfg.ReadFileInto(cfg, configFile)
	if err != nil {
		return nil, fmt.Errorf("Cannot read configuration file: " + err.Error())
	} else {
		return cfg, nil
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

func LoadExporterConfiguration(configFile string) (cfg *ExporterConfigurationData, err error) {
	cfg = new(ExporterConfigurationData)
	err = gcfg.ReadFileInto(cfg, configFile)
	if err != nil {
		return nil, fmt.Errorf("Cannot read configuration file: " + err.Error())
	} else {
		return cfg, nil
	}
}
