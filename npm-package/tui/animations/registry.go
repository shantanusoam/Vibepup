package animations

import "time"

type Kind string

const (
	Loader Kind = "loader"
	Idle   Kind = "idle"
	Event  Kind = "event"
)

type Preset struct {
	Name     string
	Kind     Kind
	Frames   []string
	Interval time.Duration
	Density  int
}

var presets = map[string]Preset{}

func Register(p Preset) {
	if p.Interval == 0 {
		p.Interval = time.Millisecond * 80
	}
	if p.Density == 0 {
		p.Density = 1
	}
	presets[p.Name] = p
}

func Get(name string) Preset {
	if p, ok := presets[name]; ok {
		return p
	}
	return presets["vhs-scan"]
}

func All() []Preset {
	out := make([]Preset, 0, len(presets))
	for _, p := range presets {
		out = append(out, p)
	}
	return out
}

func Frame(p Preset, index int) (string, int) {
	if len(p.Frames) == 0 {
		return "", index
	}
	next := (index + 1) % len(p.Frames)
	return p.Frames[index%len(p.Frames)], next
}
