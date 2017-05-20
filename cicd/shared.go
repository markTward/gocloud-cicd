package cicd

import "github.com/urfave/cli"

func isDryRun(ctx *cli.Context, wf *Workflow) bool {
	return ctx.GlobalBool("dryrun") || wf.Config.Dryrun
}

func isDebug(ctx *cli.Context, wf *Workflow) bool {
	return ctx.GlobalBool("debug") || wf.Config.Debug
}
