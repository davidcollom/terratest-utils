// Package certmanager provides Terratest-style helpers for testing cert-manager
// resources including Certificates, Issuers, ClusterIssuers, CertificateRequests,
// ACME Orders, and Challenges. These helpers use client-go and polling logic to
// wait for readiness conditions and validate associated Secrets.
package certmanager

import (
	"testing"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	cmclientset "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/davidcollom/terratest-utils/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"k8s.io/client-go/rest"
)

// NewCertManagerClient creates and returns a new cert-manager clientset.Interface using the provided testing context and kubectl options.
// If the RestConfig in options is nil, it attempts to generate a new rest.Config using the provided options.
// Returns the cert-manager clientset.Interface or an error if the configuration could not be created.
//
// Parameters:
//   - t: The testing context, used for error reporting and helper annotation.
//   - options: The kubectl options containing cluster access configuration. If options.RestConfig is nil,
//     the function will attempt to generate a rest.Config using utils.GetRestConfigE.
//
// Returns:
//   - cmclientset.Interface: The cert-manager clientset for interacting with cert-manager resources.
//   - error: An error if the configuration or clientset could not be created.
//
// Example usage:
//
//	client, err := NewCertManagerClient(t, options)
//	require.NoError(t, err)
//	issuers, err := client.CertmanagerV1().Issuers("default").List(context.TODO(), metav1.ListOptions{})
//	require.NoError(t, err)
func NewCertManagerClient(t *testing.T, options *k8s.KubectlOptions) (cmclientset.Interface, error) {
	t.Helper()
	var cfg *rest.Config
	var err error
	if options.RestConfig == nil {
		cfg, err = utils.GetRestConfigE(t, options)
		if err != nil {
			return nil, err
		}
	} else {
		cfg = options.RestConfig
	}

	return cmclientset.NewForConfig(cfg)
}

// HasCondition checks if a slice of CertificateRequestCondition contains a condition
// with the specified type and status.
//
// Parameters:
//
//	conds    - Slice of CertificateRequestCondition to search.
//	condType - The condition type to look for.
//	status   - The condition status to match.
//
// Returns:
//
//	true if a condition with the specified type and status exists, false otherwise.
func HasCondition(conds []cmv1.CertificateRequestCondition, condType cmv1.CertificateRequestConditionType, status cmmetav1.ConditionStatus) bool {
	for _, cond := range conds {
		if cond.Type == condType && cond.Status == status {
			return true
		}
	}
	return false
}
