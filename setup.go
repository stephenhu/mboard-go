package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	//"github.com/skip2/go-qrcode"
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

  name, err := getMode()

	if err != nil {
		log.Fatal(err)
	} else {

		if name == INTERFACE_CLOUD {
			return fmt.Sprintf(":%s", app.Server.Port), nil
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
							return fmt.Sprintf("%s:%s", ipnet.IP.String(), app.Server.Port), nil
						}

					}

				}

			}

		}

	}

	return TEST_ADDRESS, errors.New("Unable to configure given mode, interface has no IP address")

} // getAddress


func getAddress2() (string, error) {

  ifs, err := net.Interfaces()

	if err != nil {
		log.Fatal("Unable to identify interfaces", err)
	}

	for _, iface := range ifs {

		if strings.HasPrefix(iface.Name, INTERFACE_WIRED) || strings.HasPrefix(iface.Name, INTERFACE_WIFI) {

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
						return fmt.Sprintf("%s:%s", ipnet.IP.String(), app.Server.Port), nil
					}

				}

			}

		}

	}

	return TEST_ADDRESS, errors.New("Unable to configure given mode, interface has no IP address")

} // getAddress2

