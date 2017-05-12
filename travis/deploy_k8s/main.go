package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/markTward/gocloud-cicd/travis/config"
)

var configFile, buildTag, containerRepo, branch, serviceName, namespace, chartPath *string
var dryrun *bool

func init() {
	const (
		defaultConfigFile  = "./cicd.yaml"
		configFileUsage    = "configuration file containing project workflow values"
		buildTagUsage      = "existing image tag used as basis for further tags (required)"
		containerRepoUsage = "repository for images (default config.workflow.registry)"
		branchUsage        = "build branch (required)"
		serviceNameUsage   = "service name"
		namespaceUsage     = "target namespace"
		chartPathUsage     = "path to chart for service"
		dryrunUsage        = "generate deployment artifacts without deploying"
	)

	configFile = flag.String("config", defaultConfigFile, configFileUsage)
	buildTag = flag.String("tag", "", buildTagUsage)
	containerRepo = flag.String("repo", "", containerRepoUsage)
	branch = flag.String("branch", "", branchUsage)
	serviceName = flag.String("service", "", serviceNameUsage)
	namespace = flag.String("namespace", "", namespaceUsage)
	chartPath = flag.String("chart", "", chartPathUsage)
	dryrun = flag.Bool("dryrun", false, dryrunUsage)

}

func main() {
	// parse and validate CLI
	flag.Parse()

	// initialize configuration object
	cfg := config.New()
	if err := config.Load(*configFile, &cfg); err != nil {
		exitScript(err, true)
	}

	log.Printf("flag values: --config %v, --tag %v, -branch %v, --repo %v,--service %v, --namespace %v, --chartpath %v, --dryrun %v\n",
		*configFile, *buildTag, *branch, *containerRepo, *serviceName, *namespace, *chartPath, *dryrun)

	// point to active registry (docker, gcr, ...)
	var activeRegistry interface{}
	switch cfg.Workflow.Registry {
	case "gcr":
		activeRegistry = &cfg.Registry.GCR
	case "docker":
		activeRegistry = &cfg.Registry.Docker
	default:
		log.Println("unknown registry")
		exitScript(fmt.Errorf("unknown workflow registry: <%v>", cfg.Workflow.Registry), true)
	}

	// assert activeRegistry as type Deployer to access methods
	ar := activeRegistry.(config.Deployer)

	if err := validateCLInput(&cfg, ar); err != nil {
		exitScript(err, true)
	}

	// TODO: make release construction a func/rule that could vary by project?
	release := *serviceName + "-" + *branch

	fmt.Println(cfg)
	fmt.Println("service:", *serviceName)
	fmt.Println("release:", release)
	fmt.Println("namespace:", *namespace)
	fmt.Println("chartpath:", *chartPath)
	fmt.Println("repo url:", ar.GetRepoURL())

}

func validateCLInput(cfg *config.Config, ar config.Deployer) (err error) {

	if *buildTag == "" {
		err = fmt.Errorf("%v\n", "build tag a required value; use --tag option")
	}

	if *branch == "" {
		err = fmt.Errorf("%v\n", "build branch a required value; use --branch option")
	}

	if *namespace == "" {
		if ns := cfg.Workflow.CDProvider.Helm.Namespace; ns == "" {
			err = fmt.Errorf("%v\n", "namespace required when not defined in cicd.yaml")
		} else {
			*namespace = ns
		}
	}

	if *chartPath == "" {
		if cp := cfg.Workflow.CDProvider.Helm.ChartPath; cp == "" {
			err = fmt.Errorf("%v\n", "chart path required when not defined in cicd.yaml")
		} else {
			*chartPath = cp
		}
	}

	if *serviceName == "" {
		if svc := cfg.App.Name; svc == "" {
			err = fmt.Errorf("%v\n", "service name required when not defined in cicd.yaml")
		} else {
			*serviceName = svc
		}
	}

	if *containerRepo == "" {
		log.Println("cicd workfloe repo:", ar.GetRepoURL())
		if cr := ar.GetRepoURL(); cr == "" {
			err = fmt.Errorf("%v\n", "repoitory url required when not defined in cicd.yaml")
		} else {
			*containerRepo = cr
		}
	}

	return err
}

func exitScript(err error, exit bool) {
	s := strings.TrimSpace(err.Error())
	log.Printf("error: %v", s)
	if exit {
		fmt.Fprintf(os.Stderr, "error: %v\n", s)
		os.Exit(1)
	}
}
