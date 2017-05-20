package main

import (
	"os"
	"sort"

	"github.com/urfave/cli"
)

var configFile, branch string
var debug, dryrun bool

func main() {

	app := cli.NewApp()
	app.Name = "CICD Tools"
	app.Usage = "Continuous Intergration and Deployment Tools"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "dryrun",
			Usage:       "Show command output without execution",
			Destination: &dryrun,
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
