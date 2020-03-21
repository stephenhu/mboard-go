package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/eknkc/amber"
)


func pageHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
  case http.MethodGet:

		compiler := amber.New()
		
		resource := strings.Trim(r.URL.Path, "/")

		err := compiler.ParseFile(fmt.Sprintf("mboard-www/%s.amber", resource))

		if err != nil {
			
			log.Printf("[%s][Error] %s", version(), err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		template, err2 := compiler.Compile()

		if err2 != nil {
			
			log.Printf("[%s][Error] %s", version(), err2)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		template.Execute(w, nil)

	case http.MethodPost:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // pageHandler
