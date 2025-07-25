package workflows

import (
	"testing"

	workflowv1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListArgoWorkflowTaskResults retrieves a list of Argo WorkflowTaskResult resources from the specified namespace.
// It uses the provided testing context and kubectl options to create an Argo Workflows client.
// The function fails the test if the client cannot be created or if listing the WorkflowTaskResults fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options used to configure the client.
//   - namespace: The Kubernetes namespace from which to list WorkflowTaskResults.
//
// Returns:
//   - A slice of WorkflowTaskResult resources found in the specified namespace.
func ListArgoWorkflowTaskResults(t *testing.T, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.WorkflowTaskResult {
	t.Helper()

	client, err := NewArgoWorkflowsClient(t, options)
	require.NoError(t, err, "Failed to create Argo Workflows clientset")

	ctx := t.Context()
	workflowTemplateList, err := client.ArgoprojV1alpha1().WorkflowTaskResults(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list WorkflowTemplates in namespace %s", namespace)

	return workflowTemplateList.Items
}

// ListArgoWorkflowTaskSet retrieves all Argo WorkflowTaskSet resources in the specified namespace.
// It uses the provided testing context and kubectl options to create an Argo Workflows client,
// then lists the WorkflowTaskSets in the given namespace. If any error occurs during client creation
// or listing, the test will fail. Returns a slice of WorkflowTaskSet objects found in the namespace.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list WorkflowTaskSets.
//
// Returns:
//   - A slice of WorkflowTaskSet objects present in the specified namespace.
func ListArgoWorkflowTaskSet(t *testing.T, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.WorkflowTaskSet {
	t.Helper()

	client, err := NewArgoWorkflowsClient(t, options)
	require.NoError(t, err, "Failed to create Argo Workflows clientset")

	ctx := t.Context()
	workflowTemplateList, err := client.ArgoprojV1alpha1().WorkflowTaskSets(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list WorkflowTemplates in namespace %s", namespace)

	return workflowTemplateList.Items
}
