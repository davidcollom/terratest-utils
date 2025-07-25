package utils

import (
	"testing"

	"k8s.io/client-go/rest"

	"github.com/gruntwork-io/terratest/modules/k8s"
)

// GetRestConfigE retrieves a Kubernetes REST client configuration based on the provided KubectlOptions.
// If the RestConfig field in options is already set, it returns that configuration directly.
// Otherwise, it attempts to load the configuration from the kubeconfig file path specified in options,
// using the provided context name. Returns the REST config or an error if loading fails.
//
// Parameters:
//   - t: The testing context, used for logging and helper tracking.
//   - options: The KubectlOptions containing configuration and context information.
//
// Returns:
//   - *rest.Config: The Kubernetes REST client configuration.
//   - error: An error if the configuration could not be loaded.
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
