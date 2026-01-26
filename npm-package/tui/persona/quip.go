package persona

import (
	"math/rand"
	"strings"
)

type SnarkLevel int

const (
	Mild SnarkLevel = iota
	Spicy
	Unhinged
)

type State string

const (
	StateWaiting State = "waiting"
	StateSuccess State = "success"
	StateError   State = "error"
	StateIdle    State = "idle"
)

type CopyPack map[State][]string

var defaultCopy = map[SnarkLevel]CopyPack{
	Mild: {
		StateWaiting: {
			"Loading... promise I'm not doomscrolling.",
			"Buffering vibes...",
		},
		StateSuccess: {
			"Done. Sparkles delivered.",
			"Shipped. No take-backs.",
		},
		StateError: {
			"Oopsie. Let's pretend that didn't happen.",
			"Error? Never heard of her.",
		},
		StateIdle: {
			"Idle mode. Hydrate, maybe?",
			"Waiting for chaos to resume...",
		},
	},
	Spicy: {
		StateWaiting: {
			"Cooking the bits. Might burn them a little.",
			"Hold up—optimizing my snark cache...",
		},
		StateSuccess: {
			"Boom. Pixel-perfect flex.",
			"Mission accomplished, no body count (this time).",
		},
		StateError: {
			"I swear that compiled in my head.",
			"Stacktrace says it's your fault. Kidding. Mostly.",
		},
		StateIdle: {
			"Idle. Manifesting a raise.",
			"BRB, updating my LinkedIn to 'professional gremlin'.",
		},
	},
	Unhinged: {
		StateWaiting: {
			"Buffering the chaos. Please enjoy this existential dread.",
			"Spinning up demons... I mean daemons.",
		},
		StateSuccess: {
			"Shipped. If it breaks, that's a feature drop.",
			"Done. Tell compliance I was never here.",
		},
		StateError: {
			"Error 500: vibes deceased.",
			"It exploded. Deploy to prod?",
		},
		StateIdle: {
			"Idle. Practicing my villain arc.",
			"Loading memes from forbidden archives...",
		},
	},
}

func ParseSnark(level string) SnarkLevel {
	switch strings.ToLower(level) {
	case "spicy":
		return Spicy
	case "unhinged":
		return Unhinged
	default:
		return Mild
	}
}

func Quip(state State, level SnarkLevel, packs map[SnarkLevel]CopyPack) string {
	pack := defaultCopy[level]
	if packs != nil {
		if custom, ok := packs[level]; ok {
			pack = custom
		}
	}
	lines := pack[state]
	if len(lines) == 0 {
		return ""
	}
	return lines[rand.Intn(len(lines))]
}
