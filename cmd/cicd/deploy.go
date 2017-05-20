package main

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/markTward/gocloud-cicd"
	"github.com/urfave/cli"
)

var buildTag, containerRepo, serviceName, namespace, chartPath string

var deployCmd = cli.Command{
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

func deploy(ctx *cli.Context) error {

	// initialize configuration object
	wf := cicd.New()
	if err := cicd.Load(configFile, wf); err != nil {
		cicd.LogError(err)
		return err
	}
	cicd.LogDebug(ctx, fmt.Sprintf("%v", spew.Sdump(wf)))

	// initialize active Registry indicated by config and assert as Registrator
	var activeRegistry interface{}
	var err error
	if activeRegistry, err = wf.GetActiveRegistry(); err != nil {
		cicd.LogError(err)
		return err
	}
	ar := activeRegistry.(cicd.Registrator)

	// validate args and apply defaults
	if err = validateDeployArgs(ctx, wf, ar); err != nil {
		cicd.LogError(err)
		return err
	}
	log.Println("deploy command args:", cicd.GetAllFlags(ctx))

	// get active CD provider indicated by config and assert as Deployer
	var activeCDProvider interface{}
	if activeCDProvider, err = wf.GetActiveCDProvider(); err != nil {
		cicd.LogError(err)
		return err
	}
	ad := activeCDProvider.(cicd.Deployer)

	// deploy using active CD provider
	if err = ad.Deploy(ctx, wf); err != nil {
		cicd.LogError(err)
	}

	return err
}

func validateDeployArgs(ctx *cli.Context, wf *cicd.Workflow, ar cicd.Registrator) (err error) {

	// handle globals from cli and/or workflow config
	if cicd.IsDebug(ctx, wf) {
		debug = true
	}

	if cicd.IsDryRun(ctx, wf) {
		dryrun = true
	}

	//
	if buildTag == "" {
		err = fmt.Errorf("%v", "build tag a required value")
		return err
	}

	if branch == "" {
		err = fmt.Errorf("%v", "branch a required value")
	}

	if namespace == "" {
		if ns := wf.Provider.CD.Helm.Namespace; ns == "" {
			err = fmt.Errorf("%v", "namespace required when not defined in cicd.yaml")
			return err
		} else {
			namespace = ns
		}
	}

	if chartPath == "" {
		if cp := wf.Provider.CD.Helm.Chartpath; cp == "" {
			err = fmt.Errorf("%v", "chart path required when not defined in cicd.yaml")
			return err
		} else {
			chartPath = cp
		}
	}

	if isNotExist(chartPath) {
		cicd.LogDebug(ctx, fmt.Sprintf("is not exist chartpath: %v", chartPath))
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
		if svc := wf.App.Name; svc == "" {
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