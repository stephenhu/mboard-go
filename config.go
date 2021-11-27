package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Service struct {
	Host						string				`json:"host"`
	Port						string  			`json:"port"`
	Protocol				string				`json:"protocol"`
}

type AppConfig struct {
	Store						Service				`json:"redis"`
	Server					Service				`json:"server"`
}

var app AppConfig


func appConfig() {

	var f string

	if *conf != "" {
		f = *conf
	} else {
		f = APP_CONFIG
	}

	_, err := os.Stat(f)

	if err != nil || os.IsNotExist(err) {
		log.Fatal(err)
	} else {

		buf, err := ioutil.ReadFile(f)

		if err != nil {
			log.Fatal(err)
		} else {

			err := json.Unmarshal(buf, &app)

			if err != nil {
				log.Fatal(err)
			}

		}

	}

} // appConfig
