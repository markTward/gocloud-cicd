package cicd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func IsDryRun(ctx *cobra.Command, wf *Workflow) bool {
	// return wf.Config.Dryrun
	return viper.GetBool("dryrun") || wf.Config.Dryrun
}

func IsDebug(ctx *cobra.Command, wf *Workflow) bool {
	return viper.GetBool("debug") || wf.Config.Debug
}

func LogError(err error) {
	log.Printf("error: %v\n", strings.TrimSpace(err.Error()))
}

func LogDebug(ctx *cobra.Command, s string) {
	if viper.GetBool("debug") {
		log.Printf("debug: %v\n", strings.TrimSpace(s))
	}
}

// func GetAllFlags(ctx *cobra.Command) map[string]map[string]string {
//
// 	// collection for all global and user assigned flags
// 	allFlags := make(map[string]map[string]string)
//
// 	// get global flags
// 	globalFlags := make(map[string]string)
// 	for i := 0; i < len(ctx.GlobalFlagNames()); i++ {
// 		flag := ctx.GlobalFlagNames()[i]
// 		value := ctx.GlobalString(flag)
// 		globalFlags[flag] = value
// 	}
// 	allFlags["global"] = globalFlags
//
// 	// get user assigned flags
// 	userFlags := make(map[string]string)
// 	for i := 0; i < len(ctx.FlagNames()); i++ {
// 		flag := ctx.FlagNames()[i]
// 		value := ctx.String(flag)
// 		userFlags[flag] = value
// 	}
// 	allFlags["user"] = userFlags
//
// 	return allFlags
// }
