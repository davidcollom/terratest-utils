package certmanager

import (
	"testing"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmclientset "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	fakecm "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned/fake"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"k8s.io/apimachinery/pkg/runtime"
)

// k8soptions a global k8s.KubectlOptions instance to be used within many tests..
var k8soptions = &k8s.KubectlOptions{}

// NewTestClient creates a new test client with the given objects.
func NewTestClient(t *testing.T, objs ...runtime.Object) cmclientset.Interface {
	// Register everything to scheme
	scheme := runtime.NewScheme()
	_ = cmv1.AddToScheme(scheme)

	// Create a fake CertManager Client
	client := fakecm.NewSimpleClientset(objs...)

	// Override the function to return our expected objects
	NewClient = func(t *testing.T, options *k8s.KubectlOptions) (cmclientset.Interface, error) {
		return client, nil
	}
	// Ensure we have a cleanup to reset the function!
	t.Cleanup(func() {
		NewClient = newClient
	})

	return client
}
