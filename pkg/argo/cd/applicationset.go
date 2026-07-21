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

// ListApplicationSets retrieves all Argo CD ApplicationSet resources in the specified namespace.
// It uses the provided testing context and kubectl options to create an Argo CD client,
// then lists the ApplicationSets in the given namespace. The function fails the test if
// client creation or listing fails, and returns the list of ApplicationSet items.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list ApplicationSets.
//
// Returns:
//   - A slice of ApplicationSet resources found in the specified namespace.
//
// ListApplicationSets lists matching resources.
func ListApplicationSets(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []argocdv1alpha1.ApplicationSet {
	applicationSets, err := ListApplicationSetsE(t, options, namespace)
	require.NoError(t, err, "Failed to list ApplicationSets in namespace %s", namespace)
	return applicationSets
}

// ListApplicationSetsE retrieves all Argo CD ApplicationSet resources in the specified namespace.
// It returns an error to the caller instead of failing the test directly.
// ListApplicationSetsE lists matching resources.
func ListApplicationSetsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]argocdv1alpha1.ApplicationSet, error) {
	client, err := NewArgoCDClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	applicationSetList, err := client.ArgoprojV1alpha1().ApplicationSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return applicationSetList.Items, nil
}

// WaitForApplicationSetHealthyAndSynced waits until the specified Argo CD ApplicationSet in the given namespace
// is healthy and its resources are up to date, or until the provided timeout is reached.
// It polls the ApplicationSet status every 2 seconds, checking for the "ResourcesUpToDate" condition with a "True" status.
// If the ApplicationSet does not become healthy and synced within the timeout, the test fails.
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the Kubernetes REST config.
//   - name: The name of the ApplicationSet.
//   - namespace: The namespace where the ApplicationSet resides.
//   - timeout: The maximum duration to wait for the ApplicationSet to become healthy and synced.
//
// WaitForApplicationSetHealthyAndSynced waits for the resource condition to be satisfied.
func WaitForApplicationSetHealthyAndSynced(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForApplicationSetHealthyAndSyncedE(t, options, name, namespace, timeout)
	require.NoError(t, err, "ApplicationSet %s/%s did not become Healthy & Synced", namespace, name)
}

// WaitForApplicationSetHealthyAndSyncedE waits until the specified Argo CD ApplicationSet
// reports the resources-up-to-date condition within the provided timeout.
// WaitForApplicationSetHealthyAndSyncedE waits for the resource condition to be satisfied.
func WaitForApplicationSetHealthyAndSyncedE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	client, err := NewArgoCDClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		app, err := client.ArgoprojV1alpha1().ApplicationSets(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		for _, cond := range app.Status.Conditions {
			if cond.Type == argocdv1alpha1.ApplicationSetConditionResourcesUpToDate && cond.Status == argocdv1alpha1.ApplicationSetConditionStatusTrue {
				return true, nil
			}
		}
		return false, nil
	})
}
