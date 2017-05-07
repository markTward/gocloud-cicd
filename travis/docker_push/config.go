package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type App struct {
	Name string
}
type Github struct {
	Repo string
}

type Registry struct {
	GCRRegistry
	DockerRegistry
}

type Registrator interface {
	IsRegistryValid() error
	Push([]string) ([]string, error)
	Authenticate() error
	GetRepoURL() string
}

type GCRRegistry struct {
	Name        string
	Description string
	Host        string
	Project     string
	Repo        string
	Url         string
	KeyFile     string
}

func (r *GCRRegistry) GetRepoURL() (repoURL string) {
	repo := []string{r.Host, r.Project, r.Repo}
	repoURL = strings.Join(repo, "/")
	return repoURL
}

func (r *GCRRegistry) Authenticate() (err error) {
	var stderr bytes.Buffer

	if _, err = os.Stat(r.KeyFile); os.IsNotExist(err) {
		err = fmt.Errorf("gcloud authentication: %v", err)
		return err
	}

	cmd := exec.Command("gcloud", "auth", "activate-service-account", "--key-file", r.KeyFile)
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

func (gcr *GCRRegistry) Push(images []string) (pushed []string, err error) {
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

func (r *GCRRegistry) IsRegistryValid() (err error) {
	if r.Url == "" {
		err = fmt.Errorf("url missing from %v configuration", r.Description)
	}
	return err
}

// TODO: obsolete now that gcloud auth output captured.  but would json parse of other fields be useful?
func (r *GCRRegistry) getClientID() (email string, err error) {

	// parse google credentials for identity
	type clientSecret struct {
		ClientEmail string `json:"client_email"`
	}

	// read in service account credentials file
	var jsonInput []byte
	if jsonInput, err = ioutil.ReadFile(r.KeyFile); err != nil {
		return "", fmt.Errorf("get service account id: %v", err)
	}

	// parse json for client email
	cs := clientSecret{}
	err = json.Unmarshal([]byte(jsonInput), &cs)

	return cs.ClientEmail, err
}

type DockerRegistry struct {
	Name        string
	Description string
	Host        string
	Account     string
	Repo        string
	Url         string
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

type Workflow struct {
	Enabled bool

	Github struct {
		Repo   string
		Branch string
	}

	CIProvider struct {
		Name string
		Plan string
	}

	Platform struct {
		Name    string
		Project string
		Cluster string
	}

	Registry string

	CDProvider struct {
		Name      string
		Release   string
		Namespace string
		ChartDir  string
	}
}

func logCmdOutput(cmdOut []byte) {
	for _, o := range strings.Split(strings.TrimSpace(string(cmdOut)), "\n") {
		log.Println(o)
	}
}
