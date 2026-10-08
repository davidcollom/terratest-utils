# terratest-utils

Terratest-compatible helper libraries for testing Kubernetes-ecosystem resources, split into 12 independent Go modules on top of Terratest v2.

See [.github/copilot-instructions.md](.github/copilot-instructions.md) for the full coding conventions and API standards that all agents must follow.

## Quick Reference

### Testing Parameter

All non-test exported functions use:

```go
import "github.com/gruntwork-io/terratest/modules/core/v2/testing"

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

### Modules

The repo is a monorepo of independent Go modules wired together by `go.work`. Each `pkg/<dir>` (and each `pkg/argo/<dir>`) is its own importable module.

| Module | Import path | Domain |
| --- | --- | --- |
| `pkg/utils` | `github.com/davidcollom/terratest-utils/pkg/utils` | Shared REST config helper |
| `pkg/k8s` | `github.com/davidcollom/terratest-utils/pkg/k8s` | Core Kubernetes (CRD, StatefulSet) + `KubectlOptions` alias |
| `pkg/certmanager` | `github.com/davidcollom/terratest-utils/pkg/certmanager` | cert-manager |
| `pkg/externalsecrets` | `github.com/davidcollom/terratest-utils/pkg/externalsecrets` | External Secrets Operator |
| `pkg/flux` | `github.com/davidcollom/terratest-utils/pkg/flux` | Flux v2 |
| `pkg/istio` | `github.com/davidcollom/terratest-utils/pkg/istio` | Istio |
| `pkg/linkerd` | `github.com/davidcollom/terratest-utils/pkg/linkerd` | Linkerd |
| `pkg/velero` | `github.com/davidcollom/terratest-utils/pkg/velero` | Velero |
| `pkg/argo/cd` | `github.com/davidcollom/terratest-utils/pkg/argo/cd` | ArgoCD |
| `pkg/argo/events` | `github.com/davidcollom/terratest-utils/pkg/argo/events` | Argo Events |
| `pkg/argo/rollouts` | `github.com/davidcollom/terratest-utils/pkg/argo/rollouts` | Argo Rollouts |
| `pkg/argo/workflows` | `github.com/davidcollom/terratest-utils/pkg/argo/workflows` | Argo Workflows |
