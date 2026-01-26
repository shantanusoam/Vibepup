package animations

import "time"

var (
	dogeFrames = []string{
		"wow much wait",
		"such load very vibe",
		"plz hold pupper",
	}
	shrekFrames = []string{
		"  ／|、",
		" (°､ ｡ 7",
		" | 、`\\",
		" じしf_, )ノ",
	}
	catFrames = []string{
		"/ᐠ. ｡.ᐟ\\", // cat face
		"/ᐠ｡‸｡ᐟ\\",
		"/ᐠ – ᆽ – ᐟ\\",
	}
	wojakFrames = []string{
		"(・_・;)",
		"(・_・`)",
		"(・_・;)",
		"(・_・`)",
	}
)

func init() {
	Register(Preset{Name: "doge-wow", Kind: Event, Frames: dogeFrames, Interval: 110 * time.Millisecond, Density: 1})
	Register(Preset{Name: "shrek-blink", Kind: Event, Frames: shrekFrames, Interval: 140 * time.Millisecond, Density: 1})
	Register(Preset{Name: "cat-bounce", Kind: Event, Frames: catFrames, Interval: 90 * time.Millisecond, Density: 1})
	Register(Preset{Name: "wojak-stare", Kind: Event, Frames: wojakFrames, Interval: 120 * time.Millisecond, Density: 1})
}
