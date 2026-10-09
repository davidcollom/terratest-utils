package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

// testKubeconfigYAML is a minimal but valid kubeconfig that
// k8s.LoadAPIClientConfigE can parse. It defines a single cluster, context,
// and user. No real API call is made — the tests only assert that
// GetRestConfigE returns the expected *rest.Config (or an error).
const testKubeconfigYAML = `apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://test-server:6443
    insecure-skip-tls-verify: true
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
users:
- name: test-user
  user:
    token: fake-token-for-unit-test
`

// TestGetRestConfigE_ReturnsRestConfigIfPresent verifies that when
// options.RestConfig is already populated, GetRestConfigE returns it
// directly without consulting the kubeconfig file.
func TestGetRestConfigE_ReturnsRestConfigIfPresent(t *testing.T) {
	expectedConfig := &rest.Config{Host: "https://example.com"}
	options := &k8s.KubectlOptions{
		RestConfig: expectedConfig,
	}
	cfg, err := GetRestConfigE(t, options)
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, cfg)
}

// TestGetRestConfigE_LoadsFromKubeconfig verifies that when RestConfig is
// nil, GetRestConfigE resolves the kubeconfig path via options.GetConfigPath
// and then loads the REST config from disk via k8s.LoadAPIClientConfigE,
// returning a *rest.Config with the expected cluster server URL.
func TestGetRestConfigE_LoadsFromKubeconfig(t *testing.T) {
	kubeconfigPath := filepath.Join(t.TempDir(), "kubeconfig")
	require.NoError(t, os.WriteFile(kubeconfigPath, []byte(testKubeconfigYAML), 0o600))

	options := &k8s.KubectlOptions{
		ConfigPath:  kubeconfigPath,
		ContextName: "test-context",
	}
	cfg, err := GetRestConfigE(t, options)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, "https://test-server:6443", cfg.Host)
}

// TestGetRestConfigE_ReturnsErrorOnInvalidKubeconfig verifies that an error
// (and a nil *rest.Config) is returned when the configured kubeconfig path
// does not exist on disk.
func TestGetRestConfigE_ReturnsErrorOnInvalidKubeconfig(t *testing.T) {
	options := &k8s.KubectlOptions{
		ConfigPath:  filepath.Join(t.TempDir(), "does-not-exist"),
		ContextName: "test-context",
	}
	cfg, err := GetRestConfigE(t, options)
	assert.Error(t, err)
	assert.Nil(t, cfg)
}
