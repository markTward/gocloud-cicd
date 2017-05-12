package commands

import (
	"gopkg.in/urfave/cli.v1"
)

//globals
var verbose bool

// push
var configFile, buildTag, event, branch, baseImage, pr string
var dryrun bool

var Commands = []cli.Command{
	{
		Name:  "push",
		Usage: "push images to repository",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "image",
				Usage:       "built image used as basis for tagging (required)",
				Destination: &baseImage,
			},
			cli.StringFlag{
				Name:        "config",
				Usage:       "configuration file containing project workflow values",
				Destination: &configFile,
			},
			cli.StringFlag{
				Name:        "event",
				Usage:       "build event type from list: push, pull_request",
				Destination: &event,
			},
			cli.StringFlag{
				Name:        "pr",
				Usage:       "pull request number (required when event type is pull_request)",
				Destination: &pr,
			},
			cli.StringFlag{
				Name:        "tag, t",
				Usage:       "existing image tag used as basis for further tags (required)",
				Destination: &buildTag,
			},
			cli.StringFlag{
				Name:        "branch, b",
				Usage:       "build branch (required)",
				Destination: &branch,
			},
			cli.BoolFlag{
				Name:        "dryrun",
				Usage:       "log output but do not execute",
				Destination: &dryrun,
			},
		},
		Action: push,
	},
}
