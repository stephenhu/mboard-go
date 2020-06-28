package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
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

var subscribersMap map[string]map[*websocket.Conn] *sync.Mutex


func sendToSubscribers(id string, j []byte) {

	s, ok := subscribersMap[id]

	if ok {

		for c, mu := range s {

			mu.Lock()
			c.WriteMessage(websocket.TextMessage, j)
			mu.Unlock()

		}

	}

} // sendToSubscribers


func pushState(id string, state *GameState) {

	n := SubscriberStateResponse{
		Key: WS_GAME_STATE,
		State: state,
	}

	j, err := json.Marshal(n)

	if err != nil {
		log.Println(err)
	}

	sendToSubscribers(id,  j)

} // pushState


func pushString(id string, key string, val string) {

	n := SubscriberStringResponse{
		Key: key,
		Val: val,
	}

	j, err := json.Marshal(n)

	if err != nil {
		log.Println(err)
	}

	sendToSubscribers(id, j)

} // pushString


func pushMap(id string, msg string, options map[string] string) {

	r := SubscriberMapResponse{
		Page: msg,
		Options: options,
	}

	j, err := json.Marshal(r)

	if err != nil {
		log.Println(err)
		return
	}

	sendToSubscribers(id, j)

} // pushMap


func cleanupSubscribers(id string) {

	s, ok := subscribersMap[id]

	if ok {

		for conn, _ := range s {
			conn.Close()
		}

	}

} // cleanupSubscribers


func subscriberHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id := vars["id"]

	g, ok := gameMap[id]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {

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

		if subscribersMap == nil {
			subscribersMap = make(map[string]map[*websocket.Conn]*sync.Mutex)
		}

		_, ok := subscribersMap[id]

		if !ok {
			subscribersMap[id] = make(map[*websocket.Conn]*sync.Mutex)
		}

		subscribersMap[id][c] = &sync.Mutex{}

		gs := GameState{
			Settings: g.Game.Settings,
			Period: g.Game.GameData.Period,
			Possession: g.Game.GameData.Possession,
			Home: g.Game.GameData.Home,
			Away: g.Game.GameData.Away,
			GameClock: g.Game.GameData.Clk.PlayClock,
			ShotClock: g.Game.GameData.Clk.ShotClock,
		}

		pushState(id, &gs)

		for {

			_, msg, err := c.ReadMessage()

			if err != nil {

				log.Println("[Error] ", err)

				if websocket.IsUnexpectedCloseError(err) {
					delete(subscribersMap[id], c)
				}

				break

			}

			if msg == nil {
				continue
			}

		}

	}

} // subscriberHandler
