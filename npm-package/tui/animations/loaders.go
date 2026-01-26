package animations

import "time"

var (
	vhsFrames = []string{
		"▞▚▞▚▞▚▞▚▞▚▞▚",
		"▚▞▚▞▚▞▚▞▚▞▚▞",
		"▞▚▞▚▞▚▞▚▞▚▞▚",
	}
	crtFrames = []string{
		"│││││││││││││",
		"┃┃┃┃┃┃┃┃┃┃┃┃┃",
		"║║║║║║║║║║║║║",
	}
	matrixFrames = []string{
		"ａａｂｂ０１０１ｚｚ",
		"０１０１ｚｚａａｂｂ",
		"ｚｚａａｂｂ０１０１",
	}
	slimeFrames = []string{
		"(o˶╹︿╹˶o)",
		"(o˶╹﹏╹˶o)",
		"(o˶╹︿╹˶o)~",
		"~(o˶╹︿╹˶o)",
	}
	floppyFrames = []string{
		"💾",
		"💽",
		"💿",
	}
	waveFrames = []string{
		"~    ~    ~",
		"  ~    ~   ",
		"    ~    ~  ",
		" ~    ~    ~",
	}
	fireworkFrames = []string{
		"  .  ",
		" .*. ",
		".*★*",
		" .*. ",
		"  '  ",
	}
)

func init() {
	Register(Preset{Name: "vhs-scan", Kind: Loader, Frames: vhsFrames, Interval: 70 * time.Millisecond, Density: 3})
	Register(Preset{Name: "crt-wipe", Kind: Loader, Frames: crtFrames, Interval: 60 * time.Millisecond, Density: 2})
	Register(Preset{Name: "matrix-rain", Kind: Loader, Frames: matrixFrames, Interval: 90 * time.Millisecond, Density: 2})
	Register(Preset{Name: "slime-bounce", Kind: Loader, Frames: slimeFrames, Interval: 80 * time.Millisecond, Density: 1})
	Register(Preset{Name: "floppy-spin", Kind: Loader, Frames: floppyFrames, Interval: 80 * time.Millisecond, Density: 1})
	Register(Preset{Name: "vibe-wave", Kind: Idle, Frames: waveFrames, Interval: 120 * time.Millisecond, Density: 1})
	Register(Preset{Name: "fireworks", Kind: Event, Frames: fireworkFrames, Interval: 100 * time.Millisecond, Density: 2})
}
