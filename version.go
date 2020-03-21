package main

import (
	"encoding/json"
	"log"
  "net/http"

)

func versionHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
  case http.MethodGet:

		v := map[string]string{
			"version": VERSION,
		}

		j, jsonErr := json.Marshal(v)

		if jsonErr != nil {
			log.Printf("[%s] %s", version(), jsonErr)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write(j)
		}

		 
  case http.MethodPost:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // versionHandler
