package api

import (
	"encoding/json"
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
var cfg *config.Config
var err error

func init() {
	cfg, err = config.NewConfig(config.CfgPath)
	if err != nil {
		log.Fatal(err)
	}
}

func validaNaves(s []Satelite) (d []float32, m [][]string, count int) {
	// Ordena los satelites segun su nombre.
	// El orden es: Kenobi, SkyWalker, Sato, segun definido en la configuracion
	m = make([][]string, 3)
	d = make([]float32, 3)
	if len(s) > 0 {
		for i := 0; i < 3; i++ {
			if strings.ToLower(s[i].Name) == strings.ToLower(cfg.RebelShip1.Name) {
				d[0] = s[i].Distance
				m[0] = s[i].Message
				count++
			}
			if strings.ToLower(s[i].Name) == strings.ToLower(cfg.RebelShip2.Name) {
				d[1] = s[i].Distance
				m[1] = s[i].Message
				count++
			}
			if strings.ToLower(s[i].Name) == strings.ToLower(cfg.RebelShip3.Name) {
				d[2] = s[i].Distance
				m[2] = s[i].Message
				count++
			}
		}
	}
	return
}

func getResponse(d []float32, m [][]string) (resp BasicResponse, status bool) {

	resp.Position.X, resp.Position.Y = internal.GetLocation(d...)
	resp.Message = internal.GetMessage(m...)

	if resp.Position.X != -0.09 && resp.Message != "" {
		status = true
	}
	return resp, status
}

//PostTopSecret - POST - /topsecret
func PostTopSecret(w http.ResponseWriter, r *http.Request) {
	var response BasicResponse
	var request SatelitesRequest
	var status bool = false

	err1 := json.NewDecoder(r.Body).Decode(&request)

	d, m, count := validaNaves(request.Satelites)

	if err1 != nil {
		log.Fatal(err1)
	}

	w.Header().Set("Content-Type", "application/json")
	if count == 3 {
		response, status = getResponse(d, m)
	} else {
		log.Println("Las naves no fueron cargadas correctamente. Verifique los nombres")
	}
	j, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	if status {
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
	var status bool
	w.Header().Set("Content-Type", "application/json")
	d, m, count := validaNaves(misNaves)
	if count == 3 {
		response, status = getResponse(d, m)
	} else {
		status = true
		response.Message = "No se cargaron las tres naves, o los nombres no coinciden, no se puede calcular la posicion. Vuelva a cargarlas correctamente"
	}
	misNaves = nil
	j, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	if status {
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	} else {
		w.WriteHeader(http.StatusNotFound)

	}

}

func Serve() {

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
