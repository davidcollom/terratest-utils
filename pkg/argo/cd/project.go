package cd

import (
	"context"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForAppProjectExists waits until an Argo CD AppProject with the specified name exists in the given namespace.
// It polls the Kubernetes API at regular intervals until the AppProject is found or the timeout is reached.
// If the AppProject does not appear within the timeout, the test fails.
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the REST config for the Kubernetes cluster.
//   - name: The name of the AppProject to wait for.
//   - namespace: The namespace in which to look for the AppProject.
//   - timeout: The maximum duration to wait for the AppProject to appear.
func WaitForAppProjectExists(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := NewArgoCDClient(t, options)
	require.NoError(t, err, "Unable to create Argo CD client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := client.ArgoprojV1alpha1().AppProjects(namespace).Get(ctx, name, metav1.GetOptions{})
		return err == nil, nil
	})

	if err != nil {
		t.Fatalf("AppProject %s/%s did not appear: %v", namespace, name, err)
	}
}
