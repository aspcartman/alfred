package env

import (
	"github.com/aspcartman/darkside"
	"github.com/aspcartman/darkside/g"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"github.com/aspcartman/exceptions"
)

var Log logrus.FieldLogger

func init() {
	Log = &logrus.Logger{
		Out:       os.Stdout,
		Hooks:     logrus.LevelHooks{},
		Formatter: &logrus.TextFormatter{ForceColors: true, DisableTimestamp: true},
		Level:     logrus.DebugLevel,
	}

	e.RegisterHook(func(ex *e.Exception) {
		Log.WithError(ex.BottommostError()).Error(ex.Info)
	})

	darkside.SetUnrecoveredPanicHandler(func(p *g.Panic) {
		fmt.Println("the fuck?")
	})
}
