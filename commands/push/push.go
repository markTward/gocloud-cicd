package push

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	cmd "github.com/markTward/gocloud-cicd/commands"
	"github.com/markTward/gocloud-cicd/config"
	"gopkg.in/urfave/cli.v1"
)

var configFile, event, branch, baseImage, pr string
var dryrun bool

var PushCmd = cli.Command{
	Name:  "push",
	Usage: "push images to repository (gcr or docker)",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "image, i",
			Usage:       "built image used as basis for tagging (required)",
			Destination: &baseImage,
		},
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "configuration file containing project workflow values",
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
			Name:        "pr",
			Usage:       "pull request number (required when event type is pull_request)",
			Destination: &pr,
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

	cmd.LogDebug(c, fmt.Sprintf("flag values: --config %v, --branch %v, --image %v, --event %v, --pr %v --debug %v, --dryrun %v",
		configFile, branch, baseImage, event, pr, c.GlobalBool("debug"), dryrun))

	if err := validateCLInput(c); err != nil {
		cmd.LogError(err)
		return err
	}

	// initialize configuration object
	cfg := config.New()
	if err := config.Load(configFile, &cfg); err != nil {
		cmd.LogError(err)
		return err
	}

	cmd.LogDebug(c, fmt.Sprintf("Config: %#v", cfg))

	// initialize active registry indicated by config
	var activeRegistry interface{}
	var err error
	if activeRegistry, err = cfg.GetActiveRegistry(); err != nil {
		cmd.LogError(err)
		return err
	}
	ar := activeRegistry.(config.Registrator)

	// validate registry has required values
	if err := ar.IsRegistryValid(); err != nil {
		cmd.LogError(err)
		return err
	}

	// authenticate credentials for registry
	if err := ar.Authenticate(); err != nil {
		cmd.LogError(err)
		return err
	}

	// make list of images to tag
	var images []string
	if images, err = makeTagList(ar.GetRepoURL(), baseImage, event, branch, pr); err != nil {
		cmd.LogError(err)
		return err
	}

	// tag images
	if err := tagImages(baseImage, images); err != nil {
		cmd.LogError(err)
		return err
	}
	log.Println("tagged images:", images)

	// TODO: add global flags debug, verbose and dryrun as flags args
	// push images
	var result []string
	if result, err = ar.Push(images); err != nil {
		cmd.LogError(err)
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

func validateCLInput(c *cli.Context) (err error) {
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
