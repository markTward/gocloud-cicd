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

	RootCmd.PersistentFlags().StringVar(&configFile, "config", "./cicd.yaml", "config file (default is ./cicd.yaml)")
	RootCmd.PersistentFlags().BoolP("debug", "", false, "Show detailed debugging output")
	RootCmd.PersistentFlags().BoolP("dryrun", "", false, "Show command output without execution")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if configFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")    // adding home directory as first search path
		viper.SetConfigName("cicd") // name of config file (without extension)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
		err = viper.Unmarshal(&wf)
		if err != nil {
			log.Fatalf("unable to decode into struct: %v", err)
		}
	} else {
		log.Fatalf("unable to read config file: %v", err)
	}

	// override config file settings
	viper.SetDefault("isDryRun", wf.Config.Dryrun)
	if RootCmd.PersistentFlags().Lookup("dryrun").Changed {
		viper.BindPFlag("isDryRun", RootCmd.PersistentFlags().Lookup("dryrun"))
	}

	viper.SetDefault("isDebug", wf.Config.Debug)
	if RootCmd.PersistentFlags().Lookup("debug").Changed {
		viper.BindPFlag("isDebug", RootCmd.PersistentFlags().Lookup("debug"))
	}

	if viper.GetBool("isDryRun") {
		log.Println("operating in dryrun mode")
	}

	if viper.GetBool("isDebug") {
		log.Println("operating in debug mode")
	}

	cicd.LogDebug(fmt.Sprintf("Config: %v", spew.Sdump(wf)))

}
