// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package defluxio

import (
	"io/ioutil"
	"reflect"
	"syscall"
	"testing"
)

func TestServerConfigSerialization(t *testing.T) {
	sc := MkDefaultServerConfiguration()
	f, err := ioutil.TempFile("", "testserverconfig.json")
	if err != nil {
		t.Error("Cannot create temp file")
	}
	defer syscall.Unlink(f.Name())
	err = sc.Save(f.Name())
	if err != nil {
		t.Error("Cannot save server configuration: " + err.Error())
	}
	sc2, err := LoadServerConfiguration(f.Name())
	if err != nil {
		t.Error("Failed to load server configuration: " + err.Error())
	}
	// Note: This test currently fails because the reading cache is
	// generated dynamically after each load.
	//if !reflect.DeepEqual(sc, *sc2) {
	//	t.Error("Pre-stored server configuration not equal to post-stored one")
	//}
	if !(sc.Network == sc2.Network && sc.Assets == sc2.Assets &&
		sc.InfluxDB == sc2.InfluxDB) {
		t.Error("Pre-stored server configuration not equal to post-stored one")
	}
}

func TestProviderConfigSerialization(t *testing.T) {
	sc := MkDefaultProviderConfiguration()
	f, err := ioutil.TempFile("", "testproviderconfig.json")
	if err != nil {
		t.Error("Cannot create temp file")
	}
	defer syscall.Unlink(f.Name())
	err = sc.Save(f.Name())
	if err != nil {
		t.Error("Cannot save server configuration: " + err.Error())
	}
	sc2, err := LoadProviderConfiguration(f.Name())
	if err != nil {
		t.Error("Failed to load server configuration: " + err.Error())
	}
	if !reflect.DeepEqual(sc, *sc2) {
		t.Error("Pre-stored server configuration not equal to post-stored one")
	}
}

func TestExporterConfigSerialization(t *testing.T) {
	sc := MkDefaultExporterConfiguration()
	f, err := ioutil.TempFile("", "testexporterconfig.json")
	if err != nil {
		t.Error("Cannot create temp file")
	}
	defer syscall.Unlink(f.Name())
	err = sc.Save(f.Name())
	if err != nil {
		t.Error("Cannot save server configuration: " + err.Error())
	}
	sc2, err := LoadExporterConfiguration(f.Name())
	if err != nil {
		t.Error("Failed to load server configuration: " + err.Error())
	}
	if !reflect.DeepEqual(sc, *sc2) {
		t.Error("Pre-stored server configuration not equal to post-stored one")
	}
}
