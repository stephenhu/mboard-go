package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eknkc/amber"
)


func notFoundHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodGet:

		compiler := amber.New()

		err := compiler.ParseFile(fmt.Sprintf("%s/%s", MBOARD_WWW, "notfound.amber"))

		template, err := compiler.Compile()

		if err != nil {

			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			template.Execute(w, nil)
		}


	case http.MethodPost:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // notFoundHandler
