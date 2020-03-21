package main

import (
	"database/sql"
	"encoding/json"
  "log"
	"net/http"

  "github.com/gorilla/mux"
)

type GameRecord struct {
	Home		*Team			`json:"home"`
	Away		*Team			`json:"away"`
}

type GameTbl struct {
	ID			string 	`json:"id"`
  Data    sql.NullString `json:"data"`
	Created string 	`json:"created"`
	Updated string 	`json:"updated"`
	Status  int 		`json:"status"`
}


func addGame() int64 {

	res, err := data.Exec(
		GameCreate,
	)

	if err != nil {
		log.Println(err)
		return -1
	}

	id, err := res.LastInsertId()

	if err != nil {
		
		log.Println(err)
		return -1

	}

	return id

} // addGame

func updateGame(id int64, val string) {

	_, err := data.Exec(
		GameUpdate, val, 1, id,
	)

	if err != nil {
		log.Println(err)
	}

} // updateGame

func getGames() []GameTbl {

  rows, err := data.Query(
		GamesGet,
	)

	if err != nil {
		log.Printf("[%s][Error][DB] %s", version(), err)
		return nil
	}

	defer rows.Close()

	gt := []GameTbl{}

	for rows.Next() {

			g := GameTbl{}

			err := rows.Scan(&g.ID, &g.Data, &g.Status, &g.Created, &g.Updated)

			if err == sql.ErrNoRows || err != nil {
				log.Printf("[%s][Error] %s", version(), err)
				return nil
			}

			gt = append(gt, g)

	}

	return gt

} // getGames

func getGame(id string) *GameTbl {

  row := data.QueryRow(
		GameGet, id,
	)

	gt := GameTbl{}

	err := row.Scan(&gt.ID, &gt.Data, &gt.Status, &gt.Created, &gt.Updated)

	if err == sql.ErrNoRows || err != nil {
		log.Printf("[%s][Error] %s", version(), err)
		return nil
	}

	return &gt

} // getGame

func deleteGame(id string) {

  _, err := data.Exec(
		GameDelete, id,
	)

	if err != nil {
		log.Printf("[%s][Error] %s", version(), err)
		return
	}

} // deleteGame

func scoreHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
  case http.MethodPost:
	case http.MethodGet:

		vars := mux.Vars(r)

		id := vars["id"]

		if id != "" {

			g := getGame(id)

			if g == nil {
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
