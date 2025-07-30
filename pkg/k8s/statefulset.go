package k8s

import (
	"context"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	terrak8s "github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
)

// GetStatefulSet retrieves the specified StatefulSet from the given Kubernetes namespace using the provided KubectlOptions and GetOptions.
// It fails the test immediately if the StatefulSet cannot be retrieved, using require.NoError.
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options to use when interacting with the cluster.
//   - name: The name of the StatefulSet to retrieve.
//   - namespace: The namespace where the StatefulSet is located.
//   - opts: Additional options for the get operation.
//
// Returns:
//   - A pointer to the retrieved appsv1.StatefulSet object.
func GetStatefulSet(t *testing.T, options *terrak8s.KubectlOptions, name, namespace string, opts metav1.GetOptions) *appsv1.StatefulSet {
	t.Helper()

	sts, err := GetStatefulSetE(t, options, name, namespace, opts)
	require.NoError(t, err)
	return sts
}

// GetStatefulSetE retrieves a Kubernetes StatefulSet resource by name and namespace using the provided KubectlOptions.
// It returns the StatefulSet object and an error if the retrieval fails. The function uses the testing context from t
// and allows passing custom metav1.GetOptions for the request.
//
// Parameters:
//   - t: The testing context.
//   - options: The KubectlOptions to configure the Kubernetes client.
//   - name: The name of the StatefulSet to retrieve.
//   - namespace: The namespace where the StatefulSet resides.
//   - opts: Additional options for the Get request.
//
// Returns:
//   - *appsv1.StatefulSet: The retrieved StatefulSet object.
//   - error: An error if the StatefulSet could not be retrieved.
func GetStatefulSetE(t *testing.T, options *terrak8s.KubectlOptions, name, namespace string, opts metav1.GetOptions) (*appsv1.StatefulSet, error) {
	t.Helper()

	client, err := terrak8s.GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()
	return client.AppsV1().StatefulSets(namespace).Get(ctx, name, opts)
}

// ListStatefulSets retrieves a list of Kubernetes StatefulSets in the cluster using the provided KubectlOptions and ListOptions.
// It fails the test immediately if an error occurs during retrieval.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - opts: The options to filter the list of StatefulSets.
//
// Returns:
//   - A slice of appsv1.StatefulSet objects representing the StatefulSets found.
func ListStatefulSets(t *testing.T, options *terrak8s.KubectlOptions, opts metav1.ListOptions) []appsv1.StatefulSet {
	t.Helper()

	statefulSets, err := ListStatefulSetsE(t, options, opts)
	require.NoError(t, err)
	return statefulSets
}

// ListStatefulSetsE retrieves a list of StatefulSets from the specified Kubernetes namespace using the provided
// KubectlOptions and ListOptions. It returns a slice of StatefulSet objects and an error if the retrieval fails.
// This function is intended for use in tests and will mark the test as a helper.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the Kubernetes context and namespace.
//   - opts: The list options to filter the StatefulSets.
//
// Returns:
//   - A slice of StatefulSet objects found in the specified namespace.
//   - An error if the StatefulSets could not be listed.
func ListStatefulSetsE(t *testing.T, options *terrak8s.KubectlOptions, opts metav1.ListOptions) ([]appsv1.StatefulSet, error) {
	t.Helper()

	client, err := terrak8s.GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()
	statefulSetList, err := client.AppsV1().StatefulSets(options.Namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return statefulSetList.Items, nil
}

// WaitForStatefulSetReady waits until the specified StatefulSet in the given namespace is ready or the timeout is reached.
// It polls the StatefulSet status every 2 seconds and checks if it is up-to-date using IsStatefulSetUptoDate.
// If the StatefulSet does not become ready within the timeout, the test fails with a fatal error.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for accessing the Kubernetes cluster.
//   - name: The name of the StatefulSet to check.
//   - namespace: The namespace where the StatefulSet is located.
//   - timeout: The maximum duration to wait for the StatefulSet to become ready.
func WaitForStatefulSetReady(t *testing.T, options *terrak8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := NewClient(t, options)
	require.NoError(t, err, "Failed to create stateful set clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		sts, err := client.AppsV1().StatefulSets(namespace).Get(ctx, name, v1.GetOptions{})
		if err != nil {
			return false, nil // retry
		}
		if IsStatefulSetUptoDate(sts) {
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("StatefulSet %s/%s was not Ready in time: %v", namespace, name, err)
	}
}

// IsStatefulSetUptoDate checks whether the given StatefulSet has all its replicas updated, available, and current.
// It returns true if the number of updated, available, and current replicas all match the desired replica count.
//
// Parameters:
//   - sts: A pointer to the appsv1.StatefulSet object to check.
//
// Returns:
//   - bool: True if the StatefulSet is up-to-date (all replicas are updated, available, and current), false otherwise.
//
// This function is useful for determining if a StatefulSet rollout has completed successfully.
func IsStatefulSetUptoDate(sts *appsv1.StatefulSet) bool {
	return sts.Status.UpdatedReplicas == sts.Status.Replicas &&
		sts.Status.AvailableReplicas == sts.Status.Replicas &&
		sts.Status.CurrentReplicas == sts.Status.Replicas
}
