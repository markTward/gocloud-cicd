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

func push(ctx *cobra.Command, args []string) (err error) {

	// validate args
	if err := validatePushArgs(ctx, wf); err != nil {
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
	if err := ar.Authenticate(ctx, wf); err != nil {
		cicd.LogError(err)
		return err
	}

	// make list of images to tag
	var images []string
	log.Println("calling makeTagList branch: ", branch)
	if images = makeTagList(ctx, ar.GetRepoURL(), baseImage, event, branch, pr); len(images) == 0 {
		cicd.LogError(fmt.Errorf("no images to tag: %v", images))
		return err
	}
	log.Println("after maketaglist iamges:", images, len(images))

	// tag images
	if err := tagImages(baseImage, images); err != nil {
		cicd.LogError(err)
		return err
	}
	log.Println("tagged images:", images)

	// push tagged images
	var result []string
	log.Println("about to call gcr push with parent context", ctx.Parent().Name())
	if result, err = ar.Push(ctx.Parent(), wf, images); err != nil {
		cicd.LogError(err)
		return err
	}
	log.Println("pushed images:", result)
	return err
}

func makeTagList(ctx *cobra.Command, repoURL string, refImage string, event string, branch string, pr string) (images []string) {

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

func validatePushArgs(ctx *cobra.Command, wf *cicd.Workflow) (err error) {

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
