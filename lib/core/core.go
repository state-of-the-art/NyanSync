package core

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux" // http api

	"github.com/state-of-the-art/NyanSync/lib/config"
	"github.com/state-of-the-art/NyanSync/lib/generated"
	"github.com/state-of-the-art/NyanSync/lib/state"
)

func Init(configuration *config.Config) {
	// Init cfg variable
	cfg = configuration

	// Init state
	state.Init(cfg.StateFilePath)
}

func pingLink(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func apiV1Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}

func apiV1Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

func apiV1Put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "put called"}`))
}

func apiV1Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "delete called"}`))
}

func initApiV1(router *mux.Router) {
	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("", apiV1Get).Methods(http.MethodGet)
	api.HandleFunc("", apiV1Post).Methods(http.MethodPost)
	api.HandleFunc("", apiV1Put).Methods(http.MethodPut)
	api.HandleFunc("", apiV1Delete).Methods(http.MethodDelete)
}

// TODO: replace with actual function
func getJSMetadata(w http.ResponseWriter, r *http.Request) {
    meta, _ := json.Marshal(map[string]string{
        "deviceID": "asdasdijqwlqpowkepqwe", // TODO: Replace with actual generated device id
    })
    w.Header().Set("Content-Type", "application/javascript")
    fmt.Fprintf(w, "var metadata = %s;\n", meta)
}

func RunHTTPServer() {
	//router := mux.NewRouter().StrictSlash(true)
	//router.Handle("/", http.FileServer(generated.Gui()))
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(generated.Gui(cfg.GuiPath)))
	router.HandleFunc("/meta.js", getJSMetadata)

	//router.Handle("/", http.FileServer(http.Dir(".")))
	/*router.HandleFunc("/ping", pingLink)

	initApiV1(router)*/

	fmt.Printf("Run HTTP server on: %s\n", cfg.Endpoint.Address)
	log.Fatal(http.ListenAndServe(cfg.Endpoint.Address, router))
}

// Core configuration
var cfg = &config.Config{}
