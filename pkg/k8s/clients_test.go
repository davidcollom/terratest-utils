package k8s

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	apixv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apixcm "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apixfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"

	"k8s.io/apimachinery/pkg/runtime"
)

// NewAPIXTestClient creates a new test client with the given objects.
func NewAPIXTestClient(t *testing.T, objs []runtime.Object) apixcm.Interface {
	// Register everything to scheme
	scheme := runtime.NewScheme()
	_ = apixv1.AddToScheme(scheme)

	// Create a fake CertManager Client
	client := apixfake.NewSimpleClientset(objs...)

	// Override the function to return our expected objects
	NewAPIXClient = func(t *testing.T, options *k8s.KubectlOptions) (apixcm.Interface, error) {
		return client, nil
	}
	// Ensure we have a cleanup to reset the function!
	t.Cleanup(func() {
		NewAPIXClient = newAPIXClient
	})

	return client
}
