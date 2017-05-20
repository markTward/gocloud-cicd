package cicd

import (
	"log"
	"strings"

	"github.com/urfave/cli"
)

func IsDryRun(ctx *cli.Context, wf *Workflow) bool {
	return ctx.GlobalBool("dryrun") || wf.Config.Dryrun
}

func IsDebug(ctx *cli.Context, wf *Workflow) bool {
	return ctx.GlobalBool("debug") || wf.Config.Debug
}

func LogError(err error) {
	log.Printf("error: %v\n", strings.TrimSpace(err.Error()))
}

func LogDebug(ctx *cli.Context, s string) {
	if ctx.GlobalBool("debug") {
		log.Printf("debug: %v\n", strings.TrimSpace(s))
	}
}

func GetAllFlags(ctx *cli.Context) map[string]map[string]string {

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
