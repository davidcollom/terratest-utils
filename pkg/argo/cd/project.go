package cd

import (
	"context"
	"testing"
	"time"

	argocdv1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListAppProjects retrieves a list of Argo CD AppProject resources in the specified namespace.
//
// Parameters:
//   - t: The testing context used for logging and error handling.
//   - options: The kubectl options used to configure access to the Kubernetes cluster.
//   - namespace: The namespace from which to list AppProjects.
//
// Returns:
//   - A slice of AppProject resources found in the given namespace.
//
// This function will fail the test if it cannot create the Argo CD client or if it fails to list the AppProjects.
func ListAppProjects(t *testing.T, options *k8s.KubectlOptions, namespace string) []argocdv1alpha1.AppProject {
	t.Helper()

	client, err := NewArgoCDClient(t, options)
	require.NoError(t, err, "Failed to create Argo CD clientset")

	ctx := t.Context()
	appProjectList, err := client.ArgoprojV1alpha1().AppProjects(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list AppProjects in namespace %s", namespace)

	return appProjectList.Items
}

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
