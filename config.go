// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package defluxio

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type InfluxDBConfig struct {
	Enabled  bool
	Host     string
	Port     int
	Database string
	User     string
	Pass     string
}
type NetworkConfig struct {
	Host string
	Port int
}
type AssetConfig struct {
	ViewPath  string
	AssetPath string
}

type ServerConfiguration struct {
	Meters       Meters
	MeterTimeout time.Duration
	Assets       AssetConfig
	Network      NetworkConfig
	InfluxDB     InfluxDBConfig
}

type DeviceConfig struct {
	Path     string
	Baudrate int
}
type ValidationConfig struct {
	SpikeThreshold float64
}
type ProviderConfiguration struct {
	Meter      Meter
	Network    NetworkConfig
	Device     DeviceConfig
	Validation ValidationConfig
}

type ExporterConfiguration struct {
	InfluxDB InfluxDBConfig
}

/**
 * Exporter Configuration methods
 */

func (pc *ExporterConfiguration) Save(configFile string) (err error) {
	return saveJSON(pc, configFile)
}

func LoadExporterConfiguration(configFile string) (cfg *ExporterConfiguration, err error) {
	cfg = &ExporterConfiguration{}
	err = loadJSON(cfg, configFile)
	return cfg, err
}

func MkDefaultExporterConfiguration() (cfg ExporterConfiguration) {
	cfg = ExporterConfiguration{
		InfluxDB: InfluxDBConfig{
			Enabled:  false,
			Host:     "127.0.0.1",
			Port:     8086,
			Database: "frequency",
			User:     "root",
			Pass:     "root",
		},
	}
	return cfg
}

/**
 * Provider Configuration methods
 */
func (pc *ProviderConfiguration) Save(configFile string) (err error) {
	return saveJSON(pc, configFile)
}

func LoadProviderConfiguration(configFile string) (cfg *ProviderConfiguration, err error) {
	cfg = &ProviderConfiguration{}
	err = loadJSON(cfg, configFile)
	return cfg, err
}

func MkDefaultProviderConfiguration() (cfg ProviderConfiguration) {
	cfg = ProviderConfiguration{
		Meter: Meter{
			Rank:     0,
			ID:       "meter1",
			Key:      "secretkey1",
			Name:     "Meter 1",
			Location: "Somewhere",
		},
		Network: NetworkConfig{
			Host: "http://127.0.0.1",
			Port: 8080,
		},
		Device: DeviceConfig{
			Path:     "/dev/ttyAMA0",
			Baudrate: 115200,
		},
		Validation: ValidationConfig{
			SpikeThreshold: 150,
		},
	}
	return cfg
}

/**
* Server Configuration methods
 */
func (sc *ServerConfiguration) Save(configFile string) (err error) {
	return saveJSON(sc, configFile)
}

func LoadServerConfiguration(configFile string) (cfg *ServerConfiguration, err error) {
	cfg = new(ServerConfiguration)
	err = loadJSON(cfg, configFile)
	if err == nil {
		defaultReading := Reading{Value: 0.0, Timestamp: time.Now()}
		for idx := range cfg.Meters {
			cfg.Meters[idx].Cache = MakeReadingCache(cfg.Meters[idx].CacheSize)
			cfg.Meters[idx].AppendReading(defaultReading)
		}
	}
	return cfg, err
}

func MkDefaultServerConfiguration() (cfg ServerConfiguration) {
	meter1 := Meter{
		Rank:      0,
		ID:        "meter1",
		Key:       "secretkey1",
		Name:      "Meter 1",
		Location:  "Somewhere",
		CacheSize: 10,
	}
	meter2 := Meter{
		Rank:      1,
		ID:        "meter2",
		Key:       "secretkey2",
		Name:      "Meter 2",
		Location:  "Nowhere",
		CacheSize: 10,
	}

	cfg = ServerConfiguration{
		Meters: Meters{
			&meter1,
			&meter2,
		},
		MeterTimeout: 10,
		Network: NetworkConfig{
			Host: "127.0.0.1",
			Port: 8080,
		},
		Assets: AssetConfig{
			ViewPath:  "./views",
			AssetPath: "./assets",
		},
		InfluxDB: InfluxDBConfig{
			Enabled:  false,
			Host:     "127.0.0.1",
			Port:     8086,
			Database: "frequency",
			User:     "root",
			Pass:     "root",
		},
	}
	return cfg
}

// Helper: Saves the given configuration struct into a JSON file.
func saveJSON(configuration interface{}, configFile string) (err error) {
	var bytes []byte
	bytes, err = json.MarshalIndent(configuration, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configFile, bytes, 0644)
	return err
}

// Helper: Loads a configuration from a JSON file.
func loadJSON(configuration interface{}, configFile string) (err error) {
	var bytes []byte
	bytes, err = ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, configuration)
	if err != nil {
		return err
	}
	return nil
}
