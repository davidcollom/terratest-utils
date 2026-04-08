package k8s

import terrak8s "github.com/gruntwork-io/terratest/modules/k8s"

// KubectlOptions is an alias to terratest's KubectlOptions type so package APIs can
// match terratest function signatures.
type KubectlOptions = terrak8s.KubectlOptions

// NewKubectlOptions creates a new KubectlOptions object.
var NewKubectlOptions = terrak8s.NewKubectlOptions
