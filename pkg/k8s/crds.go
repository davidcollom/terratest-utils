package k8s

import (
	"testing"

	terrak8s "github.com/gruntwork-io/terratest/modules/k8s"

	apixv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apixclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/davidcollom/terratest-utils/pkg/utils"
	"github.com/stretchr/testify/require"
)

func GetCustomResourceDefinition(t *testing.T, options *terrak8s.KubectlOptions, crdName string) *apixv1.CustomResourceDefinition {
	t.Helper()

	crd, err := GetCustomResourceDefinitionE(t, options, crdName)
	require.NoError(t, err)
	return crd
}

func GetCustomResourceDefinitionE(t *testing.T, options *terrak8s.KubectlOptions, crdName string) (*apixv1.CustomResourceDefinition, error) {
	t.Helper()

	restConfig, err := utils.GetRestConfigE(t, options)
	if err != nil {
		return nil, err
	}

	client, err := apixclientset.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return client.ApiextensionsV1().CustomResourceDefinitions().Get(t.Context(), crdName, metav1.GetOptions{})
}
