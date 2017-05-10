package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

type DockerRegistry struct {
	Name        string
	Description string
	Host        string
	Account     string
	Repo        string
	Url         string
}

func (r *DockerRegistry) Copy(data Registry) {
	t := reflect.ValueOf(r).Elem().Type()
	for i := 0; i < t.NumField(); i++ {
		n := t.Field(i).Name
		reflect.ValueOf(r).Elem().FieldByName(n).SetString(data[n])
	}
}

func (r *DockerRegistry) Authenticate() (err error) {
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
	cmd.Stderr = &stderr
	log.Println(strings.Join(cmd.Args[:4], " "), " -p ********")

	if cmdOut, err = cmd.Output(); err != nil {
		err = fmt.Errorf("%v", stderr.String())
		return err
	}

	logCmdOutput(cmdOut)

	return err
}

func (r *DockerRegistry) IsRegistryValid() (err error) {
	if r.Url == "" {
		err = fmt.Errorf("url missing from %v configuration", r.Description)
	}
	return err
}

func (docker *DockerRegistry) Push(images []string) (pushed []string, err error) {
	var stderr bytes.Buffer
	var cmdOut []byte

	for _, image := range images {

		cmd := exec.Command("docker", "push", image)
		cmd.Stderr = &stderr

		log.Println(strings.Join(cmd.Args, " "))

		if cmdOut, err = cmd.Output(); err != nil {
			err = fmt.Errorf("%v: %v", image, stderr.String())
			break
		}

		logCmdOutput(cmdOut)

		pushed = append(pushed, image)
	}
	return pushed, err
}

func (r *DockerRegistry) GetRepoURL() (repoURL string) {
	repo := []string{r.Host, r.Account, r.Repo}
	repoURL = strings.Join(repo, "/")
	return repoURL
}
