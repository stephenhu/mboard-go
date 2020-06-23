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
	WS_CLOCK_START			= "CLOCK_START"
	WS_CLOCK_STOP				= "CLOCK_STOP"
	WS_CLOCK_RESET      = "CLOCK_RESET"
	WS_CLOCK_STEP       = "CLOCK_STEP"
	WS_SHOT_RESET       = "SHOT_RESET"
	WS_SHOT_STEP        = "SHOT_STEP"
	WS_PERIOD_UP        = "PERIOD_UP"
	WS_PERIOD_DOWN      = "PERIOD_DOWN"
	WS_POSSESSION_HOME  = "POSSESSION_HOME"
	WS_POSSESSION_AWAY  = "POSSESSION_AWAY"
	WS_FINAL            = "FINAL"
	WS_ABORT      			= "ABORT"
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
	WS_TIMEOUT_HOME  					= "TIMEOUT_HOME"
	WS_TIMEOUT_HOME_CANCEL   	= "TIMEOUT_HOME_CANCEL"
	WS_TIMEOUT_AWAY       		= "TIMEOUT_AWAY"
	WS_TIMEOUT_AWAY_CANCEL    = "TIMEOUT_AWAY_CANCEL"
)

const (
	WS_RET_POSSESSION_HOME    	= "POSSESSION_HOME"
	WS_RET_POSSESSION_AWAY    	= "POSSESSION_AWAY"
	WS_RET_HOME_SCORE    				= "HOME_SCORE"
	WS_RET_AWAY_SCORE    				= "AWAY_SCORE"
	WS_RET_CLOCK              	= "CLOCK"
	WS_RET_PERIOD             	= "PERIOD"
	WS_RET_HOME_FOUL          	= "HOME_FOUL"
	WS_RET_AWAY_FOUL          	= "AWAY_FOUL"
	WS_RET_HOME_TIMEOUT       	= "HOME_TIMEOUT"
	WS_RET_HOME_TIMEOUT_CANCEL  = "HOME_TIMEOUT_CANCEL"
	WS_RET_AWAY_TIMEOUT       	= "AWAY_TIMEOUT"
	WS_RET_AWAY_TIMEOUT_CANCEL  = "AWAY_TIMEOUT_CANCEL"
	WS_RET_TIMEOUT_FAIL_MAX     = "TIMEOUT_FAILURE_MAX"
	WS_RET_TIMEOUT_FAIL_NONE  	= "TIMEOUT_FAILURE_NONE"
	WS_RET_SHOT_VIOLATION     	= "SHOT_VIOLATION"
	WS_RET_END_PERIOD         	= "END_PERIOD"
	WS_RET_GAME_STATE           = "GAME_STATE"
	WS_RET_FINAL           			= "GAME_FINAL"
)

const (
	MSG_MAX_TIMEOUTS				= "Maximum timeouts reached."
	MSG_NO_TIMEOUTS         = "No timeouts remaining."
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

var periodNames = []string{"1st", "2nd", "3rd", "4th"}

//Games[id], has sockets, each socket has a mutex, and record

//var connections = make(map[string]map[*websocket.Conn]*sync.Mutex)


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


func incrementTimeout(id string, name string, val int, g *GameInfo) bool {

  if g == nil {
		return false
	}

  if name == HOME {

		if g.GameData.Home.Timeouts + val < 0 {

			return false

		} else if g.Settings.Timeouts < (g.GameData.Home.Timeouts + val) {

			return false

		} else {

			g.GameData.Home.Timeouts = g.GameData.Home.Timeouts + val

			if val == -1 {
				pushString(id, WS_RET_HOME_TIMEOUT, fmt.Sprintf(
					"%d", g.GameData.Home.Timeouts))
			} else if val == 1 {
				pushString(id, WS_RET_HOME_TIMEOUT_CANCEL, fmt.Sprintf(
					"%d", g.GameData.Home.Timeouts))
			}

			return true

		}

	} else if name == AWAY {

		if g.GameData.Away.Timeouts + val < 0 {

			return false

		} else if g.Settings.Timeouts < (g.GameData.Away.Timeouts + val) {

			return false

		} else {

			g.GameData.Away.Timeouts = g.GameData.Away.Timeouts + val

			if val == -1 {
				pushString(id, WS_RET_AWAY_TIMEOUT, fmt.Sprintf("%d", g.GameData.Away.Timeouts))
			} else {
				pushString(id, WS_RET_AWAY_TIMEOUT_CANCEL, fmt.Sprintf("%d", g.GameData.Away.Timeouts))
			}

			return true

		}

	} else {
		return false
	}

} // incrementTimeout


func incrementPeriod(id string, val int, g *GameInfo) {

	if (g.GameData.Period + val) < 0 {
		return
	}

  g.GameData.Period = g.GameData.Period + val

	g.GameData.Clk.GameClockReset()

  pushString(id, WS_RET_PERIOD, fmt.Sprintf("%d", g.GameData.Period))

} // incrementPeriod


func setPossession(id string, name string, stopClock bool, g *GameInfo) {

  if name == HOME {

		g.GameData.Possession = true

		pushString(id, WS_RET_POSSESSION_HOME, fmt.Sprintf("%t", stopClock))

	} else if name == AWAY {

		g.GameData.Possession = false

		pushString(id, WS_RET_POSSESSION_AWAY, fmt.Sprintf("%t", stopClock))

	} else {
		log.Println("Error: setPossession(), invalid possession string.")
	}

	g.GameData.Clk.ShotClockReset()

	if stopClock {
		g.GameData.Clk.Stop()
	} else {
		g.GameData.Clk.Start()
	}

} // setPossession


func togglePossession(id string, stopClock bool, g *GameInfo) {

	if g.GameData.Possession {

		g.GameData.Possession = false
		pushString(id, WS_RET_POSSESSION_AWAY, fmt.Sprintf("%t", stopClock))

	} else {
		g.GameData.Possession = true
	  pushString(id, WS_RET_POSSESSION_HOME, fmt.Sprintf("%t", stopClock))
	}

	g.GameData.Clk.ShotClockReset()

	if stopClock {
		g.GameData.Clk.Stop()
	} else {
		g.GameData.Clk.Start()
	}

} // togglePossession


func endGame(id string) {

	delete(gameMap, id)
	delete(subscribersMap, id)

} // endGame


func firehose(id string, g *GameInfo) {

	for {

		select {
		case <-g.GameData.Clk.ShotViolationChan:

			g.GameData.Clk.Ticker.Stop()

			// TODO: play sound
			pushString(id, WS_RET_SHOT_VIOLATION, "1")
			togglePossession(id, true, g)


		case <-g.GameData.Clk.FinalChan:

			g.GameData.Clk.Ticker.Stop()
			pushString(id, WS_RET_END_PERIOD, "1")

		case s := <-g.GameData.Clk.OutChan:
			pushString(id, WS_RET_CLOCK, string(s))
		}

	}

} // firehose


func controlHandler(w http.ResponseWriter, r *http.Request) {

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

	if g.Final {
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

	go firehose(id, &g)

	defer c.Close()

	for {

		if g.Final {
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

		req.Period = g.GameData.Period

		switch req.Cmd {
		case WS_CLOCK_START:
			go g.GameData.Clk.Start()

		case WS_CLOCK_STOP:
			g.GameData.Clk.Stop()

		case WS_CLOCK_RESET:
			g.GameData.Clk.GameClockReset()

		case WS_SHOT_RESET:
			g.GameData.Clk.ShotClockReset()

		case WS_SHOT_STEP:
			g.GameData.Clk.StepShotClock(req.Step)

		case WS_CLOCK_STEP:
			g.GameData.Clk.StepGameClock(req.Step)

		case WS_PERIOD_UP:
			incrementPeriod(id, 1, &g)

		case WS_PERIOD_DOWN:
			incrementPeriod(id, -1, &g)

		case WS_POSSESSION_HOME:
			setPossession(id, HOME, req.Meta["stop"].(bool), &g)

		case WS_POSSESSION_AWAY:
			setPossession(id, AWAY, req.Meta["stop"].(bool), &g)

		case WS_FINAL:
			g.Final = true
			endGame(id)

		case WS_ABORT:
			g.GameData.Clk.Stop()

		case WS_SCORE_HOME:
			incrementPoints(id, HOME, req.Step, &g)

		case WS_SCORE_AWAY:
			incrementPoints(id, AWAY, req.Step, &g)

		case WS_FOUL_HOME_UP:
			incrementFoul(id, HOME, 1, &g)

		case WS_FOUL_HOME_DOWN:
			incrementFoul(id, HOME, -1, &g)

		case WS_FOUL_AWAY_UP:
			incrementFoul(id, AWAY, 1, &g)

		case WS_FOUL_AWAY_DOWN:
			incrementFoul(id, AWAY, -1, &g)

		case WS_TIMEOUT_HOME:

			if !incrementTimeout(id, HOME, -1, &g) {
				req.Reason = MSG_NO_TIMEOUTS
			} else {
				g.GameData.Clk.Stop()
			}

		case WS_TIMEOUT_HOME_CANCEL:

			if !incrementTimeout(id, HOME, 1, &g) {
				req.Reason = MSG_MAX_TIMEOUTS
			}

		case WS_TIMEOUT_AWAY:

			if !incrementTimeout(id, AWAY, -1, &g) {
				req.Reason = MSG_NO_TIMEOUTS
			} else {
				g.GameData.Clk.Stop()
			}

		case WS_TIMEOUT_AWAY_CANCEL:

			if !incrementTimeout(id, AWAY, 1, &g) {
				req.Reason = MSG_MAX_TIMEOUTS
			}

		case WS_GAME_STATE:

			gs := GameState{
				Settings: g.Settings,
				Period: g.GameData.Period,
				Possession: g.GameData.Possession,
				Home: g.GameData.Home,
				Away: g.GameData.Away,
				GameClock: g.GameData.Clk.PlayClock,
				ShotClock: g.GameData.Clk.ShotClock,
			}

			pushState(id, &gs)

		default:
			log.Printf("[%s][Error] unsupported command: %s", version(), string(msg))
		}

		put(fmt.Sprintf("%s", ), g.GameData.Clk.PlayClock, req)

	}

} // controlHandler
