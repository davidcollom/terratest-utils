package linkerd

import (
	"context"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	linkerdv1alpha2 "github.com/linkerd/linkerd2/controller/gen/apis/serviceprofile/v1alpha2"
	"github.com/stretchr/testify/require"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListServiceProfiles retrieves all Linkerd ServiceProfile resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list ServiceProfiles from.
//
// Returns:
//   - A slice of pointers to ServiceProfile objects found in the namespace.
func ListServiceProfiles(t *testing.T, options *k8s.KubectlOptions, namespace string) []*linkerdv1alpha2.ServiceProfile {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	serviceProfiles, err := linkerdClient.LinkerdV1alpha2().ServiceProfiles(namespace).List(ctx, v1meta.ListOptions{})
	require.NoError(t, err, "Failed to list ServiceProfiles in namespace %s", namespace)

	// Convert slice of values to slice of pointers
	var result []*linkerdv1alpha2.ServiceProfile
	for i := range serviceProfiles.Items {
		result = append(result, &serviceProfiles.Items[i])
	}

	return result
}

// GetServiceProfile retrieves a specific Linkerd ServiceProfile resource by name in the specified namespace.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the ServiceProfile to retrieve.
//   - namespace: The namespace of the ServiceProfile.
//
// Returns:
//   - A pointer to the ServiceProfile object.
func GetServiceProfile(t *testing.T, options *k8s.KubectlOptions, name, namespace string) *linkerdv1alpha2.ServiceProfile {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	serviceProfile, err := linkerdClient.LinkerdV1alpha2().ServiceProfiles(namespace).Get(ctx, name, v1meta.GetOptions{})
	require.NoError(t, err, "Failed to get ServiceProfile %s in namespace %s", name, namespace)

	return serviceProfile
}

// WaitForServiceProfileExists waits until the specified ServiceProfile exists in the given namespace or the timeout is reached.
// It polls the ServiceProfile every 2 seconds until it exists.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the ServiceProfile to check.
//   - namespace: The namespace of the ServiceProfile.
//   - timeout: The maximum duration to wait for the resource to exist.
func WaitForServiceProfileExists(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := linkerdClient.LinkerdV1alpha2().ServiceProfiles(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		t.Fatalf("ServiceProfile %s/%s did not exist within timeout: %v", namespace, name, err)
	}
}
