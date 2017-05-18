package commands

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/markTward/gocloud-cicd/config"
	"github.com/urfave/cli"
)

var event, baseImage, pr string

var PushCmd = cli.Command{
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
		cli.BoolFlag{
			Name:        "dryrun",
			Usage:       "log output but do not execute",
			Destination: &dryrun,
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

	if err := validatePushArgs(); err != nil {
		LogError(err)
		return err
	}
	log.Println("push command args:", getAllFlags(ctx))

	// initialize configuration object
	cfg := config.New()
	if err := config.Load(configFile, &cfg); err != nil {
		LogError(err)
		return err
	}

	LogDebug(ctx, fmt.Sprintf("%v", spew.Sdump(cfg)))

	// initialize active Registry indicated by config and assert as Registrator
	var activeRegistry interface{}
	var err error
	if activeRegistry, err = cfg.GetActiveRegistry(); err != nil {
		LogError(err)
		return err
	}
	ar := activeRegistry.(config.Registrator)

	// validate registry has required values
	if err := ar.IsRegistryValid(); err != nil {
		LogError(err)
		return err
	}

	// authenticate credentials for registry
	if err := ar.Authenticate(); err != nil {
		LogError(err)
		return err
	}

	// make list of images to tag
	var images []string
	if images = makeTagList(ctx, ar.GetRepoURL(), baseImage, event, branch, pr); len(images) == 0 {
		fmt.Errorf("no images to tag: ", images)
		LogError(err)
		return err
	}

	// tag images
	if err := tagImages(baseImage, images); err != nil {
		LogError(err)
		return err
	}
	log.Println("tagged images:", images)

	// push tagged images
	var result []string
	if result, err = ar.Push(images, dryrun); err != nil {
		LogError(err)
		return err
	}
	log.Println("pushed images:", result)
	return err
}

func makeTagList(ctx *cli.Context, repoURL string, refImage string, event string, branch string, pr string) (images []string) {

	LogDebug(ctx, fmt.Sprintf("makeTagList args: repo url: %v, image: %v, event type: %v, branch: %v, pull request id: %v", repoURL, refImage, event, branch, pr))

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

func validatePushArgs() (err error) {
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
