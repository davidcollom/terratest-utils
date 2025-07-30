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

// ListGitRepositories retrieves all Flux GitRepository resources within the specified Kubernetes namespace.
// It uses the provided testing context and Kubectl options to create a Flux client, then lists the GitRepositories.
// If any error occurs during client creation or listing, the test will fail.
//
// Parameters:
//   - t: The testing context.
//   - options: The Kubectl options for connecting to the Kubernetes cluster.
//   - namespace: The namespace to search for GitRepository resources.
//
// Returns:
//   - A slice of sourcev1.GitRepository objects found in the specified namespace.
func ListGitRepositories(t *testing.T, options *k8s.KubectlOptions, namespace string, opts ...client.ListOption) []sourcev1.GitRepository {
	t.Helper()

	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	// Append the namespace to the list options
	opts = append(opts, client.InNamespace(namespace))

	ctx := t.Context()
	var repos sourcev1.GitRepositoryList
	err = fluxclient.List(ctx, &repos, opts...)
	require.NoError(t, err, "Failed to list GitRepositories in namespace %s", namespace)

	return repos.Items
}

// WaitForGitRepositoryReady waits until the specified Flux GitRepository resource becomes Ready within the given timeout.
// It polls the resource status every 2 seconds and fails the test if the resource does not become Ready in time.
//
// Parameters:
//
//	t        - The testing context.
//	options  - The kubectl options containing the REST config for the Kubernetes cluster.
//	name     - The name of the GitRepository resource.
//	namespace- The namespace where the GitRepository resource is located.
//	timeout  - The maximum duration to wait for the resource to become Ready.
//
// Fails the test if the GitRepository does not reach the Ready condition within the timeout.
func WaitForGitRepositoryReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()
	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {

		var repo sourcev1.GitRepository
		err = fluxclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &repo)
		if err != nil {
			return false, nil
		}
		return hasReadyCondition(repo.Status.Conditions), nil
	})

	if err != nil {
		t.Fatalf("GitRepository %s/%s did not become Ready: %v", namespace, name, err)
	}
}
