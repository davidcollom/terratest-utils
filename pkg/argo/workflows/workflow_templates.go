package workflows

import (
	"testing"

	workflowv1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
)

// ListArgoWorkflowTemplates retrieves all Argo WorkflowTemplates in the specified namespace.
//
// It uses the provided testing context and kubectl options to create an Argo Workflows client,
// then lists all WorkflowTemplates in the given namespace. The function fails the test if any
// errors occur during client creation or listing.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for accessing the Kubernetes cluster.
//   - namespace: The namespace from which to list WorkflowTemplates.
//
// Returns:
//   - A slice of WorkflowTemplate objects found in the specified namespace.
func ListArgoWorkflowTemplates(t *testing.T, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.WorkflowTemplate {
	t.Helper()

	client, err := NewArgoWorkflowsClient(t, options)
	require.NoError(t, err, "Failed to create Argo Workflows clientset")

	ctx := t.Context()
	workflowTemplateList, err := client.ArgoprojV1alpha1().WorkflowTemplates(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list WorkflowTemplates in namespace %s", namespace)

	return workflowTemplateList.Items
}

// ListArgoClusterWorkflowTemplates retrieves all Argo ClusterWorkflowTemplates in the specified namespace using the provided KubectlOptions.
// It returns a slice of ClusterWorkflowTemplate objects. The function will fail the test if there is an error creating the client
// or listing the ClusterWorkflowTemplates.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options to use for connecting to the cluster.
//   - namespace: The namespace to list ClusterWorkflowTemplates from.
//
// Returns:
//   - A slice of ClusterWorkflowTemplate objects.
func ListArgoClusterWorkflowTemplates(t *testing.T, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.ClusterWorkflowTemplate {
	t.Helper()

	client, err := NewArgoWorkflowsClient(t, options)
	require.NoError(t, err, "Failed to create Argo Workflows clientset")

	ctx := t.Context()
	workflowTemplateList, err := client.ArgoprojV1alpha1().ClusterWorkflowTemplates().List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list WorkflowTemplates in namespace %s", namespace)

	return workflowTemplateList.Items
}
