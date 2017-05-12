package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/markTward/gocloud-cicd/travis/config"
)

var configFile, buildTag, event, branch, baseImage, pr *string

func init() {
	const (
		defaultConfigFile = "./cicd.yaml"
		configFileUsage   = "configuration file containing project workflow values"
		buildTagUsage     = "existing image tag used as basis for further tags (required)"
		eventUsage        = "build event type from list: push, pull_request"
		branchTypeUsage   = "build branch (required)"
		prUsage           = "pull request number (required when event type is pull_request)"
		baseImageUsage    = "built image used as basis for tagging (required)"
	)
	baseImage = flag.String("image", "", baseImageUsage)
	configFile = flag.String("config", defaultConfigFile, configFileUsage)
	buildTag = flag.String("tag", "", buildTagUsage)
	event = flag.String("event", "push", eventUsage)
	branch = flag.String("branch", "", branchTypeUsage)
	pr = flag.String("pr", "", prUsage)
}

func makeTagList(repoURL string, refImage string, event string, branch string, pr string) (images []string, err error) {

	log.Println("Tagger args:", repoURL, refImage, event, branch, pr)

	// tag additional images based on build event type
	tagSep := strings.Index(refImage, ":")
	commitImage := repoURL + refImage[tagSep:]

	images = append(images, commitImage)

	switch event {
	case "push":
		images = append(images, repoURL+":"+branch)
		if branch == "master" {
			images = append(images, repoURL+":latest")
		}
	case "pull_request":
		images = append(images, repoURL+":PR-"+pr)
	}

	return images, err
}

func tagImages(src string, targets []string) (err error) {
	var stderr bytes.Buffer

	for _, target := range targets {
		cmd := exec.Command("docker", "tag", src, target)
		cmd.Stderr = &stderr
		log.Printf("docker tag from %v to %v", src, target)

		if err = cmd.Run(); err != nil {
			err = fmt.Errorf("%v", stderr.String())
			break
		}
	}

	return err
}

func validateCLInput() (err error) {

	if *baseImage == "" {
		err = fmt.Errorf("%v\n", "build image a required value; use --image option")
	}

	if *buildTag == "" {
		err = fmt.Errorf("%v\n", "build tag a required value; use --tag option")
	}

	if *branch == "" {
		err = fmt.Errorf("%v\n", "build branch a required value; use --branch option")
	}

	switch *event {
	case "push", "pull_request":
	default:
		err = fmt.Errorf("%v\n", "event type must be one of: push, pull_request")
	}

	if *event == "pull_request" && *pr == "" {
		err = fmt.Errorf("%v\n", "event type pull_request requires a PR number; use --pr option")
	}
	return err
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

	// parse and validate CLI
	flag.Parse()

	if err := validateCLInput(); err != nil {
		exitScript(err, true)
	}

	log.Printf("flag values: --config %v, --tag %v, -branch %v, --image %v, --event %v, --pr %v\n",
		*configFile, *buildTag, *branch, *baseImage, *event, *pr)

	// initialize configuration object
	cfg := config.New()
	if err := config.Load(*configFile, &cfg); err != nil {
		exitScript(err, true)
	}

	// initialize active registry indicated by config
	var activeRegistry interface{}
	var err error
	if activeRegistry, err = cfg.GetActiveRegistry(); err != nil {
		exitScript(err, true)
	}
	ar := activeRegistry.(config.Registrator)

	// validate registry has required values
	if err := ar.IsRegistryValid(); err != nil {
		exitScript(err, true)
	}

	// authenticate credentials for registry
	if err := ar.Authenticate(); err != nil {
		exitScript(err, true)
	}

	// make list of images to tag
	var images []string
	if images, err = makeTagList(ar.GetRepoURL(), *baseImage, *event, *branch, *pr); err != nil {
		exitScript(err, true)
	}

	// tag images
	if err := tagImages(*baseImage, images); err != nil {
		exitScript(err, true)
	}
	log.Println("tagged images:", images)

	// push images
	var result []string
	if result, err = ar.Push(images); err != nil {
		exitScript(err, true)
	}
	log.Println("pushed images:", result)

}
