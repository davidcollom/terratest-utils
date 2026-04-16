package cd

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

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
// ListApplications lists matching resources.
func ListApplications(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []argocdv1alpha1.Application {
	applications, err := ListApplicationsE(t, options, namespace)
	require.NoError(t, err, "Failed to list Applications in namespace %s", namespace)
	return applications
}

// ListApplicationsE retrieves a list of Argo CD Application resources from the specified namespace.
// It returns an error to the caller instead of failing the test directly.
// ListApplicationsE lists matching resources.
func ListApplicationsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]argocdv1alpha1.Application, error) {
	client, err := NewArgoCDClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	applicationList, err := client.ArgoprojV1alpha1().Applications(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return applicationList.Items, nil
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
// WaitForApplicationHealthyAndSynced waits for the resource condition to be satisfied.
func WaitForApplicationHealthyAndSynced(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForApplicationHealthyAndSyncedE(t, options, name, namespace, timeout)
	require.NoError(t, err, "Application %s/%s did not become Healthy & Synced", namespace, name)
}

// WaitForApplicationHealthyAndSyncedE waits until the specified Argo CD Application resource
// reaches both Healthy and Synced status within the provided timeout.
// WaitForApplicationHealthyAndSyncedE waits for the resource condition to be satisfied.
func WaitForApplicationHealthyAndSyncedE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	client, err := NewArgoCDClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		app, err := client.ArgoprojV1alpha1().Applications(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		if IsApplicationHealthyAndSynced(app) {
			return true, nil
		}
		return false, nil
	})
}
