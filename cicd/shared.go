package cicd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func IsDryRun(ctx *cobra.Command, wf *Workflow) bool {
	return viper.GetBool("isDryRun") || wf.Config.Dryrun
}

func IsDebug(ctx *cobra.Command, wf *Workflow) bool {
	return viper.GetBool("isDebug") || wf.Config.Debug
}

func LogError(err error) {
	log.Printf("error: %v\n", strings.TrimSpace(err.Error()))
}

func LogDebug(ctx *cobra.Command, s string) {
	if viper.GetBool("isDebug") {
		log.Printf("debug: %v\n", strings.TrimSpace(s))
	}
}
