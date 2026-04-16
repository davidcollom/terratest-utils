package cd

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

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
// ListAppProjects lists matching resources.
func ListAppProjects(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []argocdv1alpha1.AppProject {
	projects, err := ListAppProjectsE(t, options, namespace)
	require.NoError(t, err, "Failed to list AppProjects in namespace %s", namespace)
	return projects
}

// ListAppProjectsE retrieves a list of Argo CD AppProject resources in the specified namespace.
// It returns an error to the caller instead of failing the test directly.
// ListAppProjectsE lists matching resources.
func ListAppProjectsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]argocdv1alpha1.AppProject, error) {
	client, err := NewArgoCDClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	appProjectList, err := client.ArgoprojV1alpha1().AppProjects(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return appProjectList.Items, nil
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
// WaitForAppProjectExists waits for the resource condition to be satisfied.
func WaitForAppProjectExists(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForAppProjectExistsE(t, options, name, namespace, timeout)
	require.NoError(t, err, "AppProject %s/%s did not appear", namespace, name)
}

// WaitForAppProjectExistsE waits until an Argo CD AppProject with the specified name exists in the given namespace.
func WaitForAppProjectExistsE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	client, err := NewArgoCDClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := client.ArgoprojV1alpha1().AppProjects(namespace).Get(ctx, name, metav1.GetOptions{})
		return err == nil, nil
	})
}
