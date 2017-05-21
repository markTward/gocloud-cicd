package cicd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type GCR struct {
	Name        string
	Description string
	Host        string
	Project     string
	Repo        string
	Url         string
	Keyfile     string
}

func (r *GCR) GetRepoURL() (repoURL string) {
	return r.Url
}

func (r *GCR) Authenticate(wf *Workflow) (err error) {
	var stderr bytes.Buffer

	if _, err = os.Stat(r.Keyfile); os.IsNotExist(err) {
		err = fmt.Errorf("gcloud auth key: %v", err)
		return err
	}

	cmd := exec.Command("gcloud", "auth", "activate-service-account", "--key-file", r.Keyfile)
	cmd.Stderr = &stderr

	if !IsDryRun() {
		log.Println("execute:", strings.Join(cmd.Args, " "))

		if err = cmd.Run(); err != nil {
			logCmdOutput(stderr.Bytes())
			err = fmt.Errorf("%v", stderr.String())
			return err
		}

		// BUG: gcloud returning successful result over stderr (why?)
		logCmdOutput(stderr.Bytes())
	} else {
		log.Println("dryrun:", strings.Join(cmd.Args, " "))
	}

	return err

}

func (gcr *GCR) Push(wf *Workflow, images []string) (pushed []string, err error) {
	var stderr bytes.Buffer
	var cmdOut []byte

	for _, image := range images {

		cmd := exec.Command("gcloud", "docker", "--", "push", image)
		cmd.Stderr = &stderr

		if !IsDryRun() {
			log.Println("execute: ", strings.Join(cmd.Args, " "))

			if cmdOut, err = cmd.Output(); err != nil {
				logCmdOutput(stderr.Bytes())
				err = fmt.Errorf("%v: %v", image, stderr.String())
				break
			}
			pushed = append(pushed, image)
			logCmdOutput(cmdOut)

		} else {
			log.Println("dryrun: ", strings.Join(cmd.Args, " "))
		}

	}
	return pushed, err
}

func (r *GCR) IsRegistryValid() (err error) {
	// TODO: check existence of other required field and/or remove unnecessary (host, account/project, repo)
	if r.Url == "" {
		err = fmt.Errorf("registry url missing from %v configuration", r.Description)
	}
	return err
}
