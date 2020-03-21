package main

import (
	"encoding/json"
	"log"
	"time"
)

type Clock struct {
	Tenths			int		`json:"tenths"`
	Seconds 		int		`json:"seconds"`
}

type GameClocks struct {
	Ticker            *time.Ticker
	ShotViolationChan chan bool
	FinalChan         chan bool
	OutChan           chan []byte
	PlayClock         *Clock
	ShotClock         *Clock
}

type ReadableClock struct {
	GameClock		*Clock		`json:"game"`
	ShotClock   *Clock		`json:"shot"`
	Minutes     int				`json:"minutes"`
	Shot        int				`json:"shotclock"`
}

func (gc *GameClocks) ClockOut() {

  rc := ReadableClock{
		ShotClock: gc.ShotClock,
		GameClock: gc.PlayClock,
		Minutes: game.Settings.Minutes,
		Shot: game.Settings.Shot,
	}

	j, jsonErr := json.Marshal(rc)

	if jsonErr != nil {
		log.Println("[Error]", jsonErr)
	}

	gc.OutChan <- j
	
} // ClockOut

func (gc *GameClocks) Run() {

	for _ = range gc.Ticker.C {

		if gc.ShotClock.Seconds == game.Settings.Shot {

			gc.ShotClock.Tenths 	= 0
			gc.ShotClock.Seconds 	= 0
		
		}

		if gc.PlayClock.Tenths == 9 {

			gc.PlayClock.Tenths = 0
			gc.PlayClock.Seconds++

		} else {
			gc.PlayClock.Tenths++
		}

		if gc.ShotClock.Tenths == 9 {

			gc.ShotClock.Tenths = 0
			gc.ShotClock.Seconds++
		
		} else {
			gc.ShotClock.Tenths++
		}

		if gc.PlayClock.Seconds == game.Settings.Minutes * 60 {
			gc.FinalChan <- true
		}

		if gc.ShotClock.Seconds == game.Settings.Shot && game.Settings.Shot != -1 {
			gc.ShotViolationChan <- true
		}

		gc.ClockOut()
		
	}

} // Run

func (gc *GameClocks) Start() {

	//TODO: prevent multiple starts
	if gc.Ticker != nil {
		gc.Ticker.Stop()
	}

	gc.Ticker = time.NewTicker(time.Millisecond * 100)

	go gc.Run()

} // Start

func (gc *GameClocks) Stop() {

	if gc.Ticker != nil {
		gc.Ticker.Stop()
	}

} // Stop

func (gc *GameClocks) ShotClockReset() {

	if gc.Ticker != nil {

		gc.Ticker.Stop()

		gc.ShotClock.Seconds 	= 0
		gc.ShotClock.Tenths 	= 0
	
		gc.ClockOut()

	}

} // ShotClockReset

func (gc *GameClocks) GameClockReset() {

	if gc.Ticker != nil {
		gc.Ticker.Stop()
	}

	gc.PlayClock.Seconds 	= 0
	gc.PlayClock.Tenths 	= 0
	gc.ShotClock.Seconds  = 0
	gc.ShotClock.Tenths   = 0

	gc.ClockOut()

} // GameClockReset

func (gc *GameClocks) StepGameClock(ticks int) {

	if gc.Ticker != nil {
		gc.Ticker.Stop()
	}

	total := gc.PlayClock.Seconds + ticks

	if total >= 0 && total <= game.Settings.Minutes * 60 {
		gc.PlayClock.Seconds = total
	}

	if total == game.Settings.Minutes * 60 {
		gc.FinalChan <- true
	}

  gc.ClockOut()

} // StepGameClock

func (gc *GameClocks) StepShotClock(ticks int) {

	if gc.Ticker != nil {
		gc.Ticker.Stop()
	}

  total := gc.ShotClock.Seconds + ticks

	if total >= 0 && total <= game.Settings.Shot {
		gc.ShotClock.Seconds = total
	}

	if gc.ShotClock.Seconds == game.Settings.Shot && game.Settings.Shot != -1 {
		gc.ShotViolationChan <- true
	}

	gc.ClockOut()

} // StepShotClock
