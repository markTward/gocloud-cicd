package config

import (
	"io/ioutil"
	"log"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	App
	Github
	Registry
	Workflow
}

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
		GKE
	}

	Registry string

	CDProvider struct {
		Helm
	}
}

type Registrator interface {
	IsRegistryValid() error
	Push([]string) ([]string, error)
	Authenticate() error
	GetRepoURL() string
}

func New() Config {
	return Config{}
}

func Load(cf string, cfg *Config) error {
	// read in config yaml file
	yamlInput, err := ioutil.ReadFile(cf)
	if err != nil {
		return err
	}

	// parse yaml into Config object
	err = yaml.Unmarshal([]byte(yamlInput), &cfg)
	return err

}

func logCmdOutput(cmdOut []byte) {
	for _, o := range strings.Split(strings.TrimSpace(string(cmdOut)), "\n") {
		log.Println(o)
	}
}
