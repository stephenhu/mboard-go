package main

import (
	//"encoding/json"
	"log"
  "net/http"

  "github.com/gorilla/mux"
)

const (
	POWEROFF			= "POWEROFF"
	SUSPEND       = "SUSPEND"
	REBOOT        = "REBOOT"
)

func machineHandler(w http.ResponseWriter, r *http.Request) {

	mux := mux.Vars(r)

  switch r.Method {
  case http.MethodPost:

		operation := mux["operation"]

		log.Println(operation)

		switch operation {
		case POWEROFF:
		  log.Println("p")
		case SUSPEND:
		  log.Println("s")

		case REBOOT:
		  log.Println("r")
		}

  case http.MethodGet:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // machineHandler
