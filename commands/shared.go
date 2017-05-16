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
