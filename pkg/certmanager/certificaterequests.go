package certmanager

import (
	"context"
	"testing"
	"time"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListCertificateRequests retrieves all CertificateRequest resources in the specified namespace
// using the provided kubectl options. It returns a slice of CertificateRequest objects.
// The function fails the test if there is an error creating the cert-manager client or listing the resources.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list CertificateRequests.
//
// Returns:
//   - A slice of cmv1.CertificateRequest representing the CertificateRequests found in the namespace.
func ListCertificateRequests(t *testing.T, options *k8s.KubectlOptions, namespace string) []cmv1.CertificateRequest {
	t.Helper()

	client, err := NewClient(t, options)
	require.NoError(t, err, "Failed to create cert-manager clientset")

	ctx := t.Context()
	crList, err := client.CertmanagerV1().CertificateRequests(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list CertificateRequests in namespace %s", namespace)

	return crList.Items
}

// WaitForCertificateRequestReadyE waits until the specified CertificateRequest resource in the given namespace
// reaches the Ready condition within the provided timeout duration. It polls the resource status every 2 seconds.
// If the CertificateRequest does not become Ready within the timeout, the test fails with a fatal error.
//
// Parameters:
//
//	t        - The testing context.
//	options  - The kubectl options containing the Kubernetes REST config.
//	name     - The name of the CertificateRequest resource.
//	namespace- The namespace where the CertificateRequest resides.
//	timeout  - The maximum duration to wait for the CertificateRequest to become Ready.
//
// This function requires cert-manager clientset and is intended for use in integration tests.
func WaitForCertificateRequestReadyE(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	t.Helper()

	client, err := NewClient(t, options)
	require.NoError(t, err, "Failed to create cert-manager clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		cr, err := client.CertmanagerV1().CertificateRequests(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		for _, cond := range cr.Status.Conditions {
			if cond.Type == cmv1.CertificateRequestConditionReady && cond.Status == cmmetav1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})

	return err
}

// WaitForCertificateRequestReady waits until the specified CertificateRequest resource in the given namespace
// reaches the "Ready" condition or the timeout is exceeded. It fails the test if the CertificateRequest does not
// become ready within the specified duration.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options to use for interacting with the Kubernetes cluster.
//   - name: The name of the CertificateRequest resource.
//   - namespace: The namespace where the CertificateRequest is located.
//   - timeout: The maximum duration to wait for the CertificateRequest to become ready.
func WaitForCertificateRequestReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForCertificateRequestReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err)
}
