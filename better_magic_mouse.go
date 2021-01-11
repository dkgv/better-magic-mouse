package main

import (
	"math"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

type mousePoint struct {
	x int16
	y int16
}

func (p mousePoint) distanceTo(other mousePoint) float64 {
	dx := math.Abs(float64(other.x) - float64(p.x))
	dy := math.Abs(float64(other.y) - float64(p.y))
	return math.Sqrt(dx*dx + dy*dy)
}

// Amount of millis to hold left button before right clicking
const MsPressThreshold = 250
const LeftButton = 1

var pressPoint mousePoint
var pressStart int64 = 0

func main() {
	s := hook.Start()
	defer hook.End()

	for event := range s {
		if event.Button != LeftButton {
			continue
		}

		switch event.Kind {
		case hook.MouseHold:
			handleMouseHold(event)

		case hook.MouseDown:
			// Mouse was released
			pressStart = 0
		}
	}
}

func handleMouseHold(event hook.Event) {
	pressStart = nanoToMs(event.When.UnixNano())
	pressPoint = mousePoint{event.X, event.Y}

	go func() {
		time.Sleep(MsPressThreshold * time.Millisecond)

		if pressStart <= 0 {
			return
		}

		var pressPeriod = nanoToMs(time.Now().UnixNano()) - pressStart

		tryRightClick(pressPeriod)
	}()
}

func currMouse() mousePoint {
	x, y := robotgo.GetMousePos()
	return mousePoint{int16(x), int16(y)}
}

func tryRightClick(pressPeriod int64) {
	isDragging := pressPoint.distanceTo(currMouse()) > 15

	if isDragging || pressPeriod < MsPressThreshold {
		return
	}

	robotgo.Click("right", false)
}

func nanoToMs(time int64) int64 {
	return time / 1e6
}
