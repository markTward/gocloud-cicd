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

var configFile, buildTag, event, branch, baseImage, pr string
var dryrun bool

var PushCmd = cli.Command{
	Name:  "push",
	Usage: "push images to repository",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "image",
			Usage:       "built image used as basis for tagging (required)",
			Destination: &baseImage,
		},
		cli.StringFlag{
			Name:        "config",
			Usage:       "configuration file containing project workflow values",
			Destination: &configFile,
		},
		cli.StringFlag{
			Name:        "event",
			Usage:       "build event type from list: push, pull_request",
			Destination: &event,
		},
		cli.StringFlag{
			Name:        "pr",
			Usage:       "pull request number (required when event type is pull_request)",
			Destination: &pr,
		},
		cli.StringFlag{
			Name:        "tag, t",
			Usage:       "existing image tag used as basis for further tags (required)",
			Destination: &buildTag,
		},
		cli.StringFlag{
			Name:        "branch, b",
			Usage:       "build branch (required)",
			Destination: &branch,
		},
		cli.BoolFlag{
			Name:        "dryrun",
			Usage:       "log output but do not execute",
			Destination: &dryrun,
		},
	},
	Action: push,
}

func push(c *cli.Context) error {

	if debug(c) {
		log.Printf("flag values: --config %v, --tag %v, --branch %v, --image %v, --event %v, --pr %v --debug %v, --verbose %v\n",
			configFile, buildTag, branch, baseImage, event, pr, c.GlobalBool("debug"), c.GlobalBool("verbose"))
	}

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

func debug(c *cli.Context) bool {
	return c.GlobalBool("debug")
}
