// Package k8s provides utilities for working with Kubernetes resources in tests.
//
// This package includes helper functions for interacting with various Kubernetes resources such as StatefulSets,
// Deployments, Services, and more. It is designed to simplify common testing tasks, such as retrieving resources,
// waiting for resource readiness, and asserting resource states, by providing convenient wrappers around the
// Kubernetes client-go library.
//
// Example usage:
//
//	import (
//	    "testing"
//	    "github.com/davidcollom/terratest-utils/pkg/k8s"
//	)
//
//	func TestStatefulSetReady(t *testing.T) {
//	    options := k8s.NewKubectlOptions("", "", "default")
//	    k8s.WaitForStatefulSetReady(t, options, "my-statefulset", "default", 5*time.Minute)
//	}
//
// The package is intended for use in automated tests and supports both fatal and non-fatal error handling patterns.
package k8s
