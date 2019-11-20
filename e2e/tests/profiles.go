package tests

import (
	"fmt"
	"path/filepath"

	"github.com/devspace-cloud/devspace/cmd"
	"github.com/devspace-cloud/devspace/cmd/flags"
	"github.com/devspace-cloud/devspace/cmd/use"
	"github.com/devspace-cloud/devspace/e2e/utils"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/configutil"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/generated"
	"github.com/devspace-cloud/devspace/pkg/devspace/kubectl"
	"github.com/devspace-cloud/devspace/pkg/devspace/services"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RunProfiles runs the test for the kustomize example
func RunProfiles(namespace string) error {
	log.Info("Run Profiles")

	// We reset the previous config
	configutil.ResetConfig()
	generated.ResetConfig()

	var deployConfig = &cmd.DeployCmd{
		GlobalFlags: &flags.GlobalFlags{
			Namespace: namespace,
			NoWarn:    true,
		},
		ForceBuild:  true,
		ForceDeploy: true,
		SkipPush:    true,
	}

	wd, err := filepath.Abs("../examples/profiles/")
	fmt.Println(wd)

	if err != nil {
		return err
	}
	utils.ChangeWorkingDir(wd)
	if err != nil {
		return err
	}

	// Create kubectl client
	var client kubectl.Client
	client, err = kubectl.NewClientFromContext(deployConfig.KubeContext, deployConfig.Namespace, deployConfig.SwitchContext)
	if err != nil {
		return errors.Errorf("Unable to create new kubectl client: %v", err)
	}

	// At last, we delete the current namespace
	defer utils.DeleteNamespaceAndWait(client, deployConfig.Namespace)

	err = runProfile(deployConfig, "dev-service2-only", client, namespace, 2, false)
	if err != nil {
		return err
	}

	err = runProfile(deployConfig, "", client, namespace, 2, true)
	if err != nil {
		return err
	}

	return nil
}

func runProfile(deployConfig *cmd.DeployCmd, profile string, client kubectl.Client, namespace string, numberPodsExpected int, reset bool) error {
	var profileConfig = &use.ProfileCmd{
		Reset: reset,
	}

	err := profileConfig.RunUseProfile(nil, []string{profile})
	if err != nil {
		return err
	}

	err = deployConfig.Run(nil, nil)
	if err != nil {
		return err
	}

	// Checking if pods are running correctly
	utils.AnalyzePods(client, namespace)

	pods, errp := client.KubeClient().CoreV1().Pods(namespace).List(v1.ListOptions{})
	if errp != nil {
		return err
	}
	fmt.Println(len(pods.Items))
	if len(pods.Items) != numberPodsExpected {
		return errors.Errorf("There should be %v pod(s) running", numberPodsExpected)
	}

	// Load generated config
	generatedConfig, err := generated.LoadConfig(deployConfig.Profile)
	if err != nil {
		return errors.Errorf("Error loading generated.yaml: %v", err)
	}

	// Add current kube context to context
	configOptions := deployConfig.ToConfigOptions()
	config, err := configutil.GetConfig(configOptions)
	if err != nil {
		return err
	}

	servicesClient := services.NewClient(config, generatedConfig, client, nil, log.GetInstance())

	// Port-forwarding
	err = utils.PortForwardAndPing(servicesClient)
	if err != nil {
		return err
	}

	return nil
}
