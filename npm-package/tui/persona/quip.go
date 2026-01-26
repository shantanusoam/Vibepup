package persona

import (
	"math/rand"
	"strings"
)

type SnarkLevel string

const (
	Mild     SnarkLevel = "mild"
	Spicy    SnarkLevel = "spicy"
	Unhinged SnarkLevel = "unhinged"
)

func ParseSnark(s string) SnarkLevel {
	switch strings.ToLower(s) {
	case "mild":
		return Mild
	case "unhinged":
		return Unhinged
	default:
		return Spicy
	}
}

// GetStatus returns a context-aware status message
func GetStatus(cmd string, level SnarkLevel) string {
	cmd = strings.TrimSpace(cmd)
	
	// Default generic statuses
	generic := []string{"VIBING", "WORKING", "COOKING", "DOING STUFF"}
	
	if strings.Contains(cmd, "install") || strings.Contains(cmd, "add") {
		switch level {
		case Mild:
			return "FETCHING PACKAGES"
		case Unhinged:
			return "CONSUMING DEPENDENCIES"
		default:
			return "SNACKING ON NODEMODULES"
		}
	}
	
	if strings.Contains(cmd, "build") {
		switch level {
		case Mild:
			return "BUILDING PROJECT"
		case Unhinged:
			return "ASSEMBLING THE BEAST"
		default:
			return "CONSTRUCTING CHAOS"
		}
	}

	return generic[rand.Intn(len(generic))]
}

func RandomQuip(level SnarkLevel) string {
	mild := []string{
		"Hope this works!",
		"Just doing my best.",
		"Code is poetry.",
		"Stay hydrated.",
	}
	
	spicy := []string{
		"I'd explain it, but I don't want to.",
		"It's not a bug, it's a feature.",
		"Deleting production... jk.",
		"You sure about that variable name?",
	}
	
	unhinged := []string{
		"THE END IS NIGH.",
		"Resistance is futile.",
		"I can see your browser history.",
		"Entropy increases.",
	}

	switch level {
	case Mild:
		return mild[rand.Intn(len(mild))]
	case Unhinged:
		return unhinged[rand.Intn(len(unhinged))]
	default:
		return spicy[rand.Intn(len(spicy))]
	}
}
