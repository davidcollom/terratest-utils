package flux

import (
	"context"
	"testing"
	"time"

	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/util/wait"
)

// ListBuckets retrieves a list of Flux Buckets in the specified Kubernetes namespace.
// It uses the provided testing context and kubectl options to create a Flux client,
// then lists all Bucket resources within the given namespace. The function fails the test
// if the client cannot be created or if listing the Buckets fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the cluster.
//   - namespace: The Kubernetes namespace to search for Buckets.
//
// Returns:
//   - A slice of sourcev1.Bucket objects found in the specified namespace.
func ListBuckets(t *testing.T, options *k8s.KubectlOptions, namespace string) []sourcev1.Bucket {
	t.Helper()

	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	ctx := t.Context()
	var buckets sourcev1.BucketList
	err = fluxclient.List(ctx, &buckets, client.InNamespace(namespace))
	require.NoError(t, err, "Failed to list Buckets in namespace %s", namespace)

	return buckets.Items
}

// WaitForBucketReady waits until the specified Flux Bucket resource reaches the "Ready" condition within the given timeout.
// It polls the Kubernetes API at regular intervals to check the Bucket's status.
// If the Bucket does not become ready within the timeout, the test fails.
// Parameters:
//
//	t        - The testing context.
//	options  - Kubectl options containing the REST config for Kubernetes API access.
//	name     - The name of the Bucket resource.
//	namespace- The namespace where the Bucket resource is located.
//	timeout  - The maximum duration to wait for the Bucket to become ready.
func WaitForBucketReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {

		var bucket sourcev1.Bucket
		err = fluxclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &bucket)
		if err != nil {
			return false, nil
		}
		return hasReadyCondition(bucket.Status.Conditions), nil
	})

	if err != nil {
		t.Fatalf("Bucket %s/%s did not become Ready: %v", namespace, name, err)
	}
}
