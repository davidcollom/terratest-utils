package linkerd

import (
	"context"
	"testing"
	"time"

	"github.com/davidcollom/terratest-utils/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
)

var (
	// TrafficSplitGVR represents the GroupVersionResource for SMI TrafficSplit
	TrafficSplitGVR = schema.GroupVersionResource{
		Group:    "split.smi-spec.io",
		Version:  "v1alpha1",
		Resource: "trafficsplits",
	}
)

// ListTrafficSplits retrieves all SMI TrafficSplit resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list TrafficSplits from.
//
// Returns:
//   - A slice of pointers to unstructured objects representing TrafficSplit resources found in the namespace.
func ListTrafficSplits(t *testing.T, options *k8s.KubectlOptions, namespace string) []*unstructured.Unstructured {
	t.Helper()

	dynamicClient := NewDynamicClient(t, options)

	ctx := t.Context()
	trafficSplits, err := dynamicClient.Resource(TrafficSplitGVR).Namespace(namespace).List(ctx, v1meta.ListOptions{})
	require.NoError(t, err, "Failed to list TrafficSplits in namespace %s", namespace)

	var result []*unstructured.Unstructured
	for i := range trafficSplits.Items {
		result = append(result, &trafficSplits.Items[i])
	}

	return result
}

// GetTrafficSplit retrieves a specific SMI TrafficSplit resource by name in the specified namespace.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the TrafficSplit to retrieve.
//   - namespace: The namespace of the TrafficSplit.
//
// Returns:
//   - An unstructured object representing the TrafficSplit resource.
func GetTrafficSplit(t *testing.T, options *k8s.KubectlOptions, name, namespace string) *unstructured.Unstructured {
	t.Helper()

	dynamicClient := NewDynamicClient(t, options)

	ctx := t.Context()
	trafficSplit, err := dynamicClient.Resource(TrafficSplitGVR).Namespace(namespace).Get(ctx, name, v1meta.GetOptions{})
	require.NoError(t, err, "Failed to get TrafficSplit %s in namespace %s", name, namespace)

	return trafficSplit
}

// WaitForTrafficSplitExists waits until the specified TrafficSplit exists in the given namespace or the timeout is reached.
// It polls the TrafficSplit every 2 seconds until it exists.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the TrafficSplit to check.
//   - namespace: The namespace of the TrafficSplit.
//   - timeout: The maximum duration to wait for the resource to exist.
func WaitForTrafficSplitExists(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	dynamicClient := NewDynamicClient(t, options)

	ctx := t.Context()
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := dynamicClient.Resource(TrafficSplitGVR).Namespace(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		t.Fatalf("TrafficSplit %s/%s did not exist within timeout: %v", namespace, name, err)
	}
}

// NewDynamicClient creates and returns a new dynamic Kubernetes client for use with custom resources.
// It initializes the Kubernetes REST configuration and the dynamic client,
// failing the test if any errors occur during setup.
//
// Parameters:
//   - t: The testing context used for logging and error handling.
//   - options: The kubectl options specifying the context and namespace.
//
// Returns:
//   - dynamic.Interface: A dynamic client for interacting with custom resources.
func NewDynamicClient(t *testing.T, options *k8s.KubectlOptions) dynamic.Interface {
	cfg, err := utils.GetRestConfigE(t, options)
	require.NoError(t, err)

	client, err := dynamic.NewForConfig(cfg)
	require.NoError(t, err, "Failed to create dynamic client")

	return client
}
