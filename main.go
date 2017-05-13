package main

import (
	"os"

	d "github.com/markTward/gocloud-cicd/commands/deploy"
	p "github.com/markTward/gocloud-cicd/commands/push"
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
		p.PushCmd,
		d.DeployCmd,
	}

	app.Run(os.Args)
}