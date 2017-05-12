package config

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

func (r *GCR) Authenticate() (err error) {
	var stderr bytes.Buffer

	if _, err = os.Stat(r.Keyfile); os.IsNotExist(err) {
		err = fmt.Errorf("gcloud authentication: %v", err)
		return err
	}

	cmd := exec.Command("gcloud", "auth", "activate-service-account", "--key-file", r.Keyfile)
	cmd.Stderr = &stderr

	log.Println(strings.Join(cmd.Args, " "))

	if err = cmd.Run(); err != nil {
		logCmdOutput(stderr.Bytes())
		err = fmt.Errorf("%v", stderr.String())
		return err
	}

	// BUG: gcloud returning successful result over stderr (why?)
	logCmdOutput(stderr.Bytes())

	return err

}

func (gcr *GCR) Push(images []string) (pushed []string, err error) {
	var stderr bytes.Buffer
	var cmdOut []byte
	// IDEA: could use single command to push all repo images: gcloud docker -- push gcr.io/k8sdemo-159622/gocloud
	// but assumes that process ALWAYS wants ALL tags for repo to be pushed.  good for isolated build env, but ...
	for _, image := range images {

		cmd := exec.Command("gcloud", "docker", "--", "push", image)
		cmd.Stderr = &stderr

		log.Println(strings.Join(cmd.Args, " "))

		if cmdOut, err = cmd.Output(); err != nil {
			logCmdOutput(stderr.Bytes())
			err = fmt.Errorf("%v: %v", image, stderr.String())
			break
		}

		logCmdOutput(cmdOut)

		pushed = append(pushed, image)

	}
	return pushed, err
}

func (r *GCR) IsRegistryValid() (err error) {
	if r.Url == "" {
		err = fmt.Errorf("url missing from %v configuration", r.Description)
	}
	return err
}
