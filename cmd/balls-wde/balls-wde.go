// Package main is the Bouncing balls (wde) demo app.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/icza/balls-wde/engine"
	wde "github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
)

const (
	version  = "v1.0.0"
	name     = "Bouncing Balls WDE"
	homePage = "https://github.com/icza/balls-wde"
	title    = name + " " + version
)

func main() {
	fmt.Println(title)
	fmt.Println("Home page:", homePage)
	rand.Seed(time.Now().UnixNano())

	go run()

	wde.Run()
}

var eng *engine.Engine // eng is the engine

// run runs the demo.
func run() {
	w, h := 800, 550

	var err error
	win, err := wde.NewWindow(w, h)
	if err != nil {
		log.Printf("Create window error: %v", err)
		return
	}

	win.SetTitle(title)
	win.LockSize(true)
	win.Show()

	eng = engine.NewEngine(win, w, h)
	go eng.Run()

	for event := range win.EventChan() {
		if quit := handleEvent(event); quit {
			break
		}
	}
	eng.Stop()

	wde.Stop()
}

// handleEvent handles events and tells if we need to quit (based on the event).
func handleEvent(event interface{}) (quit bool) {
	switch e := event.(type) {
	case wde.KeyTypedEvent:
		switch gl := e.Glyph; gl {
		case "s", "S":
			eng.ChangeSpeed(gl == "S")
		case "a", "A":
			eng.ChangeMaxBalls(gl == "A")
		case "m", "M":
			eng.ChangeMinMaxBallRatio(gl == "M")
		case "r", "R":
			eng.Restart()
		case "o", "O":
			eng.ToggleOSD()
		case "g", "G":
			eng.ChangeGravityAbs(gl == "G")
		case "t", "T":
			eng.RotateGravity(gl == "T")
		case "x", "q", "X", "Q":
			return true
		}
	case wde.CloseEvent:
		return true
	}

	return false
}
