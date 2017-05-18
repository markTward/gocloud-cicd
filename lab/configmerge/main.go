package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/markTward/gocloud-cicd/config"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	// read in project config file
	yamlInput, err := ioutil.ReadFile("../../../gocloud/cicd.yaml")
	if err != nil {
		exitScript(err, true)
	}

	// parse yaml into Config object
	cfg := config.Config{}
	err = yaml.Unmarshal([]byte(yamlInput), &cfg)
	if err != nil {
		exitScript(err, true)
	}

	// spew.Dump(cfg)
	// log.Println()
	log.Println(fmt.Sprintf("Config: %v", spew.Sdump(cfg)))

	// fmt.Println("Default APP:", cfg.App.Name)
	// yamlInput, err = ioutil.ReadFile("./cicduser.yaml")
	// if err != nil {
	// 	exitScript(err, true)
	// }
	//
	// err = yaml.Unmarshal([]byte(yamlInput), &cfg)
	// if err != nil {
	// 	exitScript(err, true)
	// }
	//
	// fmt.Println("User APP:", cfg.App.Name)
	// fmt.Printf("Config: %v\n", cfg)
}

func exitScript(err error, exit bool) {
	s := strings.TrimSpace(err.Error())
	log.Printf("error: %v", s)
	if exit {
		fmt.Fprintf(os.Stderr, "error: %v\n", s)
		os.Exit(1)
	}
}
