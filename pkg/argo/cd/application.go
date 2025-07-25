package cd

import (
	"context"
	"testing"
	"time"

	argocdv1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListApplications retrieves a list of Argo CD Application resources from the specified namespace.
// It uses the provided testing context and kubectl options to create an Argo CD client,
// then lists all Application resources in the given namespace. The function fails the test
// if the client cannot be created or if the list operation fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options used to configure the client.
//   - namespace: The Kubernetes namespace from which to list Application resources.
//
// Returns:
//   - A slice of v1alpha1.Application representing the Applications found in the namespace.
func ListApplications(t *testing.T, options *k8s.KubectlOptions, namespace string) []argocdv1alpha1.Application {
	t.Helper()

	client, err := NewArgoCDClient(t, options)
	require.NoError(t, err, "Failed to create Argo clientset")

	ctx := t.Context()
	applicationList, err := client.ArgoprojV1alpha1().Applications(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list Applications in namespace %s", namespace)

	return applicationList.Items
}

// WaitForApplicationHealthyAndSynced waits until the specified Argo CD Application resource
// in the given namespace reaches both "Healthy" and "Synced" status within the provided timeout.
// It polls the Application status every 2 seconds using the Argo CD client and fails the test
// if the desired state is not achieved within the timeout period.
//
// Parameters:
//
//	t        - The testing context.
//	options  - Kubectl options containing the Kubernetes REST config.
//	name     - The name of the Argo CD Application.
//	namespace- The namespace where the Application resides.
//	timeout  - The maximum duration to wait for the Application to become Healthy and Synced.
//
// Fails the test if the Application does not reach the desired state within the timeout.
func WaitForApplicationHealthyAndSynced(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := NewArgoCDClient(t, options)
	require.NoError(t, err, "Unable to create Argo CD client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		app, err := client.ArgoprojV1alpha1().Applications(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		if IsApplicationHealthyAndSynced(app) {
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("Application %s/%s did not become Healthy & Synced: %v", namespace, name, err)
	}
}
