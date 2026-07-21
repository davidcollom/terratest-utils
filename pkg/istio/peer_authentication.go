package istio

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	istiosecurityv1 "istio.io/client-go/pkg/apis/security/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListPeerAuthentications retrieves all Istio PeerAuthentication resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list PeerAuthentications from.
//
// Returns:
//   - A slice of pointers to PeerAuthentication objects found in the namespace.
//
// ListPeerAuthentications lists matching resources.
func ListPeerAuthentications(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*istiosecurityv1.PeerAuthentication {
	peerAuthentications, err := ListPeerAuthenticationsE(t, options, namespace)
	require.NoError(t, err, "Failed to list Peer Authentications in namespace %s", namespace)
	return peerAuthentications
}

// ListPeerAuthenticationsE lists matching resources.
func ListPeerAuthenticationsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*istiosecurityv1.PeerAuthentication, error) {
	istioClient := NewClient(t, options)

	ctx := context.Background()
	peerAuthentications, err := istioClient.SecurityV1().PeerAuthentications(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	return peerAuthentications.Items, nil
}

// WaitForPeerAuthenticationReady waits until the specified PeerAuthentication in the given namespace is Ready or the timeout is reached.
// It polls the PeerAuthentication status every 2 seconds and checks for the Ready condition.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the PeerAuthentication to check.
//   - namespace: The namespace of the PeerAuthentication.
//   - timeout: The maximum duration to wait for the resource to become Ready.
//
// WaitForPeerAuthenticationReady waits for the resource condition to be satisfied.
func WaitForPeerAuthenticationReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForPeerAuthenticationReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "PeerAuthentication %s/%s did not become Ready", namespace, name)
}

// WaitForPeerAuthenticationReadyE waits for the resource condition to be satisfied.
func WaitForPeerAuthenticationReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var peerAuthentication *istiosecurityv1.PeerAuthentication
		peerAuthentication, err := istioClient.SecurityV1().PeerAuthentications(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		if peerAuthentication.Status.Conditions != nil {
			return istioConditionReady(t, &peerAuthentication.Status), nil
		}
		return false, nil
	})
}
