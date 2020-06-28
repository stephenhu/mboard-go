package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	WS_UNDO             = "UNDO"
	WS_GAME_STATE       = "GAME_STATE"
)

const (
	WS_SCORE_HOME       			= "SCORE_HOME"
	WS_SCORE_AWAY       			= "SCORE_AWAY"
	WS_FOUL_HOME_UP     			= "FOUL_HOME_UP"
	WS_FOUL_HOME_DOWN   			=	"FOUL_HOME_DOWN"
	WS_FOUL_AWAY_UP     			= "FOUL_AWAY_UP"
	WS_FOUL_AWAY_DOWN   			= "FOUL_AWAY_DOWN"
)

const (
	WS_RET_HOME_SCORE    				= "HOME_SCORE"
	WS_RET_AWAY_SCORE    				= "AWAY_SCORE"
	WS_RET_HOME_FOUL          	= "HOME_FOUL"
	WS_RET_AWAY_FOUL          	= "AWAY_FOUL"
	WS_RET_GAME_STATE           = "GAME_STATE"
)


type Team struct {
	Name      string    		`json:"name"`
	Logo      string    		`json:"logo"`
	Fouls			int						`json:"fouls"`
	Timeouts  int     			`json:"timeouts"`
	Points    map[int]int   `json:"points"`
}

type Game struct {
	Home				*Team					`json:"home"`
	Away      	*Team					`json:"away"`
	Period    	int						`json:"period"`
	Clk					*GameClocks		`json:"clk"`
	Possession 	bool					`json:"possesion"`
}

type Req struct {
	Cmd					string 			`json:"cmd"`
	Step				int					`json:"step"`
	Meta        map[string]interface{}      `json:"meta"`
	Reason      string      `json:"reason"`
	Timestamp   string      `json:"timestamp"`
	Period      int         `json:"period"`
}


func calcTotalScore(home bool, g *GameInfo) int {

  total := 0

	if home {

    for _, v := range g.GameData.Home.Points {
			total = total + v
		}

	} else {

    for _, v := range g.GameData.Away.Points {
			total = total + v
		}

	}

  return total

} // calcTotalScore


func incrementPoints(id string, name string, val int, g *GameInfo) {

  if g == nil {
		return
	}

  if name == HOME {

		total := g.GameData.Home.Points[g.GameData.Period]

		if (total + val) < 0 {
			return
		}

		g.GameData.Home.Points[g.GameData.Period] = total +
			val

		pushString(id, WS_RET_HOME_SCORE, fmt.Sprintf("%d", calcTotalScore(true, g)))

	} else if name == AWAY {

		total := g.GameData.Away.Points[g.GameData.Period]

		if (total + val) < 0 {
			return
		}

		g.GameData.Away.Points[g.GameData.Period] = total +
			val

		pushString(id, WS_RET_AWAY_SCORE, fmt.Sprintf("%d", calcTotalScore(false, g)))

	}

} // incrementPoints


func incrementFoul(id string, name string, val int, g *GameInfo) {

  if g == nil {
		return
	}

  if name == HOME {

		if g.GameData.Home.Fouls + val < 0 {
			return
		}

		g.GameData.Home.Fouls = g.GameData.Home.Fouls + val

		pushString(id, WS_RET_HOME_FOUL, fmt.Sprintf("%d", g.GameData.Home.Fouls))

	} else if name == AWAY {

		if g.GameData.Away.Fouls + val < 0  {
			return
		}

		g.GameData.Away.Fouls = g.GameData.Away.Fouls + val

		pushString(id, WS_RET_AWAY_FOUL, fmt.Sprintf("%d", g.GameData.Away.Fouls))

	}

} // incrementFoul


func scoreControlHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id := vars["id"]

	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	g, ok := gameMap[id]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if g.Game.Final {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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

	defer c.Close()

	if g.ScoreCtl == nil {
		g.ScoreCtl = c
	} else {
		return
	}

	for {

		if g.Game.Final {
			pushString(id, WS_RET_FINAL, id)
			break
		}

		_, msg, err := c.ReadMessage()

		if err != nil {

			log.Println("[Error] ", err)
			break

		}

		if msg == nil || len(msg) == 0 {
			continue
		}

		req := Req{}

		json.Unmarshal(msg, &req)

		log.Println(req)

		req.Period = g.Game.GameData.Period

		switch req.Cmd {

		case WS_SCORE_HOME:
			incrementPoints(id, HOME, req.Step, g.Game)

		case WS_SCORE_AWAY:
			incrementPoints(id, AWAY, req.Step, g.Game)

		case WS_FOUL_HOME_UP:
			incrementFoul(id, HOME, 1, g.Game)

		case WS_FOUL_HOME_DOWN:
			incrementFoul(id, HOME, -1, g.Game)

		case WS_FOUL_AWAY_UP:
			incrementFoul(id, AWAY, 1, g.Game)

		case WS_FOUL_AWAY_DOWN:
			incrementFoul(id, AWAY, -1, g.Game)

		case WS_GAME_STATE:

			gs := GameState{
				Settings: g.Game.Settings,
				Period: g.Game.GameData.Period,
				Possession: g.Game.GameData.Possession,
				Home: g.Game.GameData.Home,
				Away: g.Game.GameData.Away,
				GameClock: g.Game.GameData.Clk.PlayClock,
				ShotClock: g.Game.GameData.Clk.ShotClock,
			}

			pushState(id, &gs)

		default:
			log.Printf("[Error] unsupported command: %s", string(msg))
		}

		put(id, g.Game.GameData.Clk.PlayClock, req)

	}

} // scoreControlHandler
