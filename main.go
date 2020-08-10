package main

import (
	"flag"
	"fmt"
  "log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gomodule/redigo/redis"
	"github.com/skip2/go-qrcode"
)

var conf = flag.String("conf", APP_CONFIG, "configuration file")
var mode = flag.Int("mode", MODE_WIRED, "configuration mode")


//var data *sql.DB = nil

var Red *redis.Pool


func addr() string {
	return fmt.Sprintf(":%s", app.Server.Port)
} // addr


func version() string {
  return fmt.Sprintf("%s v%s", APP_NAME, APP_VERSION)
} // version


func connectRedis() {

	Red = &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {return redis.Dial(app.Store.Protocol, app.Store.Port)},
	}

} // connectRedis


func initRouter() *mux.Router {

  router := mux.NewRouter()

  router.PathPrefix("/www/").Handler(http.StripPrefix("/www/",
		http.FileServer(http.Dir("./www"))))

	router.PathPrefix("/blobs/").Handler(http.StripPrefix("/blobs/",
		http.FileServer(http.Dir("./blobs"))))

	router.HandleFunc("/api/games", gameHandler)
	router.HandleFunc("/api/games/{id:[0-9a-f]+}", gameHandler)
  router.HandleFunc("/api/scores", scoreHandler)
	router.HandleFunc("/api/scores/{id:[0-9a-f]+}", scoreHandler)
	router.HandleFunc("/api/scores/{id:[0-9a-f]+}/logs", logHandler)
	router.HandleFunc("/api/media", mediaHandler)
	router.HandleFunc("/api/version", versionHandler)
	//router.HandleFunc("/blobs/{id:[0-9a-f]+}", blobHandler)

	// management apis

	router.HandleFunc("/api/mgmt/details", detailsHandler)
	router.HandleFunc("/api/mgmt/machine", machineHandler)

	router.HandleFunc("/ws/clocks/{id:[0-9a-f]+}", clockControlHandler)
	router.HandleFunc("/ws/scores/{id:[0-9a-f]+}", scoreControlHandler)
	router.HandleFunc("/ws/subscribers/{id:[0-9a-f]+}", subscriberHandler)
	router.HandleFunc("/ws/manager", managerHandler)

  return router

} // initRouter


func generateQR() {

	log.Println("generating QR code...")

	ip, err := getAddress2()

	if err != nil {
		log.Fatal(err)
	} else {

		err := qrcode.WriteFile(fmt.Sprintf("%s/home", ip), qrcode.Medium,
		  1024, fmt.Sprintf("%s/%s", MBOARD_WWW, QR_FILE))

		if err != nil {
			log.Println(err)
		}

	}

} // generateQR


func main() {

  flag.Parse()

	appConfig()

	generateQR()

	connectRedis()

  router := initRouter()

	log.Printf("[%s] listening on port %s", version(), app.Server.Port)

	log.Fatal(http.ListenAndServe(addr(), router))

} // main
