package certmanager

import (
	"context"
	"testing"
	"time"

	acmev1 "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListChallenges retrieves a list of ACME Challenge resources from the specified namespace
// using the cert-manager client. It requires a testing context, kubectl options, and the
// target namespace. The function will fail the test if the client cannot be created or if
// the challenges cannot be listed.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list ACME Challenges.
//
// Returns:
//   - A slice of acmev1.Challenge objects found in the specified namespace.
func ListChallenges(t *testing.T, options *k8s.KubectlOptions, namespace string) []acmev1.Challenge {
	t.Helper()

	client, err := NewClient(t, options)
	require.NoError(t, err, "Failed to create cert-manager clientset")

	ctx := t.Context()
	challengeList, err := client.AcmeV1().Challenges(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list Challenges in namespace %s", namespace)

	return challengeList.Items
}

// WaitForChallengeValid waits until the specified ACME Challenge resource in the given namespace
// reaches the "Valid" state or the timeout is exceeded. It polls the challenge status at regular
// intervals using the cert-manager clientset. If the challenge does not become valid within the
// timeout, the test fails with a fatal error.
//
// Parameters:
//
//	t         - The testing context.
//	options   - Kubectl options containing the REST config for Kubernetes API access.
//	name      - The name of the ACME Challenge resource.
//	namespace - The namespace where the challenge resource resides.
//	timeout   - The maximum duration to wait for the challenge to become valid.
//
// Fails the test if the challenge does not reach the "Valid" state within the timeout.
func WaitForChallengeValid(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := NewClient(t, options)
	require.NoError(t, err, "Failed to create cert-manager clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		challenge, err := client.AcmeV1().Challenges(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		return challenge.Status.State == acmev1.Valid, nil
	})

	if err != nil {
		t.Fatalf("ACME Challenge %s/%s not in Valid state: %v", namespace, name, err)
	}
}
