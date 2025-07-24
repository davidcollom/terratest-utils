package certmanager

import (
	"context"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	acmev1 "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	cmclientset "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
)

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
func WaitForOrderValid(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := cmclientset.NewForConfig(options.RestConfig)
	require.NoError(t, err, "Failed to create cert-manager clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		order, err := client.AcmeV1().Orders(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		return order.Status.State == acmev1.Valid, nil
	})

	if err != nil {
		t.Fatalf("ACME Order %s/%s not in Valid state: %v", namespace, name, err)
	}
}
