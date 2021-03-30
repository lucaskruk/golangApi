package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var s Server

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getRandomMessage() (m []string) {
	arraylength := rand.Intn(7) + 1

	for i := 0; i < arraylength; i++ {
		stringlength := rand.Intn(8)
		message := RandStringRunes(stringlength)
		m = append(m, message)
	}

	return
}

func getRandomPoint() (x, y float32) {

	min := -6000
	max := 6000
	val := max - min + 1
	x = (float32(rand.Intn(val)) + float32(min)) / 10
	y = (float32(rand.Intn(val)) + float32(min)) / 10
	return
}

func makeSatelitesRequest() (s []Satelite) {
	x, y := getRandomPoint()
	for i := 0; i < len(cfg.RebelShips); i++ {
		var sat Satelite
		sat.Message = getRandomMessage()
		sat.Distance = getDistance(cfg.RebelShips[i].X, cfg.RebelShips[i].Y, x, y)
		sat.Name = cfg.RebelShips[i].Name
		s = append(s, sat)

	}
	return s
}

func makeSatSplitRequests() (s []SateliteSplit, names []string) {
	x, y := getRandomPoint()
	for i := 0; i < len(cfg.RebelShips); i++ {
		var sat SateliteSplit
		sat.Message = getRandomMessage()
		sat.Distance = getDistance(cfg.RebelShips[i].X, cfg.RebelShips[i].Y, x, y)
		s = append(s, sat)
		names = append(names, cfg.RebelShips[i].Name)
	}
	return s, names
}

func getDistance(a, b, c, d float32) (distance float32) {
	dx1 := a - c
	dy1 := b - d
	distance = float32(math.Sqrt(float64(dx1*dx1 + dy1*dy1)))

	return distance
}

func TestPostTopSecretSplit(t *testing.T) {
	s := New()
	ts := httptest.NewServer(s.Router())
	defer ts.Close()

	sats, names := makeSatSplitRequests()
	for i := 0; i < len(sats); i++ {
		bdy, err1 := json.Marshal(sats[i])
		if err1 != nil {
			t.Errorf("No armo el body.")
		}
		req, err := http.NewRequest("POST", "/topsecret_split/"+names[i], bytes.NewBuffer(bdy))
		if err != nil {
			t.Errorf("No pude armar request %s", err.Error())
		}
		rec := httptest.NewRecorder()
		s.PostTopSecretSplit(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
			t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
			t.Log(sats[i])
		} else {
			b, err2 := ioutil.ReadAll(res.Body)
			if err2 != nil {
				t.Fatalf("No se puede leer response: %v", err2)
			}
			var resp BasicResponse
			err1 := json.Unmarshal(b, &resp)
			if err1 != nil {
				t.Fatalf("No se puede parsear el json del response: %v", err1)
			}
			t.Log(resp)
		}
	}

}

func TestGetTopSecretSplit(t *testing.T) {
	s := New()
	ts := httptest.NewServer(s.Router())
	defer ts.Close()
	req, err := http.NewRequest("GET", "/topsecret_split", nil)
	if err != nil {
		t.Errorf("No pude armar request %s", err.Error())
	}
	rec := httptest.NewRecorder()
	s.GetTopSecretSplit(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	} else {
		b, err2 := ioutil.ReadAll(res.Body)
		if err2 != nil {
			t.Fatalf("No se puede leer response: %v", err2)
		}
		var resp BasicResponse
		err1 := json.Unmarshal(b, &resp)
		if err1 != nil {
			t.Fatalf("No se puede parsear el json del response: %v", err1)
		}
		t.Log(resp)
	}
}

func TestPostTopSecret(t *testing.T) {
	s := New()
	ts := httptest.NewServer(s.Router())
	defer ts.Close()
	var satRequest SatelitesRequest

	satRequest.Satelites = makeSatelitesRequest()
	bdy, err1 := json.Marshal(satRequest)
	if err1 != nil {
		t.Errorf("No se armo correctamente el request.")
	}
	req, err := http.NewRequest("POST", "/topsecret", bytes.NewBuffer(bdy))
	rec := httptest.NewRecorder()
	s.PostTopSecret(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
		t.Log(satRequest)
	} else {
		satRequest.Satelites = nil

		b, err2 := ioutil.ReadAll(res.Body)
		if err2 != nil {
			t.Fatalf("No se puede leer response: %v", err2)
		}
		var resp BasicResponse
		err1 := json.Unmarshal(b, &resp)
		if err1 != nil {
			t.Fatalf("No se puede parsear el json del response: %v", err1)
		}
		t.Log(resp)
	}
}
