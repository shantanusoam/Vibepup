package config

import "flag"

type Flags struct {
	Quiet    bool
	NoEmoji  bool
	Dense    bool
	PerfLow  bool
	Snark    string
	Theme    string
	Anim     string
	FX       string
	NoAlt    bool
	ForceRun bool
}

func Parse() Flags {
	f := Flags{}
	flag.BoolVar(&f.Quiet, "quiet", false, "reduce motion and chatter")
	flag.BoolVar(&f.NoEmoji, "no-emoji", false, "disable emoji rendering")
	flag.BoolVar(&f.Dense, "dense", false, "increase animation density")
	flag.BoolVar(&f.PerfLow, "perf-low", false, "lower FPS and effects for slower terminals")
	flag.StringVar(&f.Snark, "snark", "mild", "snark level: mild|spicy|unhinged")
	flag.StringVar(&f.Theme, "theme", "dracula-vibe", "theme name")
	flag.StringVar(&f.Anim, "anim", "vhs-scan", "animation preset")
	flag.StringVar(&f.FX, "fx", "fire", "sysc effect: fire|matrix|none")
	flag.BoolVar(&f.NoAlt, "no-alt", true, "disable alt screen (stay in current terminal)")
	flag.BoolVar(&f.ForceRun, "force-run", false, "run child process even if stdout is not a TTY")
	flag.Parse()
	return f
}
