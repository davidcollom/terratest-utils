package certmanager

import (
	"context"
	"testing"
	"time"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	cmclientset "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListIssuers retrieves a list of cert-manager Issuer resources from the specified namespace.
// It uses the provided testing context and kubectl options to create a cert-manager client,
// then queries for Issuers in the given namespace. The function fails the test if any error occurs
// during client creation or resource listing.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options used to configure the client.
//   - namespace: The Kubernetes namespace from which to list Issuers.
//
// Returns:
//   - A slice of cmv1.Issuer objects found in the specified namespace.
func ListIssuers(t *testing.T, options *k8s.KubectlOptions, namespace string) []cmv1.Issuer {
	t.Helper()

	client, err := NewCertManagerClient(t, options)
	require.NoError(t, err, "Failed to create cert-manager clientset")

	ctx := t.Context()
	issuerList, err := client.CertmanagerV1().Issuers(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list Issuers in namespace %s", namespace)

	return issuerList.Items
}

// WaitForIssuerReady waits until the specified cert-manager Issuer resource is in the Ready condition within the given timeout.
// It polls the Issuer status every 2 seconds and fails the test if the Issuer does not become Ready within the timeout period.
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the Kubernetes REST config.
//   - name: The name of the Issuer resource.
//   - namespace: The namespace of the Issuer resource.
//   - timeout: The maximum duration to wait for the Issuer to become Ready.
//
// Fails the test if the Issuer is not Ready within the timeout or if there is an error creating the cert-manager clientset.
func WaitForIssuerReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := NewCertManagerClient(t, options)
	require.NoError(t, err, "Failed to create cert-manager clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		issuer, err := client.CertmanagerV1().Issuers(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		for _, cond := range issuer.Status.Conditions {
			if cond.Type == cmv1.IssuerConditionReady && cond.Status == cmmetav1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("Issuer %s/%s not Ready: %v", namespace, name, err)
	}
}

// ListClusterIssuers retrieves a list of cert-manager ClusterIssuer resources from the Kubernetes cluster
// using the provided KubectlOptions. It requires a testing.T instance for error handling and test context.
// The function returns a slice of ClusterIssuer objects. If the client creation or resource listing fails,
// the test will fail with an appropriate error message.
//
// Parameters:
//   - t: A pointer to testing.T, used for test context and error reporting.
//   - options: A pointer to k8s.KubectlOptions, containing configuration for accessing the Kubernetes cluster.
//
// Returns:
//   - A slice of cmv1.ClusterIssuer representing the ClusterIssuers found in the cluster.
func ListClusterIssuers(t *testing.T, options *k8s.KubectlOptions) []cmv1.ClusterIssuer {
	t.Helper()

	client, err := cmclientset.NewForConfig(options.RestConfig)
	require.NoError(t, err, "Failed to create cert-manager clientset")

	ctx := t.Context()
	issuerList, err := client.CertmanagerV1().ClusterIssuers().List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list ClusterIssuers")

	return issuerList.Items
}

// WaitForClusterIssuerReady waits until the specified cert-manager ClusterIssuer resource is in the Ready state.
// It polls the ClusterIssuer status at regular intervals until the Ready condition is true or the timeout is reached.
// If the ClusterIssuer does not become Ready within the timeout, the test fails.
//
// Parameters:
//
//	t       - The testing context.
//	options - The kubectl options containing Kubernetes REST config.
//	name    - The name of the ClusterIssuer to check.
//	timeout - The maximum duration to wait for the ClusterIssuer to become Ready.
//
// This function requires a cert-manager clientset and uses the provided REST config to interact with the Kubernetes API.
func WaitForClusterIssuerReady(t *testing.T, options *k8s.KubectlOptions, name string, timeout time.Duration) {
	t.Helper()

	client, err := cmclientset.NewForConfig(options.RestConfig)
	require.NoError(t, err, "Failed to create cert-manager clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		issuer, err := client.CertmanagerV1().ClusterIssuers().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		for _, cond := range issuer.Status.Conditions {
			if cond.Type == cmv1.IssuerConditionReady && cond.Status == cmmetav1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("ClusterIssuer %s not Ready: %v", name, err)
	}
}
