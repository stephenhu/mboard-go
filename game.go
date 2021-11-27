package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
  "log"
	"net/http"
	//"strconv"
	"strings"
	//"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

)

const (
	HOME			= "home"
	AWAY			= "away"
	PERIODS   = "periods"
	MINUTES   = "minutes"
	FOULS     = "fouls"
	TIMEOUTS  = "timeouts"
	SHOT      = "shot"
)

type GameConfig struct {
	Sport       string  `json:"sport"`
  Periods			int 		`json:"periods"`
  Minutes			int 		`json:"minutes"`
	Shot			  int 		`json:"shot"`
	Timeouts		int 		`json:"timeouts"`
	Fouls				int 		`json:"fouls"`
	Home        string  `json:"home"`
	Away        string  `json:"away"`
}

type GameInfo struct {
  Settings			*GameConfig
	GameData			*Game
	Final         bool
	Created       int64
	Active        bool
}

type GameCtl struct {
	ScoreCtl 			*websocket.Conn
	ClockCtl 			*websocket.Conn
	Game          *GameInfo
}

//TODO: remove gamestate struct?
type GameState struct {
	Settings      *GameConfig   `json:"settings"`
	Period        int						`json:"period"`
	Possession    bool					`json:"possession"`
	Home          *Team					`json:"home"`
	Away          *Team					`json:"away"`
	GameClock     *Clock    		`json:"game"`
	ShotClock     *Clock    		`json:"shot"`
	Final         bool          `json:"final"`
}

type GameRes struct {
	Msg 	string 	`json:"msg"`
}

var gameMap map[string]*GameCtl

var fields = []string{HOME, AWAY, PERIODS, MINUTES, FOULS, TIMEOUTS, SHOT}


func gameConfig(j string) *GameConfig {

    config := GameConfig{}

		// TODO: check fields

		err := json.Unmarshal([]byte(j), &config)

		if err != nil {
			log.Println(err)
			return nil
		} else {
			return &config
		}

} // gameConfig


func initTeam(name string, timeouts int) *Team {

  team := Team{
		Name: strings.Title(name),
		Points: make(map[int]int),
		Timeouts: timeouts,
	}

  return &team

} // initTeam


func initGameClocks(id string) *GameClocks {

  gc := GameClocks{
		ShotViolationChan: make(chan bool),
		FinalChan: make(chan bool),
		OutChan: make(chan []byte),
		PlayClock: &Clock{Tenths: 0, Seconds: 0},
		ShotClock: &Clock{Tenths: 0, Seconds: 0},
		GameID: id,
	}

	return &gc

} // initGameClocks


func generateId(config *GameConfig, length int) string {

  now := time.Now().String()

  digest := hmac.New(sha256.New, []byte("ABC"))

	digest.Write([]byte(fmt.Sprintf("%s%s", now, config)))

	hash := hex.EncodeToString(digest.Sum(nil))

	return hash[:length]

} // generateId


func gameHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
  case http.MethodPost:

		config := r.FormValue(API_PARAM_GAME_CONFIG)

		log.Println(config)

		if config == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {

			cf := gameConfig(config)

			if cf != nil {

				id := generateId(cf, 32)

				h := initTeam(cf.Home, cf.Timeouts)
				a := initTeam(cf.Away, cf.Timeouts)

				c := initGameClocks(id)

				ts := time.Now().Unix()

				gi := GameInfo{
					Settings:	cf,
					GameData: &Game{
						Home: h,
						Away: a,
						Clk: c,
						Possession: true,
					},
					Final: false,
					Created: ts,
					Active: true,
				}

				if gameMap == nil {
					gameMap = make(map[string]*GameCtl)
				}

				addGame(config, ts)

				gameMap[id] = &GameCtl {
					Game: &gi,
				}

				pushMap(id, WS_SCOREBOARD, nil)

				w.Write([]byte(id))

			} else {
				w.WriteHeader(http.StatusBadRequest)
			}

		}

	case http.MethodGet:

		m := mux.Vars(r)

		id := m["id"]

		if id == "" {

			games := map[string]string{}

			for k, v := range gameMap {
				games[k] = v.Game.Settings.Sport
			}

			log.Println(games)

			j, err := json.Marshal(games)

			if err != nil {
				log.Println(err)
			} else {
				w.Write(j)
			}

		} else {

			//gr := getGameRecord(id)

			g, ok := gameMap[id]

			if ok {

				gs := GameState{
					Settings: g.Game.Settings,
					Period: g.Game.GameData.Period,
					Possession: g.Game.GameData.Possession,
					Home: g.Game.GameData.Home,
					Away: g.Game.GameData.Away,
					GameClock: g.Game.GameData.Clk.PlayClock,
					ShotClock: g.Game.GameData.Clk.ShotClock,
				}

				j, jsonErr := json.Marshal(gs)

				if jsonErr != nil {
					log.Println(jsonErr)
				}

				w.Write(j)

			} else {
				w.WriteHeader(http.StatusNotFound)
			}

		}

	case http.MethodPut:

		m := mux.Vars(r)

		id := m["id"]

		g, ok := gameMap[id]

	  if ok && g.Game.Active {

			g.Game.Active = false
			g.Game.Final 	= true
/*
			gr := GameRecord{}

			gr.
			/*
			gr := GameRecord{
				Home: g.GameData.Home,
				Away: g.GameData.Away,
			}
			*/

/*
			j, jsonErr := json.Marshal(gr)

			if jsonErr != nil {
				log.Printf("[%s] %s", version(), jsonErr)
				w.WriteHeader(http.StatusInternalServerError)
			} else {

				updateGame(g.ID, string(j))

				g.GameData.Clk.Stop()

				pushString(WS_FINAL, "")

				g = &GameInfo{}

				pushMap(WS_SETUP, nil)

				w.WriteHeader(http.StatusOK)

			}
*/
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	case http.MethodDelete:
	default:
		log.Printf("[%s][Error] unsupported command", version())
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // gameHandler
