// Package events provides Terratest-style helpers for testing Argo Events resources,
// including EventSources and Sensors. These helpers use client-go polling to wait for
// resources to report a Ready condition, ensuring event-driven workflows are correctly configured.
package events

import (
	"testing"

	argoclientset "github.com/argoproj/argo-events/pkg/client/clientset/versioned"

	"github.com/davidcollom/terratest-utils/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

// HasReadyCondition checks if the provided slice of metav1.Condition contains a condition
// of the specified type with a status of "True". It returns true if such a condition is found,
// otherwise returns false.
//
// conds:        Slice of metav1.Condition to search through.
// expectedType: The condition type to look for.
//
// Returns true if a condition with the expected type and a status of "True" exists, false otherwise.
func HasReadyCondition(conds []metav1.Condition, expectedType string) bool {
	for _, cond := range conds {
		if cond.Type == expectedType && cond.Status == "True" {
			return true
		}
	}
	return false
}

// NewArgoEventsClient creates a new Argo Events client using the provided testing context and Kubernetes options.
// It retrieves the Kubernetes REST configuration from the provided options or generates one if not present.
// Returns an Argo Events clientset interface and an error if the client could not be created.
//
// Parameters:
//   - t: The testing context, used for logging and error handling.
//   - options: The Kubernetes KubectlOptions containing cluster access information.
//
// Returns:
//   - argoclientset.Interface: The Argo Events clientset interface for interacting with Argo Events resources.
//   - error: An error if the client could not be created.
func NewArgoEventsClient(t *testing.T, options *k8s.KubectlOptions) (argoclientset.Interface, error) {
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

	return argoclientset.NewForConfig(cfg)
}
