package commands

import (
	"log"
	"strings"

	cli "github.com/urfave/cli"
)

// vars shared by multiple commands
var configFile, branch string
var dryrun bool

// utility functions
func LogError(err error) {
	log.Printf("error: %v\n", strings.TrimSpace(err.Error()))
}

func LogDebug(ctx *cli.Context, s string) {
	if ctx.GlobalBool("debug") {
		log.Printf("debug: %v\n", strings.TrimSpace(s))
	}
}

func getAllFlags(ctx *cli.Context) map[string]string {
	flagValues := make(map[string]string)
	for i := 0; i < len(ctx.GlobalFlagNames()); i++ {
		flag := ctx.GlobalFlagNames()[i]
		value := ctx.GlobalString(flag)
		flagValues[flag] = value
	}

	for i := 0; i < len(ctx.FlagNames()); i++ {
		flag := ctx.FlagNames()[i]
		value := ctx.String(flag)
		flagValues[flag] = value
	}
	return flagValues
}
