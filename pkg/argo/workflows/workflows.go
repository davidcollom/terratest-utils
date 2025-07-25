// Package workflows provides Terratest-style helpers for testing Argo Workflows.
// These include functions to wait for specific workflow phases (Succeeded, Failed, etc.),
// and utilities to assert workflow conditions using the Argo Workflows clientset.
package workflows

import (
	"testing"
	"time"

	workflowv1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	workflowsClientSet "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/davidcollom/terratest-utils/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

func ListArgoWorkflows(t *testing.T, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.Workflow {
	t.Helper()

	client, err := NewArgoWorkflowsClient(t, options)
	require.NoError(t, err, "Failed to create Argo Workflows clientset")

	ctx := t.Context()
	workflowList, err := client.ArgoprojV1alpha1().Workflows(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list Workflows in namespace %s", namespace)

	return workflowList.Items
}

// NewArgoWorkflowsClient creates a new Argo Workflows client using the provided testing context and Kubernetes options.
// It returns an implementation of the workflowv1alpha1.Interface for interacting with Argo Workflows resources.
// If the provided KubectlOptions does not include a RestConfig, it attempts to generate one.
// Returns an error if the client cannot be created.
func NewArgoWorkflowsClient(t *testing.T, options *k8s.KubectlOptions) (workflowsClientSet.Interface, error) {
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

	return workflowsClientSet.NewForConfig(cfg)
}

func ListWorkflows(t *testing.T, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.Workflow {
	t.Helper()

	client, err := NewArgoWorkflowsClient(t, options)
	require.NoError(t, err, "Failed to create Argo Workflows clientset")

	ctx := t.Context()
	workflowList, err := client.ArgoprojV1alpha1().Workflows(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list Workflows in namespace %s", namespace)

	return workflowList.Items
}

func WaitForWorkflowRunning(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	WaitForWorkflowPhase(t, options, name, namespace, workflowv1alpha1.WorkflowRunning, timeout)
}

// WaitForWorkflowPending waits until the specified Argo workflow reaches the "Pending" phase within the given timeout.
// It uses the provided testing context, kubectl options, workflow name, and namespace.
// If the workflow does not reach the "Pending" phase within the timeout, the test will fail.
func WaitForWorkflowPending(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	WaitForWorkflowPhase(t, options, name, namespace, workflowv1alpha1.WorkflowPending, timeout)
}
