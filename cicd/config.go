package cicd

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/urfave/cli"
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
	GCR
	Docker
}

type Workflow struct {
	Enabled bool

	// provider ids
	Provider struct {
		CI       string
		CD       string
		Registry string
	}

	Github struct {
		Repo   string
		Branch string
	}

	CIProvider struct {
		Travis
	}

	Platform struct {
		GKE
	}

	CDProvider struct {
		Helm
	}
}

type Registrator interface {
	IsRegistryValid() error
	Push([]string, bool) ([]string, error)
	Authenticate() error
	GetRepoURL() string
}

type Deployer interface {
	Deploy(*cli.Context, *Config) error
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

func (cfg *Config) GetActiveRegistry() (activeRegistry interface{}, err error) {
	switch cfg.Workflow.Provider.Registry {
	case "gcr":
		activeRegistry = &cfg.Registry.GCR
	case "docker":
		activeRegistry = &cfg.Registry.Docker
	default:
		err = fmt.Errorf("unknown workflow registry: <%v>", cfg.Workflow.Provider.Registry)
		log.Println(err)
	}
	return activeRegistry, err
}

func (cfg *Config) GetActiveCDProvider() (activeCD interface{}, err error) {
	switch cfg.Workflow.Provider.CD {
	case "helm":
		activeCD = &cfg.CDProvider.Helm
	default:
		err = fmt.Errorf("unknown workflow CD provider: <%v>", cfg.Workflow.Provider.CD)
		log.Println(err)
	}
	return activeCD, err
}

func logCmdOutput(cmdOut []byte) {
	for _, o := range strings.Split(strings.TrimSpace(string(cmdOut)), "\n") {
		log.Println(o)
	}
}
