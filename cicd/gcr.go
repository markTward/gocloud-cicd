package cicd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
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
		err = fmt.Errorf("gcloud auth key: %v", err)
		return err
	}

	cmd := exec.Command("gcloud", "auth", "activate-service-account", "--key-file", r.Keyfile)
	cmd.Stderr = &stderr

	log.Println(viper.GetString("cmdMode"), strings.Join(cmd.Args, " "))
	if !IsDryRun() {
		if err = cmd.Run(); err != nil {
			logCmdOutput(stderr.Bytes())
			err = fmt.Errorf("%v", stderr.String())
			return err
		}
		// BUG: gcloud returning successful result over stderr (why?)
		logCmdOutput(stderr.Bytes())
	}

	return err

}

func (gcr *GCR) Push(images []string) (pushed []string, err error) {
	var stderr bytes.Buffer
	var cmdOut []byte

	for _, image := range images {

		cmd := exec.Command("gcloud", "docker", "--", "push", image)
		cmd.Stderr = &stderr

		log.Println(viper.GetString("cmdMode"), strings.Join(cmd.Args, " "))
		if !IsDryRun() {
			if cmdOut, err = cmd.Output(); err != nil {
				logCmdOutput(stderr.Bytes())
				err = fmt.Errorf("%v: %v", image, stderr.String())
				break
			}
			pushed = append(pushed, image)
			logCmdOutput(cmdOut)
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
