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

func makeShipsRequest() (ships []Ship) {
	x, y := getRandomPoint()
	for i := 0; i < len(cfg.Ships); i++ {
		var ship Ship
		ship.Message = getRandomMessage()
		ship.Distance = getDistance(cfg.Ships[i].X, cfg.Ships[i].Y, x, y)
		ship.Name = cfg.Ships[i].Name
		ships = append(ships, ship)

	}
	return ships
}

func makeShipSplitRequests() (ships []ShipSplit, names []string) {
	x, y := getRandomPoint()
	for i := 0; i < len(cfg.Ships); i++ {
		var ship ShipSplit
		ship.Message = getRandomMessage()
		ship.Distance = getDistance(cfg.Ships[i].X, cfg.Ships[i].Y, x, y)
		ships = append(ships, ship)
		names = append(names, cfg.Ships[i].Name)
	}
	return ships, names
}

func getDistance(a, b, c, d float32) (distance float32) {
	dx1 := a - c
	dy1 := b - d
	distance = float32(math.Sqrt(float64(dx1*dx1 + dy1*dy1)))

	return distance
}

func TestPostHelpMeSplit(t *testing.T) {
	server := New()
	ts := httptest.NewServer(server.Router())
	defer ts.Close()

	ships, names := makeShipSplitRequests()
	for i := 0; i < len(ships); i++ {
		bdy, err1 := json.Marshal(ships[i])
		if err1 != nil {
			t.Errorf("Can't create body.")
		}
		req, err := http.NewRequest("POST", "/helpme_split/"+names[i], bytes.NewBuffer(bdy))
		if err != nil {
			t.Errorf("Can't create request %s", err.Error())
		}
		rec := httptest.NewRecorder()
		server.PostHelpMeSplit(rec, req)
		res := rec.Result()
		defer res.Body.Close()
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
			t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
			t.Log(ships[i])
		} else {
			b, err2 := ioutil.ReadAll(res.Body)
			if err2 != nil {
				t.Fatalf("Response cannot be read: %v", err2)
			}
			var resp BasicResponse
			err1 := json.Unmarshal(b, &resp)
			if err1 != nil {
				t.Fatalf("Response Json can't be parsed: %v", err1)
			}
			t.Log(resp)
		}
	}

}

func TestGetHelpMeSplit(t *testing.T) {
	t.Log(myShips)
	server := New()
	testserver := httptest.NewServer(server.Router())
	defer testserver.Close()
	req, err := http.NewRequest("GET", "/helpme_split", nil)
	if err != nil {
		t.Errorf("Can't create request %s", err.Error())
	}
	rec := httptest.NewRecorder()
	server.GetHelpMeSplit(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	} else {
		b, err2 := ioutil.ReadAll(res.Body)
		if err2 != nil {
			t.Fatalf("Response can't be read: %v", err2)
		}
		var resp BasicResponse
		err1 := json.Unmarshal(b, &resp)
		if err1 != nil {
			t.Fatalf("Response Json can't be parsed: %v", err1)
		}
		t.Log(resp)

	}
}

func TestPostHelpMe(t *testing.T) {
	server := New()
	testserver := httptest.NewServer(server.Router())
	defer testserver.Close()
	var shipsRequest ShipsRequest

	shipsRequest.Ships = makeShipsRequest()
	bdy, err1 := json.Marshal(shipsRequest)
	if err1 != nil {
		t.Errorf("Can't create request.")
	}
	req, err := http.NewRequest("POST", "/helpme", bytes.NewBuffer(bdy))
	rec := httptest.NewRecorder()
	server.PostHelpMe(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
		t.Log(shipsRequest)
	} else {
		if res.StatusCode != http.StatusNotFound {
			b, err2 := ioutil.ReadAll(res.Body)
			if err2 != nil {
				t.Fatalf("Response can't be read: %v", err2)
			}
			var resp BasicResponse
			err1 := json.Unmarshal(b, &resp)
			if err1 != nil {
				t.Fatalf("Response Json can't be parsed: %v", err1)
			}
			t.Log(resp)
		} else {
			t.Log("Failed to get message and/or position")
			t.Log(shipsRequest)
		}
		shipsRequest.Ships = nil
	}
}
