// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/gonium/defluxio"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
)

var configFile = flag.String("config", "defluxio.conf", "configuration file")
var Cfg *defluxio.ServerConfigurationData
var templates *template.Template
var dbclient *defluxio.DBClient

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "index", r.Host)
}

func serveImpressum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "impressum", r.Host)
}

func init() {
	flag.Parse()
	var loaderror error
	Cfg, loaderror = defluxio.LoadServerConfiguration(*configFile)
	if loaderror != nil {
		log.Fatal("Error loading configuration: ", loaderror.Error())
	}
	templates = template.Must(template.ParseGlob(Cfg.Assets.ViewPath + "/*"))
	if Cfg.InfluxDB.Enabled {
		var err error
		dbclient, err = defluxio.NewDBClient(Cfg)
		if err != nil {
			log.Fatal("Cannot initialize database client:", err.Error)
		}
		go dbclient.MkDBPusher()
	}
}

func main() {
	// TODO: Refactor, H is a global right now, this is not good.
	go defluxio.H.Run()
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/impressum", serveImpressum).Methods("GET")
	r.HandleFunc("/api/submit/{meter}",
		defluxio.MkSubmitReadingHandler(dbclient, Cfg)).Methods("POST")
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
