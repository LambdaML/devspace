package deploy

import (
	"bytes"

	"github.com/devspace-cloud/devspace/cmd"
	"github.com/devspace-cloud/devspace/cmd/flags"
	"github.com/devspace-cloud/devspace/e2e/utils"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//Test 4 - helm
//1. deploy & helm (see quickstart) (v1beta5 no tiller)
//2. purge (check if everything is deleted except namespace)

// RunHelm runs the test for the kubectl test
func RunHelm(f *customFactory, logger log.Logger) error {
	buff := &bytes.Buffer{}
	f.cacheLogger = log.NewStreamLogger(buff, logrus.InfoLevel)

	logger.Info("Run sub test 'helm' of test 'deploy'")
	logger.StartWait("Run test...")
	defer logger.StopWait()

	client, err := f.NewKubeClientFromContext("", f.namespace, false)
	if err != nil {
		return errors.Errorf("Unable to create new kubectl client: %v", err)
	}

	// The client is saved in the factory ONCE for each sub test
	f.client = client

	ts := testSuite{
		test{
			name: "1. deploy & helm (see quickstart) (v1beta5 no tiller)",
			deployConfig: &cmd.DeployCmd{
				GlobalFlags: &flags.GlobalFlags{
					Namespace: f.namespace,
					NoWarn:    true,
				},
			},
			postCheck: nil,
		},
	}

	err = beforeTest(f, logger, "tests/deploy/testdata/helm")
	defer afterTest(f)
	if err != nil {
		return errors.Errorf("sub test 'helm' of 'deploy' test failed: %s %v", buff.String(), err)
	}

	for _, t := range ts {
		err := runTest(f, &t)
		utils.PrintTestResult("helm", t.name, err, logger)
		if err != nil {
			return errors.Errorf("sub test 'helm' of 'deploy' test failed: %s %v", buff.String(), err)
		}
	}

	err = testPurge(f)
	utils.PrintTestResult("helm", "purge", err, logger)
	if err != nil {
		return err
	}

	return nil
}
