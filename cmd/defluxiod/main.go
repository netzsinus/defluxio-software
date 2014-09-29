// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/netzsinus/defluxio-software"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

var configFile = flag.String("config", "defluxiod.conf", "configuration file")
var mkConfigFile = flag.Bool("genconfig", false, "generate an example configuration file")
var Cfg *defluxio.ServerConfiguration
var templates *template.Template
var dbclient *defluxio.DBClient
var dbchannel chan defluxio.MeterReading

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "index", r.Host)
}

func serveImpressum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "impressum", r.Host)
}

func serveMeter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	type TemplateData struct {
		Meters defluxio.Meters
	}
	t := TemplateData{Meters: Cfg.Meters}
	err := templates.ExecuteTemplate(w, "meter", t)
	if err != nil {
		log.Println("executing meter template: ", err.Error())
	}
}

func init() {
	flag.Parse()
	if *mkConfigFile {
		log.Println("Creating default configuration in file " + *configFile)
		cfg := defluxio.MkDefaultServerConfiguration()
		err := cfg.Save(*configFile)
		if err != nil {
			log.Fatal("Failed to create default configuration:", err.Error())
		}
		os.Exit(0)
	}
	var loaderror error
	Cfg, loaderror = defluxio.LoadServerConfiguration(*configFile)
	if loaderror != nil {
		log.Fatal("Error loading configuration: ", loaderror.Error())
	}
	log.Printf("Configured meters are:")
	for _, m := range Cfg.Meters {
		log.Printf("* %s (%s)", m.ID, m.Name)
	}
	log.Printf("Starting meter surveilance routine")
	Cfg.Meters.StartBestMeterUpdater(Cfg.MeterTimeout)
	funcMap := template.FuncMap{
		"dosomething": func() string { return "done something" },
		"doublethreedigits": func(f float64) string {
			return fmt.Sprintf("%2.3f", f)
		},
		"tstodate": func(t time.Time) string {
			return fmt.Sprintf("%d.%d.%d", t.Day(), t.Month(), t.Year())
		},
		"tstotime": func(t time.Time) string {
			return fmt.Sprintf("%d:%d:%d", t.Hour(), t.Minute(), t.Second())
		},
	}
	templates, loaderror = template.New("").Funcs(funcMap).ParseGlob(Cfg.Assets.ViewPath +
		"/*")
	if loaderror != nil {
		log.Fatal("Cannot load templates: ", loaderror.Error())
	}
	if Cfg.InfluxDB.Enabled {
		var err error
		dbclient, err = defluxio.NewDBClient(&Cfg.InfluxDB)
		if err != nil {
			log.Fatal("Cannot initialize database client:", err.Error())
		}
		dbchannel = make(chan defluxio.MeterReading)
		pusher, perr := dbclient.MkDBPusher(dbchannel)
		if perr != nil {
			log.Fatal("Cannot generate database pusher: ", perr.Error())
		}
		go pusher()
	}
}

func main() {
	// TODO: Refactor, H is a global right now, this is not good.
	go defluxio.H.Run()
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/impressum", serveImpressum).Methods("GET")
	r.HandleFunc("/meter", serveMeter).Methods("GET")
	r.HandleFunc("/api/submit/{meter}",
		defluxio.MkSubmitReadingHandler(dbchannel, Cfg)).Methods("POST")
	r.HandleFunc("/api/status", defluxio.ServerStatus).Methods("GET")
	r.HandleFunc("/ws", defluxio.ServeWs)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(Cfg.Assets.AssetPath)))
	http.Handle("/", r)
	listenAddress := fmt.Sprintf("%s:%d", Cfg.Network.Host,
		Cfg.Network.Port)
	log.Println("Starting server at " + listenAddress)
	err := http.ListenAndServe(listenAddress, nil)
	if err != nil {
		log.Fatal("Failed to start http server: ", err)
	}
}
