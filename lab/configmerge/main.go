package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	App
	Github
	Workflow
	Registry
}

type App struct {
	Name string
}
type Github struct {
	Repo string
}

type Registry struct {
	GCR
	Docker
}

type GCR struct {
	Name        string
	Description string
	Host        string
	Project     string
	Repo        string
	Url         string
	KeyFile     string
}

type Docker struct {
	Name        string
	Description string
	Host        string
	Account     string
	Repo        string
	Url         string
}

type Workflow struct {
	Enabled bool

	Github struct {
		Repo   string
		Branch string
	}

	CIProvider struct {
		Name string
		Plan string
	}

	Platform struct {
		Name    string
		Project string
		Cluster string
	}

	Registry string

	CDProvider struct {
		Name      string
		Release   string
		Namespace string
		ChartDir  string
	}
}

func exitScript(err error, exit bool) {
	s := strings.TrimSpace(err.Error())
	log.Printf("error: %v", s)
	if exit {
		fmt.Fprintf(os.Stderr, "error: %v\n", s)
		os.Exit(1)
	}
}

func main() {
	// read in project config file
	yamlInput, err := ioutil.ReadFile("./cicddefault.yaml")
	if err != nil {
		exitScript(err, true)
	}

	// parse yaml into Config object
	cfg := Config{}
	err = yaml.Unmarshal([]byte(yamlInput), &cfg)
	if err != nil {
		exitScript(err, true)
	}

	fmt.Println("Default APP:", cfg.App.Name)
	yamlInput, err = ioutil.ReadFile("./cicduser.yaml")
	if err != nil {
		exitScript(err, true)
	}

	err = yaml.Unmarshal([]byte(yamlInput), &cfg)
	if err != nil {
		exitScript(err, true)
	}

	fmt.Println("User APP:", cfg.App.Name)
	fmt.Printf("Config: %v\n", cfg)
}
