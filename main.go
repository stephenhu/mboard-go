package main

import (
	"database/sql"
	"flag"
	"fmt"
  "log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
)

const (
	APPNAME 					= "mboard-go v%s"
	TEST_ADDRESS 	    = "127.0.0.1:8000"
	CLOUD_ADDRESS     = "madsportslab.com"
	MBOARD            = "mboard"
	VERSION 					= "0.1.0"
)

const (
	MODE_WIFI				= 0
	MODE_HOTSPOT   	= 1
	MODE_WIRED      = 2
	MODE_CLOUD      = 3
	MODE_TEST   		= 4
)

const (
	INTERFACE_WIFI 		= "en"
	INTERFACE_HOTSPOT	= "wlan"
	INTERFACE_WIRED   = "eth"
	INTERFACE_CLOUD   = "cloud"
	INTERFACE_TEST		= "lo"
	INTERFACE_ERROR   = ""
)

var database 	= flag.String("database", "./data/mboard.db", "database address")
var port 			= flag.String("port", "8000", "service port")
var mode      = flag.Int("mode", MODE_WIFI, "configuration mode")
var ssl       = flag.Bool("ssl", false, "use SSL encryption")
var certFile  = flag.String("cert", "ssl.crt", "SSL certificate")
var keyFile   = flag.String("key", "ssl.key", "SSL private key")
var v         = flag.Bool("v", false, "version")

var data *sql.DB = nil

func version() string {
  return fmt.Sprintf(APPNAME, VERSION)
} // version

func initDatabase() {

  db, err := sql.Open("sqlite3", *database)

	if err != nil {
		log.Fatal("Database connection error: ", err)
	}

	data = db

} // initDatabase

func initRouter() *mux.Router {

  router := mux.NewRouter()

  router.PathPrefix("/mboard-www/").Handler(http.StripPrefix("/mboard-www/",
		http.FileServer(http.Dir("./mboard-www"))))
		
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
	
	router.HandleFunc("/clock", pageHandler)
	router.HandleFunc("/shotclock", pageHandler)
	router.HandleFunc("/scoreboard", scoreboardHandler)
	router.HandleFunc("/setup", setupHandler)
	router.HandleFunc("/logo", pageHandler)
	router.HandleFunc("/video/{id:[0-9a-f]+}", videoHandler)
	router.HandleFunc("/photo", photoHandler)
	router.HandleFunc("/advertisement", pageHandler)
	router.HandleFunc("/download", pageHandler)

	//router.HandleFunc("/ws/games/{id:[0-9a-f]+}", controlHandler)
	router.HandleFunc("/ws/game", controlHandler)
	router.HandleFunc("/ws/subscriber", subscriberHandler)
	router.HandleFunc("/ws/manager", managerHandler)
	
  return router

} // initRouter

func generateQR() {

	log.Println("generating QR code...")
	
	ip, err := getAddress()

	if err != nil {
		log.Fatal(err)
	} else {

		err := qrcode.WriteFile(ip, qrcode.Medium, 512, "mboard-www/qr.png")
		
		if err != nil {
			log.Println(err)
		}

	}

} // generateQR

func main() {

  flag.Parse()
	
	if *v {
		fmt.Println(VERSION)
		os.Exit(0)
	}
	
	generateQR()

  log.Printf("[%s] listening on port %s", version(), *port)

  initDatabase()

  router := initRouter()

	addr := fmt.Sprintf(":%s", *port)
	
	if *ssl {
		log.Fatal(http.ListenAndServeTLS(addr, *certFile, *keyFile, router))
	} else {
		log.Fatal(http.ListenAndServe(addr, router))
	}

} // main
