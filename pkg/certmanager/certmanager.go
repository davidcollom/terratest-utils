package certmanager

import (
	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
)

// HasCondition checks if a slice of CertificateRequestCondition contains a condition
// with the specified type and status.
//
// Parameters:
//
//	conds    - Slice of CertificateRequestCondition to search.
//	condType - The condition type to look for.
//	status   - The condition status to match.
//
// Returns:
//
//	true if a condition with the specified type and status exists, false otherwise.
func HasCondition(conds []cmv1.CertificateRequestCondition, condType cmv1.CertificateRequestConditionType, status cmmetav1.ConditionStatus) bool {
	for _, cond := range conds {
		if cond.Type == condType && cond.Status == status {
			return true
		}
	}
	return false
}
