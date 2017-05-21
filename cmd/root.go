// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/markTward/gocloud-cicd/cicd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var dryrun, debug bool
var wf *cicd.Workflow

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cicd",
	Short: "Continuous Intergration and Deployment Tools",
	Long:  "Continuous Intergration and Deployment Tools",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ./cicd.yaml)")
	RootCmd.PersistentFlags().BoolP("debug", "", false, "Show detailed debugging output")
	RootCmd.PersistentFlags().BoolP("dryrun", "", false, "Show command output without execution")

	viper.BindPFlag("dryrun", RootCmd.PersistentFlags().Lookup("dryrun"))
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// initialize configuration object
	wf = cicd.New()
	if err := cicd.Load(configFile, wf); err != nil {
		// cicd.LogError(err)
		// return err
		log.Println("initConfig err:", err)
	}
	cicd.LogDebug(RootCmd, fmt.Sprintf("Config: %v", spew.Sdump(wf)))

	// cicd.LogDebug(ctx, fmt.Sprintf("%v", spew.Sdump(wf)))
	log.Println("viper dryrun:", viper.GetBool("dryrun"))
}
