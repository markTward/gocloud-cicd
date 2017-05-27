package cmd

import (
	"fmt"
	"os"

	"github.com/markTward/gocloud-cicd/cicd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildTag, containerRepo, serviceName, namespace, chartPath string

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:           "deploy",
	Short:         "deploy containerzied applications",
	Long:          "deploy containerzied applications",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          deploy,
}

func init() {

	deployCmd.Flags().StringVarP(&branch, "branch", "b", "", "branch name for tagging")
	deployCmd.Flags().StringVarP(&chartPath, "chart", "", "", "path to helm charts")
	deployCmd.Flags().StringVarP(&containerRepo, "repo", "r", "", "container repository url")
	deployCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "k8s namespace for service")
	deployCmd.Flags().StringVarP(&serviceName, "service", "s", "", "app/service name")
	deployCmd.Flags().StringVarP(&buildTag, "tag", "t", "", "existing image tag used as basis for further tags (required)")

	viper.BindPFlag("branch", deployCmd.Flags().Lookup("branch"))
	viper.BindPFlag("chart", deployCmd.Flags().Lookup("chart"))
	viper.BindPFlag("repo", deployCmd.Flags().Lookup("repo"))
	viper.BindPFlag("namespace", deployCmd.Flags().Lookup("namespace"))
	viper.BindPFlag("service", deployCmd.Flags().Lookup("service"))
	viper.BindPFlag("tag", deployCmd.Flags().Lookup("tag"))

	RootCmd.AddCommand(deployCmd)

}

func deploy(ctx *cobra.Command, args []string) error {

	// initialize active Registry indicated by config and assert as Registrator
	var activeRegistry interface{}
	var err error
	if activeRegistry, err = wf.GetActiveRegistry(); err != nil {
		return err
	}
	ar := activeRegistry.(cicd.Registrator)

	// validate args and apply defaults
	if err = validateDeployArgs(ctx, wf, ar); err != nil {
		return err
	}

	// get active CD provider indicated by config and assert as Deployer
	var activeCDProvider interface{}
	if activeCDProvider, err = wf.GetActiveCDProvider(); err != nil {
		return err
	}
	ad := activeCDProvider.(cicd.Deployer)

	// deploy using active CD provider
	err = ad.Deploy(wf)
	return err
}

func validateDeployArgs(ctx *cobra.Command, wf *cicd.Workflow, ar cicd.Registrator) (err error) {

	if buildTag == "" {
		return fmt.Errorf("%v", "build tag a required value")
	}

	if branch == "" {
		return fmt.Errorf("%v", "branch a required value")
	}

	if namespace == "" {
		if ns := wf.Provider.CD.Helm.Namespace; ns == "" {
			return fmt.Errorf("%v", "namespace required when not defined in cicd.yaml")
		} else {
			namespace = ns
		}
	}

	if chartPath == "" {
		if cp := wf.Provider.CD.Helm.Chartpath; cp == "" {
			return fmt.Errorf("%v", "chart path required when not defined in cicd.yaml")
		} else {
			chartPath = cp
		}
	}

	// test existence of chart path
	_, err = os.Stat(chartPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("chart path invalid: %v", chartPath)
	}

	if containerRepo == "" {
		if cr := ar.GetRepoURL(); cr == "" {
			return fmt.Errorf("%v\n", "repoitory url required when not defined in cicd.yaml")
		} else {
			containerRepo = cr
		}
	}

	if serviceName == "" {
		if svc := wf.App.Name; svc == "" {
			return fmt.Errorf("%v", "service name required when not defined in cicd.yaml")
		} else {
			serviceName = svc
		}
	}

	return err
}
