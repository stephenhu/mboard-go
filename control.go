package main

import (
	"encoding/json"
	"fmt"
	"log"
  "net/http"
	//"sync"

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

func calcTotalScore(home bool) int {

  total := 0

	if home {

    for _, v := range game.GameData.Home.Points {
			total = total + v
		}

	} else {

    for _, v := range game.GameData.Away.Points {
			total = total + v
		}

	}

  return total

} // calcTotalScore

func incrementPoints(name string, val int) {

  if game == nil {
		return
	}

  if name == HOME {

		total := game.GameData.Home.Points[game.GameData.Period]
		
		if (total + val) < 0 {
			return
		}

		game.GameData.Home.Points[game.GameData.Period] = total +
			val

		pushString(WS_RET_HOME_SCORE, fmt.Sprintf("%d", calcTotalScore(true)))
		
	} else if name == AWAY {

		total := game.GameData.Away.Points[game.GameData.Period]
		
		if (total + val) < 0 {
			return
		}

		game.GameData.Away.Points[game.GameData.Period] = total +
			val

		pushString(WS_RET_AWAY_SCORE, fmt.Sprintf("%d", calcTotalScore(false)))
		
	}

} // incrementPoints

func incrementFoul(name string, val int) {

  if game == nil {
		return
	}

  if name == HOME {

		if game.GameData.Home.Fouls + val < 0 {
			return
		}

		game.GameData.Home.Fouls = game.GameData.Home.Fouls + val
		
		pushString(WS_RET_HOME_FOUL, fmt.Sprintf("%d", game.GameData.Home.Fouls))

	} else if name == AWAY {

		if game.GameData.Away.Fouls + val < 0  {
			return
		}

		game.GameData.Away.Fouls = game.GameData.Away.Fouls + val

		pushString(WS_RET_AWAY_FOUL, fmt.Sprintf("%d", game.GameData.Away.Fouls))

	}

} // incrementFoul

func incrementTimeout(name string, val int) bool {

  if game == nil {
		return false
	}

  if name == HOME {

		if game.GameData.Home.Timeouts + val < 0 {
				
			return false

		} else if game.Settings.Timeouts < (game.GameData.Home.Timeouts + val) {
		
			return false

		} else {
		
			game.GameData.Home.Timeouts = game.GameData.Home.Timeouts + val
			
			if val == -1 {
				pushString(WS_RET_HOME_TIMEOUT, fmt.Sprintf(
					"%d", game.GameData.Home.Timeouts))
			} else if val == 1 {
				pushString(WS_RET_HOME_TIMEOUT_CANCEL, fmt.Sprintf(
					"%d", game.GameData.Home.Timeouts))
			}

			return true

		}
		
	} else if name == AWAY {

		if game.GameData.Away.Timeouts + val < 0 {

			return false

		} else if game.Settings.Timeouts < (game.GameData.Away.Timeouts + val) {
			
			return false

		} else {

			game.GameData.Away.Timeouts = game.GameData.Away.Timeouts + val
			
			if val == -1 {
				pushString(WS_RET_AWAY_TIMEOUT, fmt.Sprintf("%d", game.GameData.Away.Timeouts))
			} else {
				pushString(WS_RET_AWAY_TIMEOUT_CANCEL, fmt.Sprintf("%d", game.GameData.Away.Timeouts))
			}
	
			return true
			
		}	

	} else {
		return false
	}

} // incrementTimeout

func incrementPeriod(val int) {

	if (game.GameData.Period + val) < 0 {
		return
	}

  game.GameData.Period = game.GameData.Period + val

	game.GameData.Clk.GameClockReset()

  pushString(WS_RET_PERIOD, fmt.Sprintf("%d", game.GameData.Period))

} // incrementPeriod

func setPossession(name string, stopClock bool) {

  if name == HOME {
 	 	
		game.GameData.Possession = true

		pushString(WS_RET_POSSESSION_HOME, fmt.Sprintf("%b", stopClock))

	} else if name == AWAY {
		
		game.GameData.Possession = false

		pushString(WS_RET_POSSESSION_AWAY, fmt.Sprintf("%b", stopClock))

	} else {
		log.Println("Error: setPossession(), invalid possession string.")
	}

	game.GameData.Clk.ShotClockReset()

	if stopClock {
		game.GameData.Clk.Stop()
	} else {
		game.GameData.Clk.Start()
	}

} // setPossession

func togglePossession(stopClock bool) {

	if game.GameData.Possession {
		
		game.GameData.Possession = false
		pushString(WS_RET_POSSESSION_AWAY, fmt.Sprintf("%b", stopClock))

	} else {
		game.GameData.Possession = true
	  pushString(WS_RET_POSSESSION_HOME, fmt.Sprintf("%b", stopClock))
	}

	game.GameData.Clk.ShotClockReset()

	if stopClock {
		game.GameData.Clk.Stop()
	} else {
		game.GameData.Clk.Start()
	}

} // togglePossession

func firehose(game *GameInfo) {

  for {

		select {
		case <-game.GameData.Clk.ShotViolationChan:
		
		  game.GameData.Clk.Ticker.Stop()
			
			// TODO: play sound
			pushString(WS_RET_SHOT_VIOLATION, "1")
			togglePossession(true)
			
		
		case <-game.GameData.Clk.FinalChan:

		  game.GameData.Clk.Ticker.Stop()
			pushString(WS_RET_END_PERIOD, "1")
		
		case s := <-game.GameData.Clk.OutChan:
		  pushString(WS_RET_CLOCK, string(s))
		}

	}

} // firehose

func controlHandler(w http.ResponseWriter, r *http.Request) {

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

	if game.Settings == nil {
		log.Println("game.Settings is nil")
		return
	}
  
	go firehose(game)

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

    req := Req{}

		json.Unmarshal(msg, &req)
		
		log.Println(req)

		req.Period = game.GameData.Period

		switch req.Cmd {
		case WS_CLOCK_START:
			go game.GameData.Clk.Start()

		case WS_CLOCK_STOP:
		  game.GameData.Clk.Stop()

		case WS_CLOCK_RESET:
		  game.GameData.Clk.GameClockReset()

		case WS_SHOT_RESET:
		  game.GameData.Clk.ShotClockReset()

    case WS_SHOT_STEP:
		  game.GameData.Clk.StepShotClock(req.Step)

		case WS_CLOCK_STEP:
		  game.GameData.Clk.StepGameClock(req.Step)

		case WS_PERIOD_UP:
			incrementPeriod(1)
		
		case WS_PERIOD_DOWN:
		  incrementPeriod(-1)

		case WS_POSSESSION_HOME:
			setPossession(HOME, req.Meta["stop"].(bool))

		case WS_POSSESSION_AWAY:
		  setPossession(AWAY, req.Meta["stop"].(bool))

		case WS_FINAL:
			game.Final = true

		case WS_ABORT:

			game.GameData.Clk.Stop()
		
		case WS_SCORE_HOME:
      incrementPoints(HOME, req.Step)

    case WS_SCORE_AWAY:
      incrementPoints(AWAY, req.Step)

		case WS_FOUL_HOME_UP:
		  incrementFoul(HOME, 1)

		case WS_FOUL_HOME_DOWN:
			incrementFoul(HOME, -1)
		
		case WS_FOUL_AWAY_UP:
		  incrementFoul(AWAY, 1)

		case WS_FOUL_AWAY_DOWN:
			incrementFoul(AWAY, -1)

		case WS_TIMEOUT_HOME:
			
			if !incrementTimeout(HOME, -1) {
				req.Reason = MSG_NO_TIMEOUTS
			} else {
				game.GameData.Clk.Stop()
			}
			
		case WS_TIMEOUT_HOME_CANCEL:
			
			if !incrementTimeout(HOME, 1) {
				req.Reason = MSG_MAX_TIMEOUTS
			}

		case WS_TIMEOUT_AWAY:
			
			if !incrementTimeout(AWAY, -1) {
				req.Reason = MSG_NO_TIMEOUTS
			} else {
				game.GameData.Clk.Stop()
			}

		case WS_TIMEOUT_AWAY_CANCEL:
			
			if !incrementTimeout(AWAY, 1) {
				req.Reason = MSG_MAX_TIMEOUTS
			}

		case WS_GAME_STATE:

			state := getGameState()

			log.Println(state);
			if state != nil {
				pushState(state)
			}
			
		
		default:
		  log.Printf("[%s][Error] unsupported command: %s", version(), string(msg))
		}

		put(fmt.Sprintf("%d", game.ID), game.GameData.Clk.PlayClock, req)

	}

} // controlHandler
