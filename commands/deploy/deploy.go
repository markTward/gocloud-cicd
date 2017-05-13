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
	},
	Action: deploy,
}

func deploy(c *cli.Context) error {

	cmd.LogDebug(c,
		fmt.Sprintf("flag values: --config %v, --tag %v, -branch %v, --repo %v,--service %v, --namespace %v, --chartpath %v --debug %v, -verbose %v\n",
			configFile, buildTag, branch, containerRepo, serviceName, namespace, chartPath, c.GlobalBool("debug"), c.GlobalBool("verbose")))

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

	if err := validateCLInput(c, &cfg, ar); err != nil {
		cmd.LogError(err)
		return err
	}

	// TODO: make release construction a func/rule that could vary by project?
	release := serviceName + "-" + branch

	cmd.LogDebug(c, fmt.Sprintf("service: %v", serviceName))
	cmd.LogDebug(c, fmt.Sprintf("release: %v", release))
	cmd.LogDebug(c, fmt.Sprintf("namespace: %v", namespace))
	cmd.LogDebug(c, fmt.Sprintf("chartpath: %v", chartPath))
	cmd.LogDebug(c, fmt.Sprintf("repo url: %v", ar.GetRepoURL()))

	// helm upgrade \
	// $DRYRUN_OPTION \
	// --debug \
	// --install $RELEASE_NAME \
	// --namespace=$NAMESPACE \
	//TODO: add --set service.gocloud... to cicd.yaml. how to render?
	// --set service.gocloudAPI.image.repository=$DOCKER_REPO \
	// --set service.gocloudAPI.image.tag=":$COMMIT_TAG" \
	// --set service.gocloudGrpc.image.repository=$DOCKER_REPO \
	// --set service.gocloudGrpc.image.tag=":$COMMIT_TAG" \
	// $CHARTPATH

	//TODO: derive activeCDProvider in similar way as registry

	// prepare arguments for helm upgrade
	args := []string{"--install", release, "--namespace", namespace}
	if c.GlobalBool("dryrun") {
		args = append(args, "--dryrun")
	}
	args = append(args, chartPath)

	// TODO: process / render workflow.cdprovider.helm.options.(set, ...)
	log.Println("helm upgrade args: ", args)

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
		err = fmt.Errorf("%v\n", "build tag a required value; use --tag option")
	}

	if branch == "" {
		err = fmt.Errorf("%v\n", "build branch a required value; use --branch option")
	}

	if namespace == "" {
		if ns := cfg.Workflow.CDProvider.Helm.Namespace; ns == "" {
			err = fmt.Errorf("%v\n", "namespace required when not defined in cicd.yaml")
		} else {
			namespace = ns
		}
	}

	if chartPath == "" {
		if cp := cfg.Workflow.CDProvider.Helm.ChartPath; cp == "" {
			err = fmt.Errorf("%v\n", "chart path required when not defined in cicd.yaml")
		} else {
			chartPath = cp
		}
	}

	if _, err = os.Stat(chartPath); os.IsNotExist(err) {
		cmd.LogDebug(c, fmt.Sprintf("chart path invalid: %v", err))
		return fmt.Errorf("chart path invalid: %v", err)
	}

	if serviceName == "" {
		if svc := cfg.App.Name; svc == "" {
			err = fmt.Errorf("%v\n", "service name required when not defined in cicd.yaml")
		} else {
			serviceName = svc
		}
	}

	if containerRepo == "" {
		if cr := ar.GetRepoURL(); cr == "" {
			err = fmt.Errorf("%v\n", "repoitory url required when not defined in cicd.yaml")
		} else {
			containerRepo = cr
		}
	}

	return err
}
