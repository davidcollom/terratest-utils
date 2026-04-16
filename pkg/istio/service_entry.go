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

// ListServiceEntries retrieves all Istio ServiceEntry resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list ServiceEntries from.
//
// Returns:
//   - A slice of pointers to ServiceEntry objects found in the namespace.
// ListServiceEntries lists matching resources.
func ListServiceEntries(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*istionetworkingv1alpha3.ServiceEntry {
	serviceEntries, err := ListServiceEntriesE(t, options, namespace)
	require.NoError(t, err, "Failed to list Service Entries in namespace %s", namespace)
	return serviceEntries
}

// ListServiceEntriesE lists matching resources.
func ListServiceEntriesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*istionetworkingv1alpha3.ServiceEntry, error) {
	istioClient := NewClient(t, options)

	ctx := context.Background()
	serviceEntries, err := istioClient.NetworkingV1alpha3().ServiceEntries(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	return serviceEntries.Items, nil
}

// WaitForServiceEntryReady waits until the specified ServiceEntry in the given namespace is Ready or the timeout is reached.
// It polls the ServiceEntry status every 2 seconds and checks for the Ready condition.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the ServiceEntry to check.
//   - namespace: The namespace of the ServiceEntry.
//   - timeout: The maximum duration to wait for the resource to become Ready.
// WaitForServiceEntryReady waits for the resource condition to be satisfied.
func WaitForServiceEntryReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForServiceEntryReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "ServiceEntry %s/%s did not become Ready", namespace, name)
}

// WaitForServiceEntryReadyE waits for the resource condition to be satisfied.
func WaitForServiceEntryReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var serviceEntry *istionetworkingv1alpha3.ServiceEntry
		serviceEntry, err := istioClient.NetworkingV1alpha3().ServiceEntries(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		if serviceEntry.Status.Conditions != nil {
			return serviceEntryConditionReady(t, &serviceEntry.Status), nil
		}
		return false, nil
	})
}
