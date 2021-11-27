package main

import (
	"encoding/json"
	//"fmt"
  "log"
	"net/http"
	//"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

type GameData struct {
	Home				*Team			`json:"home"`
	Away				*Team			`json:"away"`
	ShotClock		Clock     `json:"shotClock"`
	GameClock		Clock     `json:"gameClock"`
	Period      int       `json:"period"`
	Possession  bool      `json:"possession"`
	Status      int       `json:"status"`
}

type Plays struct {
	Logs 				[]Log			`json:"logs"`
}

type GameRecord struct {
	Config      GameConfig			`json:"config"`
	Data        GameData        `json:"data"`
	Plays           `json:"plays"`
}


func addGame(config string, ts int64) {

	rp := Red.Get()

	err := rp.Send(HSET, ts, KEY_FIELD_CONFIG, config)

	if err != nil {
		log.Println(err)
	} else {

		rp.Flush()
		rp.Close()

	}

} // addGame


func updateGameData(id string, gd GameData) {

	rp := Red.Get()

	j, err := json.Marshal(gd)

	if err != nil {
		log.Println(err)
	} else {

		_, err := rp.Do(HSET, id, KEY_FIELD_DATA, j)

		if err != nil {
			log.Println(err)
		}

	}

} // updateGameData


func getGameData(id string) []string {

	rp := Red.Get()

	d, err := redis.Strings(rp.Do(HGET, id, KEY_FIELD_DATA))

	if err != nil {
		log.Println(err)
		return nil
	} else {
		return d
	}

} // getGameData


func getGames() []GameRecord {
	return nil
} // getGames


func getGameRecord(id string) *GameRecord {

	if id == "" {
		return nil
	}

	rp := Red.Get()

	all, err := redis.StringMap(rp.Do(HGETALL, id))

	if err != nil {
		log.Println(err)
	} else {

		gr := GameRecord{}

		for k, v := range all {

			switch k {
			case KEY_FIELD_CONFIG:

				conf := GameConfig{}

				err := json.Unmarshal([]byte(v), &conf)

				if err != nil {
					log.Println(err)
				} else {
					gr.Config = conf
				}

			case KEY_FIELD_DATA:

				gd := GameData{}

				err := json.Unmarshal([]byte(v), &gd)

				if err != nil {
					log.Println(err)
				} else {
					gr.Data = gd
				}

			case KEY_FIELD_PLAYS:

				p := Plays{}

				err := json.Unmarshal([]byte(v), &p)

				if err != nil {
					log.Println(err)
				} else {
					gr.Plays = p
				}

			}

		}

		return &gr

	}

	return nil

} // getGameRecord


func deleteGame(id string) {

} // deleteGame


func scoreHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
  case http.MethodPost:
	case http.MethodGet:

		vars := mux.Vars(r)

		id := vars["id"]

		if id != "" {

			g, ok := gameMap[id]

			if !ok {
				w.WriteHeader(http.StatusNotFound)
			} else {

				j, jsonErr := json.Marshal(g)

				if jsonErr != nil {
					log.Printf("[%s] %s", version(), jsonErr)
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					w.Write(j)
				}

			}

		} else {

			gts := getGames()

			j, jsonErr := json.Marshal(gts)

			if jsonErr != nil {
				log.Printf("[%s] %s", version(), jsonErr)
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.Write(j)
			}

		}


	case http.MethodPut:
	case http.MethodDelete:

	  vars := mux.Vars(r)

		id := vars["id"]

		if id != "" {

			deleteGame(id)

			w.WriteHeader(http.StatusOK)

		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	default:
		log.Printf("[%s][Error] unsupported command", version())
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // scoreHandler
