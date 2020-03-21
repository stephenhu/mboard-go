package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/eknkc/amber"
)

func getMode() (string, error) {

  switch *mode {

	case MODE_WIFI:
		return INTERFACE_WIFI, nil
	case MODE_HOTSPOT:
	  return INTERFACE_HOTSPOT, nil
	case MODE_WIRED:
	  return INTERFACE_WIRED, nil
	case MODE_TEST:
		return INTERFACE_TEST, nil
	case MODE_CLOUD:
	  return INTERFACE_CLOUD, nil
	default:
	  return INTERFACE_ERROR, errors.New("Unsupported mode of configuration")
	}

} // getMode


func getAddress() (string, error) {

  name, modeErr := getMode()

	if modeErr != nil {
		log.Fatal(modeErr)
	}

	if name == INTERFACE_CLOUD {
		return fmt.Sprintf(":%s", *port), nil
	}

  ifs, err := net.Interfaces()

	if err != nil {
		log.Fatal("Unable to identify interfaces", err)
	}

	for _, iface := range ifs {

		if strings.HasPrefix(iface.Name, name) {
			
			addrs, err := iface.Addrs()

			if err != nil {
				log.Println(err)
				break
			}

			for _, addr := range addrs {

				ipnet, ok := addr.(*net.IPNet)

				//if ok && !ipnet.IP.IsLoopback() {
				if ok {

					if ipnet.IP.To4() != nil {
						return fmt.Sprintf("%s:%s", ipnet.IP.String(), *port), nil
					}
					
				}

			}

		}

	}

	return TEST_ADDRESS, errors.New("Unable to configure given mode, interface has no IP address")

} // getAddress

func setupHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodGet:

		compiler := amber.New()

		parseErr := compiler.ParseFile("mboard-www/setup.amber")

		if parseErr != nil {
			
			log.Printf("[%s][Error] %s", version(), parseErr)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		template, compileErr := compiler.Compile()

		if compileErr != nil {
			
			log.Printf("[%s][Error] %s", version(), compileErr)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		template.Execute(w, nil)		

	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // setupHandler
