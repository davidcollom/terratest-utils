# terratest-utils

Terratest-compatible helper libraries for testing Kubernetes-ecosystem resources.

See [.github/copilot-instructions.md](.github/copilot-instructions.md) for the full coding conventions and API standards that all agents must follow.

## Quick Reference

### Testing Parameter

All non-test exported functions use:

```go
import "github.com/gruntwork-io/terratest/modules/testing"

func FooBar(t testing.TestingT, options *k8s.KubectlOptions, ...) ReturnType
```

Never `*testing.T`. Never `t.Helper()`, `t.Context()`, or `t.Cleanup()` in module files.

### Function Pair Pattern

Every resource action is exposed as a pair:

```go
func ListFoo(t testing.TestingT, ...)            []Foo          // fails test on error
func ListFooE(t testing.TestingT, ...)           ([]Foo, error) // returns error

func GetFoo(t testing.TestingT, ..., name string) *Foo
func GetFooE(t testing.TestingT, ..., name string) (*Foo, error)

func WaitForFooReady(t testing.TestingT, ..., timeout time.Duration)
func WaitForFooReadyE(t testing.TestingT, ..., timeout time.Duration) error
```

### Packages

| Package | Domain |
| --- | --- |
| `pkg/k8s` | Core Kubernetes (CRD, StatefulSet) + `KubectlOptions` alias |
| `pkg/argo/cd` | ArgoCD |
| `pkg/argo/events` | Argo Events |
| `pkg/argo/rollouts` | Argo Rollouts |
| `pkg/argo/workflows` | Argo Workflows |
| `pkg/certmanager` | cert-manager |
| `pkg/externalsecrets` | External Secrets Operator |
| `pkg/flux` | Flux v2 |
| `pkg/istio` | Istio |
| `pkg/linkerd` | Linkerd |
| `pkg/utils` | Shared REST config helper |
| `pkg/velero` | Velero |
