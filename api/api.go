package api

import (
	"encoding/json"
	"fmt"
	"fuegoquasar/config"
	"fuegoquasar/internal"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Satelite struct {
	Name     string   `json:"name"`
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
}

type SateliteSplit struct {
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
}

type point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type BasicResponse struct {
	Position point  `json:"position"`
	Message  string `json:"message"`
}

type SatelitesRequest struct {
	Satelites []Satelite `json:"satelites"`
}

var misNaves []Satelite

//PostTopSecret - POST - /topsecret
func PostTopSecret(w http.ResponseWriter, r *http.Request) {
	var response BasicResponse
	var request SatelitesRequest
	// Ordenar los satelites segun su nombre. Validar los nombres de los tres, si me falta uno, devolver error.
	err1 := json.NewDecoder(r.Body).Decode(&request)

	if err1 != nil {
		log.Fatal(err1)
	}
	var status1, status2 bool
	w.Header().Set("Content-Type", "application/json")
	response.Position.X, response.Position.Y, status1 = internal.GetLocation(request.Satelites[0].Distance, request.Satelites[1].Distance, request.Satelites[2].Distance)

	response.Message, status2 = internal.GetMessage(request.Satelites[0].Message, request.Satelites[1].Message, request.Satelites[2].Message)

	j, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	if status1 && status2 {
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

//PostTopSecretSplit - POST - /topsecret_split
func PostTopSecretSplit(w http.ResponseWriter, r *http.Request) {
	var response BasicResponse
	var request SateliteSplit

	err1 := json.NewDecoder(r.Body).Decode(&request)

	if err1 != nil {
		log.Fatal(err1)
	}

	w.Header().Set("Content-Type", "application/json")
	var route []string
	var s Satelite
	route = strings.Split(r.URL.Path, "/")
	s.Name = route[len(route)-1]
	s.Distance = request.Distance
	s.Message = request.Message

	fmt.Println(s.Name)

	if len(misNaves) < 3 {
		misNaves = append(misNaves, s)
		response.Message = "Satelite " + s.Name + " cargado."
	} else {
		response.Message = "Todas las naves fueron cargadas. Ejecute un get para obtener la posicion y limpiar el array."
	}

	j, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

//GetTopSecretSplit - GET - /topsecret_split
func GetTopSecretSplit(w http.ResponseWriter, r *http.Request) {
	// Ordenar los satelites segun su nombre. Validar los nombres de los tres, si me falta uno, devolver error.
	var response BasicResponse
	var status1, status2 bool
	w.Header().Set("Content-Type", "application/json")
	if len(misNaves) >= 3 {
		response.Position.X, response.Position.Y, status1 = internal.GetLocation(misNaves[0].Distance, misNaves[1].Distance, misNaves[2].Distance)
		response.Message, status2 = internal.GetMessage(misNaves[0].Message, misNaves[1].Message, misNaves[2].Message)
		misNaves = nil
	} else {
		w.WriteHeader(http.StatusOK)
		status1 = true
		status2 = true
		response.Message = "No se cargaron las tres naves, no se puede calcular la posicion"
	}
	j, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	if status1 && status2 {
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	} else {
		w.WriteHeader(http.StatusNotFound)

	}

}

func Serve() {
	cfg, err := config.NewConfig(config.CfgPath)
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
		log.Printf("Defaulting to port %s", port)
	}

	r := mux.NewRouter().StrictSlash(false)
	r.PathPrefix("/topsecret_split").HandlerFunc(PostTopSecretSplit).Methods("POST")
	r.PathPrefix("/topsecret_split").HandlerFunc(GetTopSecretSplit).Methods("GET")
	r.HandleFunc("/topsecret", PostTopSecret).Methods("POST")

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    cfg.Server.Timeout.Read * time.Second,
		WriteTimeout:   cfg.Server.Timeout.Write * time.Second,
		IdleTimeout:    cfg.Server.Timeout.Idle * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Escuchando...")
	server.ListenAndServe()

}
