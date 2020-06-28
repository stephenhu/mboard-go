package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/eknkc/amber"
	"github.com/gorilla/mux"
)

var PageIndex = map[string]string {
	"": 								"index.amber",
	"ads":   						"ads.amber",
	"clockctl":   			"clockctl.amber",
	"clocks":   				"clock.amber",
	"gamectl":   				"gamectl.amber",
	"gameconfig":   		"gameconfig.amber",
	"home": 						"home.amber",
	"media":   					"media.amber",
	"monitor":   				"monitor.amber",
	"notfound":         "notfound.amber",
	"scoreboards":   		"scoreboard.amber",
	"setup": 						"setup.amber",
	"settings":   			"settings.amber",
}


func pageHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
  case http.MethodGet:

		vars := mux.Vars(r)

		id := vars["id"]

		compiler := amber.New()

		p := strings.Trim(strings.Trim(r.URL.Path, id), "/")

		log.Println(PageIndex[p])

		page, ok := PageIndex[p]

		if !ok {
			w.WriteHeader(http.StatusNotFound)
		} else {

			if id != "" {

				_, ok := gameMap[id]

				if !ok {
					page = "gamenotfound.amber"
				}

			}

			err := compiler.ParseFile(fmt.Sprintf("%s/%s", MBOARD_WWW, page))

			if err != nil {

				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return

			} else {

				template, err := compiler.Compile()

				if err != nil {

					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)

				} else {

					d := map[string]string {
						"id": id,
					}

					template.Execute(w, &d)
				}

			}

		}

	case http.MethodPost:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // pageHandler
