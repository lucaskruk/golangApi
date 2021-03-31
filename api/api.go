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

type api struct {
	router http.Handler
}
type Satelite struct {
	Name     string   `json:"name"`
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
}

type SateliteSplit struct {
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
}

type BasicResponse struct {
	Position struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
	} `json:"position"`
	Message string `json:"message"`
}

type SatelitesRequest struct {
	Satelites []Satelite `json:"satelites"`
}
type Server interface {
	Router() http.Handler
	PostTopSecret(w http.ResponseWriter, r *http.Request)
	GetTopSecretSplit(w http.ResponseWriter, r *http.Request)
	PostTopSecretSplit(w http.ResponseWriter, r *http.Request)
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
func (a *api) Router() http.Handler {
	return a.router
}

func New() Server {
	a := &api{}
	r := mux.NewRouter()
	r.PathPrefix("/topsecret_split").HandlerFunc(a.PostTopSecretSplit).Methods(http.MethodPost)
	r.PathPrefix("/topsecret_split").HandlerFunc(a.GetTopSecretSplit).Methods(http.MethodGet)
	r.HandleFunc("/topsecret", a.PostTopSecret).Methods(http.MethodPost)
	a.router = r
	return a
}

//GetTopSecretSplit - GET - /topsecret_split
func (a *api) GetTopSecretSplit(w http.ResponseWriter, r *http.Request) {
	var response BasicResponse
	var encontrado bool
	w.Header().Set("Content-Type", "application/json")
	d, m, count := validaNaves(misNaves)
	if count == len(cfg.RebelShips) {
		response, encontrado = getResponse(d, m)
		if !encontrado {
			w.WriteHeader(http.StatusNotFound)
			if response.Message == "" {
				response.Message = "No se logro descifrar el mensaje."
			} else if response.Position.X == -0.09 {
				response.Message = "No se identifica la ubicacion."
			}
			response.Message = response.Message + " Revise sus parametros e intente nuevamente"
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		response.Message = "No se cargaron los tres satelites, o los nombres no coinciden. Vuelva a cargarlos correctamente"
	}

	misNaves = nil
	j, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(j)
}

//PostTopSecret - POST - /topsecret
func (a *api) PostTopSecret(w http.ResponseWriter, r *http.Request) {
	var response BasicResponse
	var request SatelitesRequest
	var encontrado bool = false
	w.Header().Set("Content-Type", "application/json")
	err1 := json.NewDecoder(r.Body).Decode(&request)
	if err1 != nil {
		log.Println(err1)
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "Error en el request recibido, por favor verificar"
		encontrado = true
	} else {
		d, m, count := validaNaves(request.Satelites)
		if count == len(cfg.RebelShips) { // confirmo que tengo todas las naves configuradas
			response, encontrado = getResponse(d, m)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Las naves no fueron cargadas correctamente. Verifique los nombres")
		}
	}
	j, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	if encontrado {
		w.Write(j)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

//PostTopSecretSplit - POST - /topsecret_split
func (a *api) PostTopSecretSplit(w http.ResponseWriter, r *http.Request) {
	var response BasicResponse
	var request SateliteSplit
	var route []string
	var s Satelite

	w.Header().Set("Content-Type", "application/json")
	err1 := json.NewDecoder(r.Body).Decode(&request)

	if err1 != nil {
		log.Println(err1)
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "Error en el request recibido, por favor verificar"
	} else {
		route = strings.Split(r.URL.Path, "/")
		s.Name = route[len(route)-1]
		s.Distance = request.Distance
		s.Message = request.Message

		if len(misNaves) < len(cfg.RebelShips) {
			misNaves = append(misNaves, s)
			response.Message = "Satelite " + s.Name + " cargado."
		} else {
			response.Message = "Todas las naves fueron cargadas. Ejecute un get para obtener la posicion y limpiar el array."
		}
	}
	j, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(j)
}

func validaNaves(s []Satelite) (d []float32, m [][]string, count int) {
	// Ordena los satelites segun su nombre, y devuelve la cantidad
	// El orden es: Kenobi, SkyWalker, Sato, segun definido en la configuracion
	d = make([]float32, len(cfg.RebelShips))
	m = make([][]string, len(cfg.RebelShips))
	if len(s) > 0 {
		for i := 0; i < len(s); i++ {
			for j := 0; j < len(cfg.RebelShips); j++ {
				if strings.ToLower(s[i].Name) == strings.ToLower(cfg.RebelShips[j].Name) {
					d[j] = s[i].Distance
					m[j] = s[i].Message
					count++
				}
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

func Serve() {

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
		log.Printf("Defaulting to port %s", port)
	}

	s := New()

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        s.Router(),
		ReadTimeout:    cfg.Server.Timeout.Read * time.Second,
		WriteTimeout:   cfg.Server.Timeout.Write * time.Second,
		IdleTimeout:    cfg.Server.Timeout.Idle * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Escuchando...")
	server.ListenAndServe()

}
