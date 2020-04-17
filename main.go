/**
 * Copyright (C) 2020, State Of The Art https://www.state-of-the-art.io/
 */

package main

import (
    "fmt"
    "log"
    "time"
    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/gorilla/mux"     // http api

    "nyansync/config"
)

type Catalog struct {
    Movies []CatalogItem
    Music []CatalogItem
}

type CatalogItem struct {
    Url string
    Name string
    Type string
    Image string
}

type Access []struct {
    Url string
    AccessTokenHash string
    StampFrom time.Time
    StampTo time.Time
    HitLimit int
}

type event struct {
    ID          string `json:"ID"`
    Title       string `json:"Title"`
    Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
    {
        ID:          "1",
        Title:       "Introduction to Golang",
        Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
    },
}

func homeLink(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome home!")
}

func pingLink(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong")
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

func createEvent(w http.ResponseWriter, r *http.Request) {
    var newEvent event
    reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
    }

    json.Unmarshal(reqBody, &newEvent)
    events = append(events, newEvent)
    w.WriteHeader(http.StatusCreated)

    json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
    eventID := mux.Vars(r)["id"]

    for _, singleEvent := range events {
        if singleEvent.ID == eventID {
            json.NewEncoder(w).Encode(singleEvent)
        }
    }
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
    eventID := mux.Vars(r)["id"]
    var updatedEvent event

    reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
    }
    json.Unmarshal(reqBody, &updatedEvent)

    for i, singleEvent := range events {
        if singleEvent.ID == eventID {
            singleEvent.Title = updatedEvent.Title
            singleEvent.Description = updatedEvent.Description
            events = append(events[:i], singleEvent)
            json.NewEncoder(w).Encode(singleEvent)
        }
    }
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
    eventID := mux.Vars(r)["id"]

    for i, singleEvent := range events {
        if singleEvent.ID == eventID {
            events = append(events[:i], events[i+1:]...)
            fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
        }
    }
}

func initApiV1(router *mux.Router) {
    api := router.PathPrefix("/api/v1").Subrouter()

    api.HandleFunc("", apiV1Get).Methods(http.MethodGet)
    api.HandleFunc("", apiV1Post).Methods(http.MethodPost)
    api.HandleFunc("", apiV1Put).Methods(http.MethodPut)
    api.HandleFunc("", apiV1Delete).Methods(http.MethodDelete)

    api.HandleFunc("/event", createEvent).Methods(http.MethodPost)
    api.HandleFunc("/events", getAllEvents).Methods(http.MethodGet)
    api.HandleFunc("/events/{id}", getOneEvent).Methods(http.MethodGet)
    api.HandleFunc("/events/{id}", updateEvent).Methods(http.MethodPut)
    api.HandleFunc("/events/{id}", deleteEvent).Methods(http.MethodDelete)
}

func main() {
    cfg := config.Load()
    fmt.Printf("%+v\n", cfg)

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", homeLink)
    router.HandleFunc("/ping", pingLink)

    initApiV1(router)

    log.Fatal(http.ListenAndServe(":8080", router))
}
