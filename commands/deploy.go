package commands

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"github.com/markTward/gocloud-cicd/config"
	"gopkg.in/urfave/cli.v1"
)

var buildTag, containerRepo, serviceName, namespace, chartPath string

var DeployCmd = cli.Command{
	Name:  "deploy",
	Usage: "deploy services to providers (helm ==> k8s)",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "branch, b",
			Usage:       "build branch (required)",
			Destination: &branch,
		},
		cli.StringFlag{
			Name:        "chart",
			Usage:       "`PATH` to helm charts",
			Destination: &chartPath,
		},
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "load configuration file from `FILE`",
			Value:       "./cicd.yaml",
			Destination: &configFile,
		},
		cli.BoolFlag{
			Name:        "dryrun",
			Usage:       "log output but do not execute",
			Destination: &dryrun,
		},
		cli.StringFlag{
			Name:        "repo, r",
			Usage:       "repository source for images",
			Destination: &containerRepo,
		},
		cli.StringFlag{
			Name:        "namespace, n",
			Usage:       "target namespace",
			Destination: &namespace,
		},
		cli.StringFlag{
			Name:        "service, s",
			Usage:       "service name",
			Destination: &serviceName,
		},
		cli.StringFlag{
			Name:        "tag, t",
			Usage:       "existing image tag used as basis for further tags (required)",
			Destination: &buildTag,
		},
	},
	Action: deploy,
}

func deploy(c *cli.Context) error {

	LogDebug(c,
		fmt.Sprintf("flag values: --config %v, --tag %v, -branch %v, --repo %v,--service %v, --namespace %v, --chartpath %v --debug %v, --dryrun %v",
			configFile, buildTag, branch, containerRepo, serviceName, namespace, chartPath, c.GlobalBool("debug"), dryrun))

	// initialize configuration object
	cfg := config.New()
	if err := config.Load(configFile, &cfg); err != nil {
		LogDebug(c, "config error?")
		LogError(err)
		return err
	}

	LogDebug(c, fmt.Sprintf("Config: %#v", cfg))

	// initialize active Registry indicated by config and assert as Registrator
	var activeRegistry interface{}
	var err error
	if activeRegistry, err = cfg.GetActiveRegistry(); err != nil {
		LogError(err)
		return err

	}
	ar := activeRegistry.(config.Registrator)

	if err = validateDeployArgs(c, &cfg, ar); err != nil {
		LogError(err)
		return err
	}

	// initialize active CD provider indicated by config and assert as Deployer
	var activeCDProvider interface{}
	if activeCDProvider, err = cfg.GetActiveCDProvider(); err != nil {
		LogError(err)
		return err
	}
	ad := activeCDProvider.(config.Deployer)

	// TODO: pass args and move logic to helm.Deploy method
	release := serviceName + "-" + branch

	// helm required flags
	args := []string{"--install", release, "--namespace", namespace}

	// config file boolean flags
	for _, flag := range cfg.Workflow.CDProvider.Helm.Options.Flags {
		args = append(args, flag)
	}

	// cli flag conversion
	if c.GlobalBool("debug") {
		args = append(args, "--debug")
	}

	// convert cicd --dryrun arg to helm dialect
	if dryrun {
		args = append(args, "--dry-run")
	}

	// write runtime helm --values <file> using when available in config  otherwise create/remove a TempFile
	outFile := cfg.Workflow.CDProvider.Helm.Options.Values.Output
	var valuesFile *os.File
	switch {
	case outFile == "":
		valuesFile, err = ioutil.TempFile("", "runtime_values.yaml.")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(valuesFile.Name())
	default:
		valuesFile, err = os.Create(outFile)
		if err != nil {
			return err
		}
		defer valuesFile.Close()
	}

	LogDebug(c, fmt.Sprintf("helm runtime values filename: %v", valuesFile.Name()))

	// render values file from template
	err = renderHelmValuesFile(c, &cfg, valuesFile, containerRepo, buildTag)
	if err != nil {
		return fmt.Errorf("renderHelmValuesFile(): %v", err)
	}

	// join flags and positional args
	args = append(args, "--values", valuesFile.Name())
	args = append(args, chartPath)

	// deploy using active CD provider
	if err = ad.Deploy(&cfg, args); err != nil {
		LogError(err)
	}

	return err
}

func renderHelmValuesFile(c *cli.Context, cfg *config.Config, valuesFile *os.File, repo string, tag string) error {
	type Values struct {
		Repo, Tag, ServiceType string
	}

	// Prepare some data to insert into the template.
	var values = Values{Repo: repo, Tag: tag}

	// initialize the template
	var t *template.Template
	var err error
	if t, err = template.ParseFiles(cfg.Workflow.CDProvider.Helm.Options.Values.Template); err != nil {
		return err
	}

	// render the template
	err = t.Execute(valuesFile, values)

	// verify rendered file contents
	yaml, err := ioutil.ReadFile(valuesFile.Name())
	if err != nil {
		log.Println("error read yaml:", err)
		return err
	}

	log.Println(string(yaml))

	return err
}

func validateDeployArgs(c *cli.Context, cfg *config.Config, ar config.Registrator) (err error) {

	if buildTag == "" {
		err = fmt.Errorf("%v", "build tag a required value")
		return err
	}

	if branch == "" {
		err = fmt.Errorf("%v", "branch a required value")
	}

	if namespace == "" {
		if ns := cfg.Workflow.CDProvider.Helm.Namespace; ns == "" {
			err = fmt.Errorf("%v", "namespace required when not defined in cicd.yaml")
			return err
		} else {
			namespace = ns
		}
	}

	if chartPath == "" {
		if cp := cfg.Workflow.CDProvider.Helm.Chartpath; cp == "" {
			err = fmt.Errorf("%v", "chart path required when not defined in cicd.yaml")
			return err
		} else {
			chartPath = cp
		}
	}

	if isNotExist(chartPath) {
		LogDebug(c, fmt.Sprintf("is not exist chartpath: %v", chartPath))
		err = fmt.Errorf("chart path invalid: %v", chartPath)
		return err
	}

	if containerRepo == "" {
		if cr := ar.GetRepoURL(); cr == "" {
			err = fmt.Errorf("%v\n", "repoitory url required when not defined in cicd.yaml")
			return err
		} else {
			containerRepo = cr
		}
	}

	if serviceName == "" {
		if svc := cfg.App.Name; svc == "" {
			err = fmt.Errorf("%v", "service name required when not defined in cicd.yaml")
			return err
		} else {
			serviceName = svc
		}
	}

	return err
}

func isNotExist(f string) bool {
	_, err := os.Stat(f)
	return os.IsNotExist(err)
}
