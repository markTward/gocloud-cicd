package cicd

import "github.com/urfave/cli"

func isDryRun(ctx *cli.Context) bool {
	return ctx.Bool("dryrun")
}
