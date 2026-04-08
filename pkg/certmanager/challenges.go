package certmanager

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

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
// ListChallenges lists matching resources.
func ListChallenges(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []acmev1.Challenge {
	challenges, err := ListChallengesE(t, options, namespace)
	require.NoError(t, err, "Failed to list Challenges in namespace %s", namespace)
	return challenges
}

// ListChallengesE lists matching resources.
func ListChallengesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]acmev1.Challenge, error) {
	client, err := NewClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	challengeList, err := client.AcmeV1().Challenges(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return challengeList.Items, nil
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
// WaitForChallengeValid waits for the resource condition to be satisfied.
func WaitForChallengeValid(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForChallengeValidE(t, options, name, namespace, timeout)
	require.NoError(t, err, "ACME Challenge %s/%s not in Valid state", namespace, name)
}

// WaitForChallengeValidE waits for the resource condition to be satisfied.
func WaitForChallengeValidE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	client, err := NewClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		challenge, err := client.AcmeV1().Challenges(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		return challenge.Status.State == acmev1.Valid, nil
	})
}
