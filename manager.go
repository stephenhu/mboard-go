package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	WS_LOGIN						=	"LOGIN"
	WS_LOGO             = "LOGO"
  WS_SETUP            = "SETUP"
	WS_SCOREBOARD				= "SCOREBOARD"
	WS_ADVERTISEMENT		= "ADVERTISEMENT"
	WS_VIDEO_PLAY       = "VIDEO_PLAY"
	WS_VIDEO_STOP				= "VIDEO_STOP"
	WS_PHOTO_PLAY				= "PHOTO_PLAY"
	WS_PHOTO_STOP       = "PHOTO_STOP"
	WS_AUDIO_PLAY       = "AUDIO_PLAY"
	WS_AUDIO_STOP       = "AUDIO_STOP"
)

var manager *websocket.Conn

type ManagerCommand struct {
  Cmd 			string								`json:"cmd"`
	Options		map[string]string			`json:"options"`
}

func managerHandler(w http.ResponseWriter, r *http.Request) {

  upgrader := websocket.Upgrader {
		ReadBufferSize:		1024,
		WriteBufferSize: 	1024,
		CheckOrigin:		func(r *http.Request) bool { return true },
	}

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		
		log.Println("[Error]", err)
		return

	}

	manager = c

	defer c.Close()

	for {
   
		_, msg, err := c.ReadMessage()

		if err != nil {

			log.Println("[Error] ", err)
			break

		}

    if msg == nil {
			log.Println(msg)
			break
		}

    mc := ManagerCommand{}

		json.Unmarshal(msg, &mc)

		switch mc.Cmd {
		case WS_LOGO:
			pushMap(WS_LOGO, mc.Options)

		case WS_LOGIN:
			log.Println(WS_LOGIN)

		case WS_SCOREBOARD:

			if game != nil && game.Active {
				pushMap(WS_SCOREBOARD, mc.Options)
			} else {
				pushMap(WS_SETUP, mc.Options)
			}

		case WS_ADVERTISEMENT:
			pushMap(WS_ADVERTISEMENT, mc.Options)

		case WS_VIDEO_PLAY:

			go videoPlay(mc.Options)

			//pushMap(WS_VIDEO_PLAY, mc.Options)

		case WS_VIDEO_STOP:

			err := videoStop()

			if err != nil {
				log.Println(err)
			}

			//pushMap(WS_VIDEO_STOP, mc.Options)

		case WS_AUDIO_PLAY:
			pushMap(WS_AUDIO_PLAY, mc.Options)

		case WS_AUDIO_STOP:
			pushMap(WS_AUDIO_STOP, mc.Options)

		case WS_PHOTO_PLAY:
			pushMap(WS_PHOTO_PLAY, mc.Options)

		case WS_PHOTO_STOP:
			pushMap(WS_PHOTO_STOP, mc.Options)

		default:
			log.Println("unknown")
		}

	}

} // managerHandler
