package api

import (
	"encoding/json"
	"locator/config"
	"locator/internal"
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
type Ship struct {
	Name     string   `json:"name"`
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
}

type ShipSplit struct {
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

type ShipsRequest struct {
	Ships []Ship `json:"ships"`
}
type Server interface {
	Router() http.Handler
	PostHelpMe(w http.ResponseWriter, r *http.Request)
	GetHelpMeSplit(w http.ResponseWriter, r *http.Request)
	PostHelpMeSplit(w http.ResponseWriter, r *http.Request)
}

var myShips []Ship
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
	r.PathPrefix("/helpme_split").HandlerFunc(a.PostHelpMeSplit).Methods(http.MethodPost)
	r.PathPrefix("/helpme_split").HandlerFunc(a.GetHelpMeSplit).Methods(http.MethodGet)
	r.HandleFunc("/helpme", a.PostHelpMe).Methods(http.MethodPost)
	a.router = r
	return a
}

//GetHelpMeSplit - GET - /helpme_split
func (a *api) GetHelpMeSplit(w http.ResponseWriter, r *http.Request) {
	var response BasicResponse
	var encontrado bool
	w.Header().Set("Content-Type", "application/json")
	d, m, count := checkShips(myShips)
	if count == len(cfg.Ships) {
		response, encontrado = getResponse(d, m)
		if !encontrado {
			w.WriteHeader(http.StatusNotFound)
			if response.Message == "" {
				response.Message = "Message can't be recovered."
			} else if response.Position.X == -0.09 {
				response.Message = "Location can't be identified."
			}
			response.Message = response.Message + " Check your parameters and try again."
		}
	} else {
		response.Message = "Please check that all configured ships were loaded before trying to get position. Loaded ships cleared. Please try again"
	}

	myShips = nil
	j, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(j)
}

//PostHelpMe - POST - /helpme
func (a *api) PostHelpMe(w http.ResponseWriter, r *http.Request) {
	var response BasicResponse
	var request ShipsRequest
	var encontrado bool = false
	w.Header().Set("Content-Type", "application/json")
	err1 := json.NewDecoder(r.Body).Decode(&request)
	if err1 != nil {
		log.Println(err1)
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "Bad Request Syntax, please check"
		encontrado = true
	} else {
		d, m, count := checkShips(request.Ships)
		if count == len(cfg.Ships) {
			response, encontrado = getResponse(d, m)
		} else {
			response.Message = "Ship name not found in configured names."
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
		if response.Message == "" {
			response.Message = "Message can't be recovered."
		} else if response.Position.X == -0.09 {
			response.Message = "Location can't be identified."
		}
		response.Message = response.Message + " Check your parameters and try again."
	}
}

//PostHelpMeSplit - POST - /helpme_split
func (a *api) PostHelpMeSplit(w http.ResponseWriter, r *http.Request) {
	var response BasicResponse
	var request ShipSplit
	var route []string
	var s Ship

	w.Header().Set("Content-Type", "application/json")
	err1 := json.NewDecoder(r.Body).Decode(&request)

	if err1 != nil {
		log.Println(err1)
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "Wrong request syntax, please check"
	} else {
		route = strings.Split(r.URL.Path, "/")
		s.Name = route[len(route)-1]
		s.Distance = request.Distance
		s.Message = request.Message

		if len(myShips) < len(cfg.Ships) {
			myShips = append(myShips, s)
			response.Message = "Ship " + s.Name + " loaded."
		} else {
			response.Message = "All ships loaded. Run get methot to obtain position and clean array."
		}
	}
	j, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(j)
}

func checkShips(s []Ship) (d []float32, m [][]string, count int) {
	d = make([]float32, len(cfg.Ships))
	m = make([][]string, len(cfg.Ships))
	if len(s) > 0 {
		for i := 0; i < len(s); i++ {
			for j := 0; j < len(cfg.Ships); j++ {
				if strings.EqualFold(s[i].Name, cfg.Ships[j].Name) {
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

	log.Println("Listening...")
	server.ListenAndServe()

}
