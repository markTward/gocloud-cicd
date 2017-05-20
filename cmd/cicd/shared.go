package main

import (
	"log"
	"strings"

	"github.com/markTward/gocloud-cicd/cicd"
	cli "github.com/urfave/cli"
)

// vars shared by multiple commands
var configFile, branch string

// utility functions
func logError(err error) {
	log.Printf("error: %v\n", strings.TrimSpace(err.Error()))
}

func logDebug(ctx *cli.Context, s string) {
	if ctx.GlobalBool("debug") {
		log.Printf("debug: %v\n", strings.TrimSpace(s))
	}
}

func getAllFlags(ctx *cli.Context) map[string]map[string]string {

	// collection for all global and user assigned flags
	allFlags := make(map[string]map[string]string)

	// get global flags
	globalFlags := make(map[string]string)
	for i := 0; i < len(ctx.GlobalFlagNames()); i++ {
		flag := ctx.GlobalFlagNames()[i]
		value := ctx.GlobalString(flag)
		globalFlags[flag] = value
	}
	allFlags["global"] = globalFlags

	// get user assigned flags
	userFlags := make(map[string]string)
	for i := 0; i < len(ctx.FlagNames()); i++ {
		flag := ctx.FlagNames()[i]
		value := ctx.String(flag)
		userFlags[flag] = value
	}
	allFlags["user"] = userFlags

	return allFlags
}

func isDryRun(ctx *cli.Context, wf *cicd.Workflow) bool {
	return ctx.GlobalBool("dryrun") || wf.Config.Dryrun
}

func isDebug(ctx *cli.Context, wf *cicd.Workflow) bool {
	return ctx.GlobalBool("debug") || wf.Config.Debug
}
