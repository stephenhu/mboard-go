package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var invalidURL = []string{
	"127.0.0.1:8000/games",
	"ws://127.0.0.1:8000/games",
	"ws://127.0.0.1:8000/games/",
	"ws://127.0.0.1:8000/api/games/",
	"ws://127.0.0.1:8000/ws",
	"http://127.0.0.1:8000/ws/games",
	"ws://127.0.0.1:8000/ws/games/g",
}

var validURL = []string {
	"ws://127.0.0.1:8000/ws/games/a",
	"ws://127.0.0.1:8000/ws/games/0",
}

var gameId string

func parseBody(body io.ReadCloser, s interface {}) {

  j, readErr := ioutil.ReadAll(body)

  if readErr != nil {
		log.Println(readErr)
	}

  json.Unmarshal(j, &s)

} // parseBody

func TestInvalidConnect(t *testing.T) {

  for _, u := range invalidURL {

    _, _, err := websocket.DefaultDialer.Dial(u, nil)

		if err == nil {
			t.Fatal("Should not connect successfully.")
		}
		
	}

} // TestInvalidConnect

func TestConnect(t *testing.T) {

  for _, u := range validURL {

    ws, _, err := websocket.DefaultDialer.Dial(u, nil)

		if err != nil {
			t.Fatal(err)
		}

		ws.Close()
	
	}

} // TestConnect

func TestNewGame(t *testing.T) {

  form := url.Values{}
	form.Add("periods", "19")
	form.Add("minutes", "12")

  r, postErr := http.Post("http://127.0.0.1:8000/api/games",
	  "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	
	if postErr != nil {
		t.Fatal(postErr)
	}

	defer r.Body.Close()

  b := PostRes{}

  parseBody(r.Body, &b)
	
	if b.GameId == "" {
		t.Fatal("No GameId returned.")
	}

	gameId = b.GameId

  url := fmt.Sprintf("%s%s",
	  "http://127.0.0.1:8000/api/games/", b.GameId)
	
	r2, getErr := http.Get(url)

	if getErr != nil {
		t.Fatal(getErr)
	}

  config := Config{}

	parseBody(r2.Body, &config)

  if config.Periods != 4 {
		t.Fatal("Returned incorrect periods")
	}

	if config.Minutes != 12 {
		t.Fatal("Returned incorrect minutes")
	}
	
} // TestNewGame


func TestScore(t *testing.T) {

  url := fmt.Sprintf("ws://127.0.0.1:8000/ws/games/%s", gameId)

  ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatal(err)
	}
	
	req := Req{
		Cmd: "SCORE_UP",
		Data: &ReqData{
			Step: 2,
			Value: HOME,
		},
	}

  j, jsonErr := json.Marshal(req)

	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	writeErr := ws.WriteMessage(websocket.TextMessage, j)

	if writeErr != nil {
		t.Fatal(writeErr)
	}
	
	req2 := Req{
		Cmd: "SCORE_DOWN",
		Data: &ReqData{
			Step: -2,
			Value: HOME,
		},
	}

  j2, jsonErr2 := json.Marshal(req2)

	if jsonErr2 != nil {
		t.Fatal(jsonErr2)
	}

	writeErr2 := ws.WriteMessage(websocket.TextMessage, j2)

	if writeErr2 != nil {
		t.Fatal(writeErr2)
	}

	ws.Close()

} // TestScore

func TestClock(t *testing.T) {

  url := fmt.Sprintf("ws://127.0.0.1:8000/ws/games/%s", gameId)

  ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatal(err)
	}
	
	req := Req{
		Cmd: "CLOCK_START",
	}

  j, jsonErr := json.Marshal(req)

	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	writeErr := ws.WriteMessage(websocket.TextMessage, j)

	if writeErr != nil {
		t.Fatal(writeErr)
	}
	
	time.Sleep(2 * time.Second)

	req2 := Req{
		Cmd: "CLOCK_STOP",
	}

  j2, jsonErr2 := json.Marshal(req2)

	if jsonErr2 != nil {
		t.Fatal(jsonErr2)
	}

	writeErr2 := ws.WriteMessage(websocket.TextMessage, j2)

	if writeErr2 != nil {
		t.Fatal(writeErr2)
	}

	ws.Close()

} // TestClock
