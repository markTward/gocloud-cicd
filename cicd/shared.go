package cicd

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

func IsDryRun() bool {
	return viper.GetBool("isDryRun")
}

func IsDebug() bool {
	return viper.GetBool("isDebug")
}

func LogError(err error) {
	log.Printf("error: %v\n", strings.TrimSpace(err.Error()))
}

func LogDebug(s string) {
	if IsDebug() {
		log.Printf("debug: %v\n", strings.TrimSpace(s))
	}
}
