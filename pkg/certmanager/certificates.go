package certmanager

import (
	"context"
	"testing"
	"time"

	certv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListCertificates retrieves all cert-manager Certificate resources in the specified namespace.
// It uses the provided testing context and kubectl options to create a cert-manager client,
// then lists and returns the Certificate objects found in the namespace. The function will
// fail the test if there is an error creating the client or listing the certificates.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list Certificate resources.
//
// Returns:
//   - A slice of certv1.Certificate objects found in the specified namespace.
func ListCertificates(t *testing.T, options *k8s.KubectlOptions, namespace string) []certv1.Certificate {
	t.Helper()

	client, err := NewCertManagerClient(t, options)
	require.NoError(t, err, "Failed to create cert-manager clientset")

	ctx := t.Context()
	certList, err := client.CertmanagerV1().Certificates(namespace).List(ctx, v1.ListOptions{})
	require.NoError(t, err, "Failed to list Certificates in namespace %s", namespace)

	return certList.Items
}

// WaitForCertificateReady waits until the specified cert-manager Certificate resource is in the Ready state.
// It polls the Certificate status at regular intervals until the Ready condition is true or the timeout is reached.
// If the Certificate does not become Ready within the timeout, the test fails.
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the Kubernetes REST config.
//   - name: The name of the Certificate resource.
//   - namespace: The namespace of the Certificate resource.
//   - timeout: The maximum duration to wait for the Certificate to become Ready.
func WaitForCertificateReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := NewCertManagerClient(t, options)
	require.NoError(t, err, "Failed to create cert-manager clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		cert, err := client.CertmanagerV1().Certificates(namespace).Get(ctx, name, v1.GetOptions{})
		if err != nil {
			return false, nil // retry
		}

		for _, cond := range cert.Status.Conditions {
			if cond.Type == certv1.CertificateConditionReady && cond.Status == cmmetav1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("Certificate %s/%s was not Ready in time: %v", namespace, name, err)
	}
}

// ValidateCertificateSecret verifies that the Kubernetes Secret referenced by the given
// cert-manager Certificate contains both the "tls.crt" and "tls.key" data fields.
// It fails the test if either field is missing.
// The function sets the namespace in the provided KubectlOptions to match the Certificate's namespace
// before retrieving and validating the Secret.
//
// Parameters:
//
//	t       - The testing context.
//	options - Kubectl options for accessing the Kubernetes cluster.
//	cert    - The cert-manager Certificate resource whose Secret should be validated.
func ValidateCertificateSecret(t *testing.T, options *k8s.KubectlOptions, cert *certv1.Certificate) {
	// We need to ensure we're looking in the right namespace
	options.Namespace = cert.Namespace
	secret := k8s.GetSecret(t, options, cert.Spec.SecretName)

	if _, ok := secret.Data["tls.crt"]; !ok {
		t.Fatalf("Secret %s missing tls.crt", cert.Spec.SecretName)
	}
	if _, ok := secret.Data["tls.key"]; !ok {
		t.Fatalf("Secret %s missing tls.key", cert.Spec.SecretName)
	}
}
