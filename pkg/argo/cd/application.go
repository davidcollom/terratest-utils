package cd

import (
	"context"
	"testing"
	"time"

	argocd "github.com/argoproj/argo-cd/v3/pkg/client/clientset/versioned"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForApplicationHealthyAndSynced waits until the specified Argo CD Application resource
// in the given namespace reaches both "Healthy" and "Synced" status within the provided timeout.
// It polls the Application status every 2 seconds using the Argo CD client and fails the test
// if the desired state is not achieved within the timeout period.
//
// Parameters:
//
//	t        - The testing context.
//	options  - Kubectl options containing the Kubernetes REST config.
//	name     - The name of the Argo CD Application.
//	namespace- The namespace where the Application resides.
//	timeout  - The maximum duration to wait for the Application to become Healthy and Synced.
//
// Fails the test if the Application does not reach the desired state within the timeout.
func WaitForApplicationHealthyAndSynced(t *testing.T, options k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := argocd.NewForConfig(options.RestConfig)
	require.NoError(t, err, "Unable to create Argo CD client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		app, err := client.ArgoprojV1alpha1().Applications(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		if IsApplicationHealthyAndSynced(app) {
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("Application %s/%s did not become Healthy & Synced: %v", namespace, name, err)
	}
}
