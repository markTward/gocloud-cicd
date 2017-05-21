package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/markTward/gocloud-cicd/cicd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var branch, event, baseImage, pr string

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push a list of containers to a registry",
	Long:  "push a list of containers to a registry",
	RunE:  push,
}

func init() {
	pushCmd.Flags().StringVarP(&branch, "branch", "b", "", "branch name for tagging")
	pushCmd.Flags().StringVarP(&event, "event", "e", "push", "build event type from list: push, pull_request")
	pushCmd.Flags().StringVarP(&baseImage, "image", "i", "", "built image used as basis for tagging (required)")
	pushCmd.Flags().StringVarP(&pr, "pr", "", "", "pull request number (required when event type is pull_request)")

	viper.BindPFlag("branch", pushCmd.Flags().Lookup("branch"))
	viper.BindPFlag("url", pushCmd.PersistentFlags().Lookup("dryrun"))

	RootCmd.AddCommand(pushCmd)

}

func push(ccmd *cobra.Command, args []string) (err error) {

	// validate args
	if err := validatePushArgs(); err != nil {
		cicd.LogError(err)
		return err
	}

	// initialize active Registry indicated by config and assert as Registrator
	var activeRegistry interface{}
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
	if err := ar.Authenticate(wf); err != nil {
		cicd.LogError(err)
		return err
	}

	// make list of images to tag
	var images []string
	if images = makeTagList(ar.GetRepoURL()); len(images) == 0 {
		cicd.LogError(fmt.Errorf("no images to tag: %v", images))
		return err
	}

	// tag images
	if err = tagImages(images); err != nil {
		cicd.LogError(err)
		return err
	}
	log.Println("tagged images:", images)

	// push tagged images
	var result []string
	if result, err = ar.Push(wf, images); err != nil {
		cicd.LogError(err)
		return err
	}
	log.Println("pushed images:", result)
	return err
}

func makeTagList(repoURL string) (images []string) {

	// tag additional images based on build event type
	tagSep := strings.Index(baseImage, ":")
	commitImage := repoURL + baseImage[tagSep:]

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

func tagImages(images []string) (err error) {
	var stderr bytes.Buffer

	for _, image := range images {
		cmd := exec.Command("docker", "tag", baseImage, image)
		cmd.Stderr = &stderr
		log.Printf("docker tag from %v to %v", baseImage, image)

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
