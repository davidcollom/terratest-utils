package workflows

import (
	"testing"

	workflowv1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListArgoWorkflowTaskResults(t *testing.T, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.WorkflowTaskResult {
	t.Helper()

	client, err := NewArgoWorkflowsClient(t, options)
	require.NoError(t, err, "Failed to create Argo Workflows clientset")

	ctx := t.Context()
	workflowTemplateList, err := client.ArgoprojV1alpha1().WorkflowTaskResults(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list WorkflowTemplates in namespace %s", namespace)

	return workflowTemplateList.Items
}

func ListArgoWorkflowTaskSet(t *testing.T, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.WorkflowTaskSet {
	t.Helper()

	client, err := NewArgoWorkflowsClient(t, options)
	require.NoError(t, err, "Failed to create Argo Workflows clientset")

	ctx := t.Context()
	workflowTemplateList, err := client.ArgoprojV1alpha1().WorkflowTaskSets(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list WorkflowTemplates in namespace %s", namespace)

	return workflowTemplateList.Items
}
