package main

import (
	"os"
	"sort"

	"github.com/urfave/cli"
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
		deployCmd,
		pushCmd,
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.FlagsByName(deployCmd.Flags))
	sort.Sort(cli.FlagsByName(pushCmd.Flags))

	app.Run(os.Args)

}
