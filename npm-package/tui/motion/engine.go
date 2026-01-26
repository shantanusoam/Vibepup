package motion

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Nomadcxx/sysc-Go/animations"
)

type TickMsg time.Time

type Engine struct {
	Quiet   bool
	LowPerf bool
	Effect  animations.Animation
}

func New(lowPerf, quiet bool) Engine {
	// Initialize default effect (Matrix Rain)
	// sysc-Go animations usually need initialization
	// For now, we'll just placeholder it or use a simple one if available
	return Engine{
		Quiet:   quiet,
		LowPerf: lowPerf,
	}
}

func (e Engine) Next() tea.Cmd {
	d := time.Millisecond * 1000 / 60
	if e.LowPerf {
		d = time.Millisecond * 1000 / 15
	}
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// ASCII Dog Sprites
var DogSleeping = []string{
	"      z",
	"     z ",
	"    Z  ",
	"૮ – ﻌ – ა",
}

var DogRunning = []string{
	"            ",
	"            ",
	"            ",
	"૮ ˶ᵔ ᵕ ᵔ˶ ა", // Frame 1
	"            ",
	"            ",
	"            ",
	"૮ ˶• ﻌ •˶ ა", // Frame 2
}

var DogBarking = []string{
	"    WOOF!   ",
	"            ",
	"            ",
	"૮ ≧ ﻌ ≦ ა",
}

// GetDogFrame returns the current frame for the dog state
func GetDogFrame(state string, frame int) string {
	switch state {
	case "sleeping":
		// Animate Zzz
		zzz := []string{"", "z", "zz", "zzz"}
		return "      " + zzz[(frame/30)%4] + "\n" + DogSleeping[3]
	case "barking":
		if (frame/10)%2 == 0 {
			return DogBarking[0] + "\n" + DogBarking[3]
		}
		return "            \n" + DogBarking[3]
	case "happy":
		// Jump animation
		if (frame/5)%2 == 0 {
			return "            \n" + "૮ ˶ᵔ ᵕ ᵔ˶ ა" // Up
		}
		return "            \n" + "૮ ˶• ﻌ •˶ ა" // Down
	default: // Running
		frames := []string{"૮ ˶ᵔ ᵕ ᵔ˶ ა", "૮ ˶• ﻌ •˶ ა"}
		return "            \n" + frames[(frame/10)%len(frames)]
	}
}
