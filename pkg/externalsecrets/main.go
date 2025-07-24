package externalsecrets

import (
	esov1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewESOClient creates and returns a new controller-runtime client for interacting with ExternalSecrets resources.
// It initializes a runtime scheme and adds the ExternalSecrets API types to it.
// Returns the client or an error if the client could not be created.
//
// Parameters:
//
//	cfg - The Kubernetes REST configuration to use for the client.
//
// Returns:
//
//	client.Client - The initialized client for ExternalSecrets resources.
//	error         - An error if the client could not be created.
func NewESOClient(cfg *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	_ = esov1.AddToScheme(scheme)
	return client.New(cfg, client.Options{Scheme: scheme})
}
