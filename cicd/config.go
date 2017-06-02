package cicd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/viper"

	yaml "gopkg.in/yaml.v2"
)

type Workflow struct {
	Config
	App
	Provider
}

type Config struct {
	Enabled  bool
	Debug    bool
	Dryrun   bool
	Provider struct {
		CI struct {
			ID      string
			Enabled bool
		}
		CD struct {
			ID      string
			Enabled bool
		}
		Registry struct {
			ID      string
			Enabled bool
		}
		Platform struct {
			ID      string
			Enabled bool
		}
	}
}

type App struct {
	Name string
	Repo string
}

type reqs map[string][]interface{}

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
		MiniKube
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
	Push([]string) ([]string, error)
	Authenticate() error
	GetRepoURL() string
}

type Deployer interface {
	Deploy(*Workflow) error
}

func New() *Workflow {
	wf := Workflow{}
	return &wf
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

//TODO: create getActive funcs for CIProvider and Platform
func (wf *Workflow) GetActiveRegistry() (activeRegistry interface{}, err error) {
	switch wf.Config.Provider.Registry.ID {
	case "gcr":
		activeRegistry = &wf.Provider.Registry.GCR
	case "docker":
		activeRegistry = &wf.Provider.Registry.Docker
	default:
		err = fmt.Errorf("unknown workflow registry: <%v>", wf.Config.Provider.Registry.ID)
		log.Println(err)
	}
	return activeRegistry, err
}

func (wf *Workflow) GetActiveCDProvider() (activeCD interface{}, err error) {
	switch wf.Config.Provider.CD.ID {
	case "helm":
		activeCD = &wf.Provider.CD.Helm
	default:
		err = fmt.Errorf("unknown workflow CD provider: <%v>", wf.Config.Provider.CD.ID)
		log.Println(err)
	}
	return activeCD, err
}

func (wf *Workflow) UseContext() (err error) {
	var stderr bytes.Buffer
	var cmdOut []byte
	var ctx string

	switch wf.Config.Provider.Platform.ID {
	case "gke":
		ctx = wf.Provider.Platform.GKE.Context
	case "minikube":
		ctx = wf.Provider.Platform.MiniKube.Context
	default:
		LogError(fmt.Errorf("unknown platform provider: <%v>", wf.Config.Provider.Platform.ID))
	}

	cmd := exec.Command("kubectl", "config", "use-context", ctx)
	cmd.Stderr = &stderr

	log.Println(viper.GetString("cmdMode"), strings.Join(cmd.Args, " "))

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

func logCmdOutput(cmdOut []byte) {
	for _, o := range strings.Split(strings.TrimSpace(string(cmdOut)), "\n") {
		log.Println(o)
	}
}
