package istio

import (
	"context"
	"testing"
	"time"

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
func ListWorkloadEntries(t *testing.T, options *k8s.KubectlOptions, namespace string) []*istionetworkingv1alpha3.WorkloadEntry {
	t.Helper()

	istioClient := NewClient(t, options)

	ctx := t.Context()
	workloadEntries, err := istioClient.NetworkingV1alpha3().WorkloadEntries(namespace).List(ctx, v1meta.ListOptions{})
	require.NoError(t, err, "Failed to list Workload Entries in namespace %s", namespace)

	return workloadEntries.Items
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
func WaitForWorkloadEntryReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := t.Context()
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
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

	if err != nil {
		t.Fatalf("WorkloadEntry %s/%s did not become Ready: %v", namespace, name, err)
	}
}
