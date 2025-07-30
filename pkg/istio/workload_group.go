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

// ListWorkloadGroups retrieves all Istio WorkloadGroup resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list WorkloadGroups from.
//
// Returns:
//   - A slice of pointers to WorkloadGroup objects found in the namespace.
func ListWorkloadGroups(t *testing.T, options *k8s.KubectlOptions, namespace string) []*istionetworkingv1alpha3.WorkloadGroup {
	t.Helper()

	istioClient := NewClient(t, options)

	ctx := t.Context()
	workloadGroups, err := istioClient.NetworkingV1alpha3().WorkloadGroups(namespace).List(ctx, v1meta.ListOptions{})
	require.NoError(t, err, "Failed to list Workload Groups in namespace %s", namespace)

	return workloadGroups.Items
}

// WaitForWorkloadGroupReady waits until the specified WorkloadGroup in the given namespace is Ready or the timeout is reached.
// It polls the WorkloadGroup status every 2 seconds and checks for the Ready condition.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the WorkloadGroup to check.
//   - namespace: The namespace of the WorkloadGroup.
//   - timeout: The maximum duration to wait for the resource to become Ready.
func WaitForWorkloadGroupReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := t.Context()
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var workloadGroup *istionetworkingv1alpha3.WorkloadGroup
		workloadGroup, err := istioClient.NetworkingV1alpha3().WorkloadGroups(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		if workloadGroup.Status.Conditions != nil {
			return istioConditionReady(t, &workloadGroup.Status), nil
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("WorkloadGroup %s/%s did not become Ready: %v", namespace, name, err)
	}
}
