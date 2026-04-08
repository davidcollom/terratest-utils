package linkerd

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

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
// ListServiceProfiles lists matching resources.
func ListServiceProfiles(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*linkerdv1alpha2.ServiceProfile {
	serviceProfiles, err := ListServiceProfilesE(t, options, namespace)
	require.NoError(t, err, "Failed to list ServiceProfiles in namespace %s", namespace)
	return serviceProfiles
}

// ListServiceProfilesE lists matching resources.
func ListServiceProfilesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*linkerdv1alpha2.ServiceProfile, error) {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	serviceProfiles, err := linkerdClient.LinkerdV1alpha2().ServiceProfiles(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Convert slice of values to slice of pointers
	var result []*linkerdv1alpha2.ServiceProfile
	for i := range serviceProfiles.Items {
		result = append(result, &serviceProfiles.Items[i])
	}

	return result, nil
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
// GetServiceProfile gets a resource by name.
func GetServiceProfile(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string) *linkerdv1alpha2.ServiceProfile {
	serviceProfile, err := GetServiceProfileE(t, options, name, namespace)
	require.NoError(t, err, "Failed to get ServiceProfile %s in namespace %s", name, namespace)
	return serviceProfile
}

// GetServiceProfileE gets a resource by name.
func GetServiceProfileE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string) (*linkerdv1alpha2.ServiceProfile, error) {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	serviceProfile, err := linkerdClient.LinkerdV1alpha2().ServiceProfiles(namespace).Get(ctx, name, v1meta.GetOptions{})
	if err != nil {
		return nil, err
	}

	return serviceProfile, nil
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
// WaitForServiceProfileExists waits for the resource condition to be satisfied.
func WaitForServiceProfileExists(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForServiceProfileExistsE(t, options, name, namespace, timeout)
	require.NoError(t, err, "ServiceProfile %s/%s did not exist within timeout", namespace, name)
}

// WaitForServiceProfileExistsE waits for the resource condition to be satisfied.
func WaitForServiceProfileExistsE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := linkerdClient.LinkerdV1alpha2().ServiceProfiles(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
}
