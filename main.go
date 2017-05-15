package main

import (
	"os"

	"github.com/markTward/gocloud-cicd/commands"
	"gopkg.in/urfave/cli.v1"
)

func main() {

	var debug bool
	var verbose bool

	app := cli.NewApp()
	app.Name = "CICD Tools"
	app.Usage = "Continuous Intergration and Deployment Tools"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "verbose",
			Usage:       "Show more output",
			Destination: &verbose,
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Show detailed debugging output",
			Destination: &debug,
		},
	}
	app.Commands = []cli.Command{
		commands.DeployCmd,
		commands.PushCmd,
	}

	app.Run(os.Args)
}
