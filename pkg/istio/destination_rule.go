package istio

import (
	"context"
	"github.com/gruntwork-io/terratest/modules/testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	isitonetworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListDestinationRules retrieves all Istio DestinationRule resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list DestinationRules from.
//
// Returns:
//   - A slice of pointers to DestinationRule objects found in the namespace.
func ListDestinationRules(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*isitonetworkingv1alpha3.DestinationRule {
	istioClient := NewClient(t, options)

	ctx := context.Background()
	destinationRules, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).List(ctx, v1meta.ListOptions{})
	require.NoError(t, err, "Failed to list Destination Rules in namespace %s", namespace)

	return destinationRules.Items
}

// WaitForDestinationRuleReady waits until the specified DestinationRule in the given namespace is Ready or the timeout is reached.
// It polls the DestinationRule status every 2 seconds and checks for the Ready condition.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the DestinationRule to check.
//   - namespace: The namespace of the DestinationRule.
//   - timeout: The maximum duration to wait for the resource to become Ready.
func WaitForDestinationRuleReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := context.Background()
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var destinationRule *isitonetworkingv1alpha3.DestinationRule
		destinationRule, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		if destinationRule.Status.Conditions != nil {
			return istioConditionReady(t, &destinationRule.Status), nil
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("DestinationRule %s/%s did not become Ready: %v", namespace, name, err)
	}
}
