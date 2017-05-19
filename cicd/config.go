package cicd

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

type Workflow struct {
	Config
	App
	Provider
}

type Config struct {
	Enabled  bool
	Provider struct {
		CI       string
		CD       string
		Registry string
	}
}

type App struct {
	Name string
	Repo string
}

type Provider struct {
	CICD struct {
		Repo   string
		Branch string
	}

	CI struct {
		Travis
	}

	Platform struct {
		GKE
	}

	CD struct {
		Helm
	}

	Registry struct {
		GCR
		Docker
	}
}

type Registrator interface {
	IsRegistryValid() error
	Push(*cli.Context, []string) ([]string, error)
	Authenticate() error
	GetRepoURL() string
}

type Deployer interface {
	Deploy(*cli.Context, *Workflow) error
}

func New() Workflow {
	return Workflow{}
}

func Load(cf string, wf *Workflow) error {
	// read in config yaml file
	yamlInput, err := ioutil.ReadFile(cf)
	if err != nil {
		return err
	}

	// parse yaml into Workflow object
	err = yaml.Unmarshal([]byte(yamlInput), &wf)
	return err

}

func (wf *Workflow) GetActiveRegistry() (activeRegistry interface{}, err error) {
	switch wf.Config.Provider.Registry {
	case "gcr":
		activeRegistry = &wf.Provider.Registry.GCR
	case "docker":
		activeRegistry = &wf.Provider.Registry.Docker
	default:
		err = fmt.Errorf("unknown workflow registry: <%v>", wf.Config.Provider.Registry)
		log.Println(err)
	}
	return activeRegistry, err
}

func (wf *Workflow) GetActiveCDProvider() (activeCD interface{}, err error) {
	switch wf.Config.Provider.CD {
	case "helm":
		activeCD = &wf.Provider.CD.Helm
	default:
		err = fmt.Errorf("unknown workflow CD provider: <%v>", wf.Config.Provider.CD)
		log.Println(err)
	}
	return activeCD, err
}

func logCmdOutput(cmdOut []byte) {
	for _, o := range strings.Split(strings.TrimSpace(string(cmdOut)), "\n") {
		log.Println(o)
	}
}
