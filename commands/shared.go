package commands

import (
	"log"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

// vars shared by multiple commands
var configFile, branch string
var dryrun bool

// utility functions
func LogError(err error) {
	log.Printf("error: %v\n", strings.TrimSpace(err.Error()))
}

func LogDebug(c *cli.Context, s string) {
	if c.GlobalBool("debug") {
		log.Printf("debug: %v\n", strings.TrimSpace(s))
	}
}
