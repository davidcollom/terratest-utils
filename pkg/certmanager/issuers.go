package certmanager

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
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
// ListIssuers lists matching resources.
func ListIssuers(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []cmv1.Issuer {
	issuers, err := ListIssuersE(t, options, namespace)
	require.NoError(t, err, "Failed to list Issuers in namespace %s", namespace)
	return issuers
}

// ListIssuersE lists matching resources.
func ListIssuersE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]cmv1.Issuer, error) {
	client, err := NewClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	issuerList, err := client.CertmanagerV1().Issuers(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return issuerList.Items, nil
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
// WaitForIssuerReady waits for the resource condition to be satisfied.
func WaitForIssuerReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForIssuerReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "Issuer %s/%s not Ready", namespace, name)
}

// WaitForIssuerReadyE waits for the resource condition to be satisfied.
func WaitForIssuerReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	client, err := NewClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
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
// ListClusterIssuers lists matching resources.
func ListClusterIssuers(t testing.TestingT, options *k8s.KubectlOptions) []cmv1.ClusterIssuer {
	clusterIssuers, err := ListClusterIssuersE(t, options)
	require.NoError(t, err, "Failed to list ClusterIssuers")
	return clusterIssuers
}

// ListClusterIssuersE lists matching resources.
func ListClusterIssuersE(t testing.TestingT, options *k8s.KubectlOptions) ([]cmv1.ClusterIssuer, error) {
	client, err := NewClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	issuerList, err := client.CertmanagerV1().ClusterIssuers().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return issuerList.Items, nil
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
// WaitForClusterIssuerReady waits for the resource condition to be satisfied.
func WaitForClusterIssuerReady(t testing.TestingT, options *k8s.KubectlOptions, name string, timeout time.Duration) {
	err := WaitForClusterIssuerReadyE(t, options, name, timeout)
	require.NoError(t, err, "ClusterIssuer %s not Ready", name)
}

// WaitForClusterIssuerReadyE waits for the resource condition to be satisfied.
func WaitForClusterIssuerReadyE(t testing.TestingT, options *k8s.KubectlOptions, name string, timeout time.Duration) error {
	client, err := NewClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
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
}
