package certmanager

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	acmev1 "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
)

// ListOrders retrieves a list of ACME Order resources from the specified namespace using the cert-manager client.
// It requires a testing context, kubectl options, and the target namespace as input.
// The function fails the test if the client cannot be created or if listing the Orders fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options used to configure the client.
//   - namespace: The Kubernetes namespace to search for Orders.
//
// Returns:
//   - A slice of acmev1.Order objects found in the specified namespace.
// ListOrders lists matching resources.
func ListOrders(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []acmev1.Order {
	orders, err := ListOrdersE(t, options, namespace)
	require.NoError(t, err, "Failed to list Orders in namespace %s", namespace)
	return orders
}

// ListOrdersE lists matching resources.
func ListOrdersE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]acmev1.Order, error) {
	client, err := NewClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	orderList, err := client.AcmeV1().Orders(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return orderList.Items, nil
}

// WaitForOrderValid waits until the specified ACME Order resource in the given namespace reaches the "Valid" state or the timeout is exceeded.
// It polls the Order status every 2 seconds using the cert-manager clientset.
// If the Order does not reach the "Valid" state within the timeout, the test fails with a fatal error.
//
// Parameters:
//
//	t        - The testing context.
//	options  - Kubectl options containing the REST config for Kubernetes API access.
//	name     - The name of the ACME Order resource.
//	namespace- The namespace where the ACME Order resource resides.
//	timeout  - The maximum duration to wait for the Order to become valid.
//
// Fails the test if the Order does not reach the "Valid" state within the specified timeout.
// WaitForOrderValid waits for the resource condition to be satisfied.
func WaitForOrderValid(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForOrderValidE(t, options, name, namespace, timeout)
	require.NoError(t, err, "ACME Order %s/%s not in Valid state", namespace, name)
}

// WaitForOrderValidE waits for the resource condition to be satisfied.
func WaitForOrderValidE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	client, err := NewClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		order, err := client.AcmeV1().Orders(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		return order.Status.State == acmev1.Valid, nil
	})
}
