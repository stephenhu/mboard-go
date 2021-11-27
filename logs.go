package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"

)

type Log struct {
	Clock       string  `json:"clock"`
	Period      string  `json:"period"`
	Msg					string	`json:"msg"`
	Created			string	`json:"created"`
	Updated			string	`json:"updated"`
}


func logKey(id string) string {
  return fmt.Sprintf("%s.logs", id)
} // logKey


func gameTime(clk *Clock) string {

	m := 0
	s := 0

	if clk.Seconds != 0 {
		m = clk.Seconds/60
		s = clk.Seconds%60
	}

	if m < 10 {

		if s == 60 {
			return fmt.Sprintf("0%d:00.%d", m, clk.Tenths)
		} else if s < 10 {
			return fmt.Sprintf("0%d:0%d.%d", m, s, clk.Tenths)
		} else {
			return fmt.Sprintf("0%d:%d.%d", m, s, clk.Tenths)
		}

	} else {

		if s == 60 {
			return fmt.Sprintf("%d:00.%d", m, clk.Tenths)
		} else if s < 10 {
			return fmt.Sprintf("%d:0%d.%d", m, s, clk.Tenths)
		} else {
			return fmt.Sprintf("%d:%d.%d", m, s, clk.Tenths)
		}

	}

} // gameTime


func put(id string, clk *Clock, req Req) {

	l := Log{
		Clock: gameTime(clk),
		Period: fmt.Sprintf("%d", req.Period),
		Msg: req.Cmd,
	}

  j, err := json.Marshal(l)

	if err != nil {
		log.Println(err)
		return
	} else {

		rp := Red.Get()

		_, err := rp.Do(LPUSH, logKey(id), string(j))

		if err != nil {
			log.Println(err)
		}

	}

} // put


func get(id string) []string {

	rp := Red.Get()

	logs, err := redis.Strings(rp.Do(LRANGE, logKey(id), 0, -1))

	if err != nil {
		log.Println(err)
		return nil
	} else {
		return logs
	}

} // get


func del(log_id string) bool {

	/*
  _, err := data.Exec(
		LogDelete, log_id,
	)

	if err != nil {

		log.Printf("[%s][Error] %s", version(), err)
		return false

	}

	return true
	*/

	return false
} // delete


func logHandler(w http.ResponseWriter, r *http.Request) {

	mux := mux.Vars(r)

  switch r.Method {
  case http.MethodGet:

	  id := mux["id"]

		logs := get(id)

		j, err := json.Marshal(logs)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Write(j)
		}

  case http.MethodPost:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // logHandler
