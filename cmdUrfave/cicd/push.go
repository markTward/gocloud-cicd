package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/markTward/gocloud-cicd"
	"github.com/urfave/cli"
)

var event, baseImage, pr string

var pushCmd = cli.Command{
	Name:  "push",
	Usage: "push images to repository (gcr or docker)",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "branch, b",
			Usage:       "build branch (required)",
			Destination: &branch,
		},
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "load configuration file from `FILE`",
			Value:       "./cicd.yaml",
			Destination: &configFile,
		},
		cli.StringFlag{
			Name:        "event, e",
			Usage:       "build event type from list: push, pull_request",
			Value:       "push",
			Destination: &event,
		},
		cli.StringFlag{
			Name:        "image, i",
			Usage:       "built image used as basis for tagging (required)",
			Destination: &baseImage,
		},
		cli.StringFlag{
			Name:        "pr",
			Usage:       "pull request number (required when event type is pull_request)",
			Destination: &pr,
		},
	},
	Action: push,
}

func push(ctx *cli.Context) error {

	// initialize configuration object
	wf := cicd.New()
	if err := cicd.Load(configFile, wf); err != nil {
		cicd.LogError(err)
		return err
	}

	if err := validatePushArgs(ctx, wf); err != nil {
		cicd.LogError(err)
		return err
	}
	log.Println("push command args:", cicd.GetAllFlags(ctx))

	cicd.LogDebug(ctx, fmt.Sprintf("%v", spew.Sdump(wf)))

	// initialize active Registry indicated by config and assert as Registrator
	var activeRegistry interface{}
	var err error
	if activeRegistry, err = wf.GetActiveRegistry(); err != nil {
		cicd.LogError(err)
		return err
	}
	ar := activeRegistry.(cicd.Registrator)

	// validate registry has required values
	if err := ar.IsRegistryValid(); err != nil {
		cicd.LogError(err)
		return err
	}

	// authenticate credentials for registry
	if err := ar.Authenticate(ctx, wf); err != nil {
		cicd.LogError(err)
		return err
	}

	// make list of images to tag
	var images []string
	if images = makeTagList(ctx, ar.GetRepoURL(), baseImage, event, branch, pr); len(images) == 0 {
		cicd.LogError(fmt.Errorf("no images to tag: %v", images))
		return err
	}

	// tag images
	if err := tagImages(baseImage, images); err != nil {
		cicd.LogError(err)
		return err
	}
	log.Println("tagged images:", images)

	// push tagged images
	var result []string
	if result, err = ar.Push(ctx, wf, images); err != nil {
		cicd.LogError(err)
		return err
	}
	log.Println("pushed images:", result)
	return err
}

func makeTagList(ctx *cli.Context, repoURL string, refImage string, event string, branch string, pr string) (images []string) {

	cicd.LogDebug(ctx, fmt.Sprintf("makeTagList args: repo url: %v, image: %v, event type: %v, branch: %v, pull request id: %v",
		repoURL, refImage, event, branch, pr))

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

	return images
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

func validatePushArgs(ctx *cli.Context, wf *cicd.Workflow) (err error) {

	// handle globals from cli and/or workflow config
	if cicd.IsDebug(ctx, wf) {
		debug = true
	}

	if cicd.IsDryRun(ctx, wf) {
		dryrun = true
	}

	switch {
	case baseImage == "":
		err = fmt.Errorf("%v", "build image a required value; use --image option")

	case branch == "":
		err = fmt.Errorf("%v", "build branch a required value; use --branch option")

	case !(event == "push" || event == "pull_request"):
		err = fmt.Errorf("%v", "event type must be one of: push, pull_request")

	case event == "pull_request" && pr == "":
		err = fmt.Errorf("%v", "event type pull_request requires a PR number; use --pr option")
	}
	return err
}
