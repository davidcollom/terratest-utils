package flux

import (
	"context"
	"testing"
	"time"

	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/util/wait"
)

// ListOCIRepositories retrieves a list of OCIRepository resources from the specified namespace
// using the provided Kubernetes options. It returns a slice of sourcev1.OCIRepository objects.
// The function requires a testing.T instance for error handling and test context propagation.
// It fails the test if the Flux client cannot be created or if listing the OCIRepositories fails.
//
// Parameters:
//   - t:        The testing.T instance used for test context and assertions.
//   - options:  The KubectlOptions specifying the Kubernetes context and configuration.
//   - namespace: The namespace from which to list OCIRepository resources.
//
// Returns:
//   - []sourcev1.OCIRepository: A slice containing the retrieved OCIRepository resources.
func ListOCIRepositories(t *testing.T, options *k8s.KubectlOptions, namespace string, opts ...client.ListOption) []sourcev1.OCIRepository {
	t.Helper()

	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	// Append the namespace to the list options
	opts = append(opts, client.InNamespace(namespace))

	ctx := t.Context()
	var repos sourcev1.OCIRepositoryList
	err = fluxclient.List(ctx, &repos, opts...)
	require.NoError(t, err, "Failed to list OCIRepositories in namespace %s", namespace)

	return repos.Items
}

// WaitForOCIRepositoryReady waits until the specified Flux OCIRepository resource becomes Ready within the given timeout.
// It polls the resource status at regular intervals and fails the test if the resource does not become Ready in time.
//
// Parameters:
//
//	t        - The testing context.
//	options  - Kubectl options containing the REST config for the Kubernetes cluster.
//	name     - The name of the OCIRepository resource.
//	namespace- The namespace where the OCIRepository resource is located.
//	timeout  - The maximum duration to wait for the resource to become Ready.
//
// Fails the test if the OCIRepository does not reach the Ready condition within the timeout.
func WaitForOCIRepositoryReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {

		var ocirepo sourcev1.OCIRepository
		err = fluxclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &ocirepo)
		if err != nil {
			return false, nil
		}
		return hasReadyCondition(ocirepo.Status.Conditions), nil
	})

	if err != nil {
		t.Fatalf("OCIRepository %s/%s did not become Ready: %v", namespace, name, err)
	}
}
