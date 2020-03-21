package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Log struct {
	ID          string  `json:"id"`
	Clock       string  `json:"clock"`
	Msg					string	`json:"msg"`
	Created			string	`json:"created"`
	Updated			string	`json:"updated"`
}

const (

	LogCreate = "INSERT into logs" +
	  "(game_id, clock, msg) " + 
		"VALUES ($1, $2, $3)"
	
	LogGet = "SELECT " +
	  "id, clock, msg, created, updated " +
		"FROM logs " + 
		"WHERE game_id=? ORDER BY created DESC"

  LogDelete = "DELETE from logs WHERE id=?"

)

func gameTime(clk *Clock) string {

	m := 0
	s := 0

	if clk.Seconds != 0 {
		m = clk.Seconds/60
		s = clk.Seconds%60
		log.Println(m)
		log.Println(s)
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

func put(game_id string, clk *Clock, req Req) {

  j, errJson := json.Marshal(req)

	if errJson != nil {
		log.Println(errJson)
		return
	}

  _, err := data.Exec(
		LogCreate, game_id, gameTime(clk), j,
	)

	if err != nil {
		
		log.Printf("[%s][Error] %s", version(), err)
		return

	}

} // put

func get(game_id string) []Log {

	rows, err := data.Query(
		LogGet, game_id,
	)

	if err != nil {

		log.Printf("[%s][Error] %s", version(), err)
		return nil

	}

	defer rows.Close()

	logs := []Log{}

	for rows.Next() {

		l := Log{}

		err := rows.Scan(&l.ID, &l.Clock, &l.Msg, &l.Created, &l.Updated)

		if err == sql.ErrNoRows || err != nil {
			
			log.Printf("[%s][Error] %s", version(), err)
			return nil

		}

		logs = append(logs, l)

	}

	return logs

} // get

func del(log_id string) bool {

  _, err := data.Exec(
		LogDelete, log_id,
	)

	if err != nil {
		
		log.Printf("[%s][Error] %s", version(), err)
		return false

	}

	return true

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