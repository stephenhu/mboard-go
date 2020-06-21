package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/eknkc/amber"
)

var PageIndex = map[string]string {
	"": 					"index.amber",
	"ads":   			"ads.amber",
	"home": 			"home.amber",
	"media":   "media.amber",
	"monitor":   "monitor.amber",
	"setup": 			"setup.amber",
	"settings":   "settings.amber",
}


func pageHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
  case http.MethodGet:

		compiler := amber.New()

		p := strings.Trim(r.URL.Path, "/")

		log.Println(PageIndex[p])

		page, ok := PageIndex[p]

		if !ok {
			w.WriteHeader(http.StatusNotFound)
		} else {

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
					template.Execute(w, nil)
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
