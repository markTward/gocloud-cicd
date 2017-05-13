package commands

import (
	"log"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

func LogError(err error) {
	s := strings.TrimSpace(err.Error())
	log.Printf("ERROR: %v", s)
}

func LogDebug(c *cli.Context, s string) {
	if c.GlobalBool("debug") {
		log.Printf("DEBUG: %v\n", strings.TrimSpace(s))
	}
}
