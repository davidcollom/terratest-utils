package utils

import (
	"testing"

	"k8s.io/client-go/rest"

	"github.com/gruntwork-io/terratest/modules/k8s"
)

func GetRestConfigE(t *testing.T, options *k8s.KubectlOptions) (*rest.Config, error) {
	t.Helper()

	if options.RestConfig != nil {
		return options.RestConfig, nil
	}

	cfgPath, err := options.GetConfigPath(t)
	if err != nil {
		return nil, err
	}
	return k8s.LoadApiClientConfigE(cfgPath, options.ContextName)
}
