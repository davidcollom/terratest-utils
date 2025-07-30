package linkerd

import (
	"testing"

	"github.com/davidcollom/terratest-utils/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/k8s"
	linkerdclientset "github.com/linkerd/linkerd2/controller/gen/client/clientset/versioned"
	"github.com/stretchr/testify/require"
)

// NewClient creates and returns a new Linkerd Client for use in tests.
// It initializes the Kubernetes REST configuration and the Linkerd clientset,
// failing the test if any errors occur during setup.
//
// Parameters:
//   - t: The testing context used for logging and error handling.
//   - options: The kubectl options specifying the context and namespace.
//
// Returns:
//   - *linkerdclientset.Clientset: A pointer to the initialized Linkerd client.
func NewClient(t *testing.T, options *k8s.KubectlOptions) *linkerdclientset.Clientset {
	cfg, err := utils.GetRestConfigE(t, options)
	require.NoError(t, err)

	client, err := linkerdclientset.NewForConfig(cfg)
	require.NoError(t, err, "Failed to create Linkerd client")

	return client
}
