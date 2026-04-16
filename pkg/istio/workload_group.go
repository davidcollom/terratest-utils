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
// ListWorkloadGroups lists matching resources.
func ListWorkloadGroups(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*istionetworkingv1alpha3.WorkloadGroup {
	workloadGroups, err := ListWorkloadGroupsE(t, options, namespace)
	require.NoError(t, err, "Failed to list Workload Groups in namespace %s", namespace)
	return workloadGroups
}

// ListWorkloadGroupsE lists matching resources.
func ListWorkloadGroupsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*istionetworkingv1alpha3.WorkloadGroup, error) {
	istioClient := NewClient(t, options)

	ctx := context.Background()
	workloadGroups, err := istioClient.NetworkingV1alpha3().WorkloadGroups(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	return workloadGroups.Items, nil
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
// WaitForWorkloadGroupReady waits for the resource condition to be satisfied.
func WaitForWorkloadGroupReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForWorkloadGroupReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "WorkloadGroup %s/%s did not become Ready", namespace, name)
}

// WaitForWorkloadGroupReadyE waits for the resource condition to be satisfied.
func WaitForWorkloadGroupReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
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
}
