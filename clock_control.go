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
)

const (
	WS_TIMEOUT_HOME  					= "TIMEOUT_HOME"
	WS_TIMEOUT_HOME_CANCEL   	= "TIMEOUT_HOME_CANCEL"
	WS_TIMEOUT_AWAY       		= "TIMEOUT_AWAY"
	WS_TIMEOUT_AWAY_CANCEL    = "TIMEOUT_AWAY_CANCEL"
)

const (
	WS_RET_POSSESSION_HOME    	= "POSSESSION_HOME"
	WS_RET_POSSESSION_AWAY    	= "POSSESSION_AWAY"
	WS_RET_CLOCK              	= "CLOCK"
	WS_RET_PERIOD             	= "PERIOD"
	WS_RET_HOME_TIMEOUT       	= "HOME_TIMEOUT"
	WS_RET_HOME_TIMEOUT_CANCEL  = "HOME_TIMEOUT_CANCEL"
	WS_RET_AWAY_TIMEOUT       	= "AWAY_TIMEOUT"
	WS_RET_AWAY_TIMEOUT_CANCEL  = "AWAY_TIMEOUT_CANCEL"
	WS_RET_TIMEOUT_FAIL_MAX     = "TIMEOUT_FAILURE_MAX"
	WS_RET_TIMEOUT_FAIL_NONE  	= "TIMEOUT_FAILURE_NONE"
	WS_RET_SHOT_VIOLATION     	= "SHOT_VIOLATION"
	WS_RET_END_PERIOD         	= "END_PERIOD"
	WS_RET_FINAL           			= "GAME_FINAL"
)

const (
	MSG_MAX_TIMEOUTS				= "Maximum timeouts reached."
	MSG_NO_TIMEOUTS         = "No timeouts remaining."
)

var periodNames = []string{"1st", "2nd", "3rd", "4th"}


func checkGameOver(g *GameInfo) bool {

	home := calcTotalScore(true, g)
	away := calcTotalScore(false, g)

	if g.GameData.Period >= 3 && away == home {
		return true
	} else {
		return false
	}

} // checkGameOver


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

	// TODO: check game over condition

	if checkGameOver(g) {
		endGame(id)
	} else {

		g.GameData.Period = g.GameData.Period + val

		g.GameData.Clk.GameClockReset()

		pushString(id, WS_RET_PERIOD, fmt.Sprintf("%d", g.GameData.Period))

	}

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

	g, ok := gameMap[id]

	if ok {

		if g.ScoreCtl != nil {
			g.ScoreCtl.Close()
		}

		if g.ClockCtl != nil {
			g.ClockCtl.Close()
		}

		cleanupSubscribers(id)
		delete(gameMap, id)
		delete(subscribersMap, id)

	}

} // endGame


func firehose(id string, g *GameInfo) {

	for {

		select {
		case <-g.GameData.Clk.ShotViolationChan:

			g.GameData.Clk.Ticker.Stop()

			// TODO: play sound

			pushString(id, WS_RET_SHOT_VIOLATION, "1")
			//togglePossession(id, true, g)

		case <-g.GameData.Clk.FinalChan:

			g.GameData.Clk.Ticker.Stop()

			endGame(id)

			pushString(id, WS_RET_END_PERIOD, "1")

		case s := <-g.GameData.Clk.OutChan:
			pushString(id, WS_RET_CLOCK, string(s))

		}

	}

} // firehose


func clockControlHandler(w http.ResponseWriter, r *http.Request) {

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

	if g.ClockCtl == nil {
		log.Println("socket connection doesn't exist, store socket connection")
		g.ClockCtl = c
	} else {
		c.Close()
		return
	}

	defer c.Close()

	go firehose(id, g.Game)

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
		case WS_CLOCK_START:
			g.Game.GameData.Clk.Start()

		case WS_CLOCK_STOP:
			g.Game.GameData.Clk.Stop()

		case WS_CLOCK_RESET:
			g.Game.GameData.Clk.GameClockReset()

		case WS_SHOT_RESET:
			g.Game.GameData.Clk.ShotClockReset()

		case WS_SHOT_STEP:
			g.Game.GameData.Clk.StepShotClock(req.Step)

		case WS_CLOCK_STEP:
			g.Game.GameData.Clk.StepGameClock(req.Step)

		case WS_PERIOD_UP:
			incrementPeriod(id, 1, g.Game)

		case WS_PERIOD_DOWN:
			incrementPeriod(id, -1, g.Game)

		case WS_POSSESSION_HOME:
			setPossession(id, HOME, req.Meta["stop"].(bool), g.Game)

		case WS_POSSESSION_AWAY:
			setPossession(id, AWAY, req.Meta["stop"].(bool), g.Game)

		case WS_FINAL:

			g.Game.Final = true
			//endGame(id)
			g.Game.GameData.Clk.FinalChan <- true
			break

		case WS_ABORT:
			g.Game.GameData.Clk.Stop()

		case WS_TIMEOUT_HOME:

			if !incrementTimeout(id, HOME, -1, g.Game) {
				req.Reason = MSG_NO_TIMEOUTS
			} else {
				g.Game.GameData.Clk.Stop()
			}

		case WS_TIMEOUT_HOME_CANCEL:

			if !incrementTimeout(id, HOME, 1, g.Game) {
				req.Reason = MSG_MAX_TIMEOUTS
			}

		case WS_TIMEOUT_AWAY:

			if !incrementTimeout(id, AWAY, -1, g.Game) {
				req.Reason = MSG_NO_TIMEOUTS
			} else {
				g.Game.GameData.Clk.Stop()
			}

		case WS_TIMEOUT_AWAY_CANCEL:

			if !incrementTimeout(id, AWAY, 1, g.Game) {
				req.Reason = MSG_MAX_TIMEOUTS
			}

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

} // clockControlHandler
