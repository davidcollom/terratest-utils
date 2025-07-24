// Package events provides Terratest-style helpers for testing Argo Events resources,
// including EventSources and Sensors. These helpers use client-go polling to wait for
// resources to report a Ready condition, ensuring event-driven workflows are correctly configured.
package events

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
