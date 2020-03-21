package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type SubscriberMapResponse struct {
	Page 			string							`json:"page"`
	Options   map[string]string 	`json:"options"`
}

type SubscriberStringResponse struct {
  Key 			string				`json:"key"`
	Val				string				`json:"val"`
}

type SubscriberStateResponse struct {
	Key				string				`json:"key"`
	State     *GameState    `json:"state"`
}

var subscribers map[*websocket.Conn] *sync.Mutex

func sendToSubscribers(j []byte) {

	for c, mu := range subscribers {
		
		mu.Lock()
		c.WriteMessage(websocket.TextMessage, j)
		mu.Unlock()

	}

} // sendToSubscribers

func pushState(state *GameState) {

	n := SubscriberStateResponse{
		Key: WS_GAME_STATE,
		State: state,
	}

	j, jsonErr := json.Marshal(n)

	if jsonErr != nil {
		log.Println(jsonErr)
	}

	sendToSubscribers(j)

} // pushState

func pushString(key string, val string) {

	n := SubscriberStringResponse{
		Key: key,
		Val: val,
	}

	j, jsonErr := json.Marshal(n)

	if jsonErr != nil {
		log.Println(jsonErr)
	}

	sendToSubscribers(j)

} // pushString

func pushMap(msg string, options map[string] string) {

	r := SubscriberMapResponse{
		Page: msg,
		Options: options,
	}

	j, err := json.Marshal(r)

	if err != nil {
		log.Println(err)
		return
	}

	sendToSubscribers(j)

} // pushMap

func subscriberHandler(w http.ResponseWriter, r *http.Request) {

  upgrader := websocket.Upgrader {
		ReadBufferSize:		1024,
		WriteBufferSize: 	1024,
		CheckOrigin:		func(r *http.Request) bool { return true },
	}

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		
		log.Println("[Error]", err)
		return

	}

	defer c.Close()

	if subscribers == nil {
		subscribers = make(map[*websocket.Conn]*sync.Mutex)
	}

	subscribers[c] = &sync.Mutex{}

	for {
   
		_, msg, err := c.ReadMessage()

		if err != nil {

			log.Println("[Error] ", err)

			if websocket.IsUnexpectedCloseError(err) {

				//log.Println("Removing connection: ", c)
				delete(subscribers, c)

			}

			break

		}

    if msg == nil {
			log.Println(msg)
			break
		}		

	}
	
} // subscriberHandler
