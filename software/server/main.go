// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
)

var configFile = flag.String("config", "defluxio.conf", "configuration file")
var templates *template.Template

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
	loadConfiguration(*configFile)
	templates = template.Must(template.ParseGlob(Cfg.Assets.ViewPath + "/*"))
	if Cfg.InfluxDB.Enabled {
		InitDBConnector()
	}
}

func main() {
	go h.run()
	if Cfg.InfluxDB.Enabled {
		go DBPusher()
	}
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/impressum", serveImpressum).Methods("GET")
	r.HandleFunc("/api/submit/{meter}", submitReading).Methods("POST")
	r.HandleFunc("/api/status", serverStatus).Methods("GET")
	r.HandleFunc("/ws", serveWs)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(Cfg.Assets.AssetPath)))
	http.Handle("/", r)
	listenAddress := fmt.Sprintf("%s:%d", Cfg.Network.Host, Cfg.Network.Port)
	log.Println("Starting server at " + listenAddress)
	err := http.ListenAndServe(listenAddress, nil)
	if err != nil {
		log.Fatal("Failed to start http server: ", err)
	}
}
