package deploy

import (
	"fmt"
	"log"
	"os"

	cmd "github.com/markTward/gocloud-cicd/commands"
	"github.com/markTward/gocloud-cicd/config"
	"gopkg.in/urfave/cli.v1"
)

var configFile, buildTag, containerRepo, branch, serviceName, namespace, chartPath string
var dryrun bool

var DeployCmd = cli.Command{
	Name:  "deploy",
	Usage: "deploy services to providers (helm ==> k8s)",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "repo, r",
			Usage:       "repository source for images",
			Destination: &containerRepo,
		},
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "configuration file containing project workflow values",
			Value:       "./cicd.yaml",
			Destination: &configFile,
		},
		cli.StringFlag{
			Name:        "service, s",
			Usage:       "service name",
			Destination: &serviceName,
		},
		cli.StringFlag{
			Name:        "namespace, n",
			Usage:       "target namespace",
			Destination: &namespace,
		},
		cli.StringFlag{
			Name:        "tag, t",
			Usage:       "existing image tag used as basis for further tags (required)",
			Destination: &buildTag,
		},
		cli.StringFlag{
			Name:        "branch, b",
			Usage:       "build branch (required)",
			Destination: &branch,
		},
		cli.StringFlag{
			Name:        "chart",
			Usage:       "path to helm charts",
			Destination: &chartPath,
		},
		cli.BoolFlag{
			Name:        "dryrun",
			Usage:       "log output but do not execute",
			Destination: &dryrun,
		},
	},
	Action: deploy,
}

func deploy(c *cli.Context) error {

	cmd.LogDebug(c,
		fmt.Sprintf("flag values: --config %v, --tag %v, -branch %v, --repo %v,--service %v, --namespace %v, --chartpath %v --debug %v, --dryrun %v",
			configFile, buildTag, branch, containerRepo, serviceName, namespace, chartPath, c.GlobalBool("debug"), dryrun))

	// initialize configuration object
	cfg := config.New()
	if err := config.Load(configFile, &cfg); err != nil {
		cmd.LogDebug(c, "config error?")
		cmd.LogError(err)
		return err
	}

	cmd.LogDebug(c, fmt.Sprintf("Config: %#v", cfg))

	// initialize active registry indicated by config
	var activeRegistry interface{}
	var err error
	if activeRegistry, err = cfg.GetActiveRegistry(); err != nil {
		cmd.LogError(err)
		return err

	}
	ar := activeRegistry.(config.Registrator)

	if err = validateCLInput(c, &cfg, ar); err != nil {
		cmd.LogError(err)
		return err
	}

	// TODO: pass args and move logic to helm.Deploy method
	// TODO: make release construction a func/rule that could vary by project/plan?
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

	if dryrun {
		args = append(args, "--dry-run")
	}

	// TODO: values file used  from config is static for testing only
	// convert to template and render with dynamic repo and tag values
	for _, v := range cfg.CDProvider.Helm.Options.Values {
		cmd.LogDebug(c, fmt.Sprintf("add values file: --values %v", v))
		args = append(args, "--values", v)
	}

	// chart must be last positional argument
	args = append(args, chartPath)

	// TODO: process / render workflow.cdprovider.helm.options.(set, ...)
	cmd.LogDebug(c, fmt.Sprintf("helm upgrade args: %v", args))

	// for testing only
	helm := cfg.Workflow.CDProvider.Helm
	if err = helm.Deploy(&cfg, args); err != nil {
		cmd.LogError(err)
		return err
	}
	log.Println("deploy_k8s: helm deploy successful")

	return err
}

func validateCLInput(c *cli.Context, cfg *config.Config, ar config.Registrator) (err error) {

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
		cmd.LogDebug(c, fmt.Sprintf("is not exist chartpath: %v", chartPath))
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
