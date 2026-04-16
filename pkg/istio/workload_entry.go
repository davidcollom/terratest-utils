package istio

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	istionetworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListWorkloadEntries retrieves all Istio WorkloadEntry resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list WorkloadEntries from.
//
// Returns:
//   - A slice of pointers to WorkloadEntry objects found in the namespace.
// ListWorkloadEntries lists matching resources.
func ListWorkloadEntries(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*istionetworkingv1alpha3.WorkloadEntry {
	workloadEntries, err := ListWorkloadEntriesE(t, options, namespace)
	require.NoError(t, err, "Failed to list Workload Entries in namespace %s", namespace)
	return workloadEntries
}

// ListWorkloadEntriesE lists matching resources.
func ListWorkloadEntriesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*istionetworkingv1alpha3.WorkloadEntry, error) {
	istioClient := NewClient(t, options)

	ctx := context.Background()
	workloadEntries, err := istioClient.NetworkingV1alpha3().WorkloadEntries(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	return workloadEntries.Items, nil
}

// WaitForWorkloadEntryReady waits until the specified WorkloadEntry in the given namespace is Ready or the timeout is reached.
// It polls the WorkloadEntry status every 2 seconds and checks for the Ready condition.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the WorkloadEntry to check.
//   - namespace: The namespace of the WorkloadEntry.
//   - timeout: The maximum duration to wait for the resource to become Ready.
// WaitForWorkloadEntryReady waits for the resource condition to be satisfied.
func WaitForWorkloadEntryReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForWorkloadEntryReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "WorkloadEntry %s/%s did not become Ready", namespace, name)
}

// WaitForWorkloadEntryReadyE waits for the resource condition to be satisfied.
func WaitForWorkloadEntryReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var workloadEntry *istionetworkingv1alpha3.WorkloadEntry
		workloadEntry, err := istioClient.NetworkingV1alpha3().WorkloadEntries(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		if workloadEntry.Status.Conditions != nil {
			return istioConditionReady(t, &workloadEntry.Status), nil
		}
		return false, nil
	})
}
