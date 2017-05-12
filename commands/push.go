package commands

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/markTward/gocloud-cicd/travis/config"
	"gopkg.in/urfave/cli.v1"
)

func push(c *cli.Context) error {

	log.Printf("flag values: --config %v, --tag %v, --branch %v, --image %v, --event %v, --pr %v --debug %v, --verbose %v\n",
		configFile, buildTag, branch, baseImage, event, pr, c.GlobalBool("debug"), c.GlobalBool("verbose"))

	if err := validateCLInput(); err != nil {
		logError(err)
		return err
	}

	// initialize configuration object
	cfg := config.New()
	if err := config.Load(configFile, &cfg); err != nil {
		logError(err)
		return err
	}

	// initialize active registry indicated by config
	var activeRegistry interface{}
	var err error
	if activeRegistry, err = cfg.GetActiveRegistry(); err != nil {
		logError(err)
		return err
	}
	ar := activeRegistry.(config.Registrator)

	// validate registry has required values
	if err := ar.IsRegistryValid(); err != nil {
		logError(err)
		return err
	}

	// authenticate credentials for registry
	if err := ar.Authenticate(); err != nil {
		logError(err)
		return err
	}

	// make list of images to tag
	var images []string
	if images, err = makeTagList(ar.GetRepoURL(), baseImage, event, branch, pr); err != nil {
		logError(err)
		return err
	}

	// tag images
	if err := tagImages(baseImage, images); err != nil {
		logError(err)
		return err
	}
	log.Println("tagged images:", images)

	// push images
	var result []string
	if result, err = ar.Push(images); err != nil {
		logError(err)
		return err
	}
	log.Println("pushed images:", result)
	return err
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

	if baseImage == "" {
		err = fmt.Errorf("%v", "build image a required value; use --image option")
	}

	if buildTag == "" {
		err = fmt.Errorf("%v", "build tag a required value; use --tag option")
	}

	if branch == "" {
		err = fmt.Errorf("%v", "build branch a required value; use --branch option")
	}

	switch event {
	case "push", "pull_request":
	default:
		err = fmt.Errorf("%v", "event type must be one of: push, pull_request")
	}

	if event == "pull_request" && pr == "" {
		err = fmt.Errorf("%v", "event type pull_request requires a PR number; use --pr option")
	}
	return err
}

func logError(err error) {
	s := strings.TrimSpace(err.Error())
	log.Printf("error: %v", s)
}
