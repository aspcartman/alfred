package env

import (
	"github.com/aspcartman/darkside"
	"github.com/aspcartman/darkside/g"
	"fmt"
)

func init() {
	darkside.SetUnrecoveredPanicHandler(func(p *g.Panic) {
		fmt.Println("the fuck?")
	})
}