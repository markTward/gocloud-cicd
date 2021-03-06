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

type Docker struct {
	Name        string
	Description string
	Host        string
	Account     string
	Repo        string
	Url         string
}

func (r *Docker) Authenticate() (err error) {
	var stderr bytes.Buffer
	var cmdOut []byte

	dockerUser := os.Getenv("DOCKER_USER")
	if dockerUser == "" {
		err = fmt.Errorf("DOCKER_USER environment variable not set")
		return err
	}

	dockerPass := os.Getenv("DOCKER_PASSWORD")
	if dockerPass == "" {
		err = fmt.Errorf("DOCKER_PASSWORD environment variable not set")
		return err
	}

	cmd := exec.Command("docker", "login", "-u", dockerUser, "-p", dockerPass)
	log.Println(viper.GetString("cmdMode"), strings.Join(cmd.Args[:4], " "), " -p ********")
	if !IsDryRun() {
		cmd.Stderr = &stderr
		if cmdOut, err = cmd.Output(); err != nil {
			err = fmt.Errorf("%v", stderr.String())
			return err
		}
		logCmdOutput(cmdOut)
	}

	return err
}

func (r *Docker) IsRegistryValid() (err error) {
	if r.Url == "" {
		err = fmt.Errorf("url missing from %v configuration", r.Description)
	}
	return err
}

func (docker *Docker) Push(images []string) (pushed []string, err error) {
	var stderr bytes.Buffer
	var cmdOut []byte

	for _, image := range images {

		cmd := exec.Command("docker", "push", image)
		cmd.Stderr = &stderr

		log.Println(viper.GetString("cmdMode"), strings.Join(cmd.Args, " "))

		if !IsDryRun() {
			if cmdOut, err = cmd.Output(); err != nil {
				err = fmt.Errorf("%v: %v", image, stderr.String())
				break
			}

			logCmdOutput(cmdOut)
			pushed = append(pushed, image)

		}
	}

	return pushed, err
}

func (r *Docker) GetRepoURL() (repoURL string) {
	return r.Url
}
