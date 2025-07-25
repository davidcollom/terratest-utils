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

func ListCertificateRequests(t *testing.T, options *k8s.KubectlOptions, namespace string) []cmv1.CertificateRequest {
	t.Helper()

	client, err := NewCertManagerClient(t, options)
	require.NoError(t, err, "Failed to create cert-manager clientset")

	ctx := t.Context()
	crList, err := client.CertmanagerV1().CertificateRequests(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list CertificateRequests in namespace %s", namespace)

	return crList.Items
}

// WaitForCertificateRequestReady waits until the specified CertificateRequest resource in the given namespace
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
func WaitForCertificateRequestReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := NewCertManagerClient(t, options)
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

	if err != nil {
		t.Fatalf("CertificateRequest %s/%s not Ready: %v", namespace, name, err)
	}
}
