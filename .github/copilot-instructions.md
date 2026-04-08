# terratest-utils — Copilot Instructions

## Project Purpose

`terratest-utils` is a collection of **Terratest-compatible helper libraries** for testing Kubernetes-ecosystem resources. Each sub-package under `pkg/` provides typed, idiomatic Go helpers that work alongside [gruntwork-io/terratest](https://github.com/gruntwork-io/terratest) without depending on its concrete `*testing.T`.

## Repository Layout

```
pkg/
  argo/
    cd/            ArgoCD Application, ApplicationSet, Project
    events/        Argo Events — EventBus, EventSource, Sensor
    rollouts/      Argo Rollouts
    workflows/     Argo Workflows, CronWorkflows, WorkflowPhase, etc.
  certmanager/     cert-manager Certificate, Issuer, ClusterIssuer, Order, Challenge
  externalsecrets/ External Secrets Operator
  flux/            Flux v2 — HelmRelease, GitRepository, Kustomization, etc.
  istio/           Istio networking and security resources
  k8s/             Core Kubernetes helpers (CRDs, StatefulSets) + KubectlOptions alias
  linkerd/         Linkerd policy and traffic resources
  utils/           Shared utilities (REST config)
  velero/          Velero Backup, Restore, Schedule, BackupStorageLocation
```

## API Conventions (Non-Negotiable)

Every exported public-facing function must follow the Terratest job-helper pattern exactly.

### Function Pair Pattern

All resource interactions are exposed as a pair:

```go
// Non-E variant: fails the test on error (uses require.NoError)
func VerbResource(t testing.TestingT, options *k8s.KubectlOptions, ...) ReturnType

// E variant: returns error to the caller
func VerbResourceE(t testing.TestingT, options *k8s.KubectlOptions, ...) (ReturnType, error)
```

### Standard Verbs

| Verb prefix | Behaviour |
|---|---|
| `List`       | List all resources (with filters), non-E fails on error |
| `Get`        | Fetch single resource by name, non-E fails on error |
| `WaitFor`    | Poll until ready/condition met, non-E fails on timeout |
| `WaitUntil`  | Same intent as WaitFor — use `WaitUntilXSucceedE` for retry-loop style |
| `Create`     | Create a resource, non-E fails on error |
| `Validate`   | Assert a property (no E variant needed — calls Fatalf directly) |
| `Is`         | Pure boolean predicate — no `t` param, no error return |
| `New*Client` | Construct a typed client — always exported as `var New*Client = new*Client` |

### Testing Parameter

Always use `testing.TestingT` from `github.com/gruntwork-io/terratest/modules/testing`, imported as `testing` (not aliased):

```go
import "github.com/gruntwork-io/terratest/modules/testing"

func ListFoo(t testing.TestingT, options *k8s.KubectlOptions, ...) []Foo {
```

**Never** use stdlib `*testing.T` in non-test module files.

### Prohibited in non-test module files

- `t.Helper()` — not on `testing.TestingT` interface
- `t.Context()` — not on `testing.TestingT` interface
- `t.Cleanup()` — not on `testing.TestingT` interface
- `t.Run()` — not on `testing.TestingT` interface

Use `context.Background()` wherever you need a context.

### KubectlOptions

All resource helpers accept `*k8s.KubectlOptions` where `k8s` is `github.com/gruntwork-io/terratest/modules/k8s`. The `pkg/k8s` package exposes a re-export alias:

```go
// pkg/k8s/aliases.go
type KubectlOptions = terrak8s.KubectlOptions
var NewKubectlOptions = terrak8s.NewKubectlOptions
```

Other packages (certmanager, flux, etc.) import `github.com/gruntwork-io/terratest/modules/k8s` directly — they do **not** depend on the local `pkg/k8s` package for `KubectlOptions`.

### Client Construction Pattern

Each domain package has an internal constructor exported as a function variable so tests can override it:

```go
var NewClient = newClient   // exported var for test injection

func newClient(t testing.TestingT, options *k8s.KubectlOptions) (SomeClientset, error) {
    cfg, err := utils.GetRestConfigE(t, options)
    if err != nil {
        return nil, err
    }
    return someclientset.NewForConfig(cfg)
}
```

### Polling Pattern

`WaitFor*` functions use `k8s.io/apimachinery/pkg/util/wait.PollUntilContextTimeout` with 2-second poll interval:

```go
err = wait.PollUntilContextTimeout(context.Background(), 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
    resource, err := client.Get(ctx, name, metav1.GetOptions{})
    if err != nil {
        return false, nil // retry on transient errors
    }
    return isReady(resource), nil
})
if err != nil {
    t.Fatalf("Resource %s/%s was not ready in time: %v", namespace, name, err)
}
```

### Return types

- `List*` → `[]ResourceType` (slice, not pointer to list object)
- `Get*` → `*ResourceType`
- `WaitFor*` / `WaitUntil*` (non-E) → nothing (void)
- `WaitFor*E` / `WaitUntil*E` → `error`
- `Create*` → `*ResourceType`
- `Create*E` → `(*ResourceType, error)`
- `Is*` → `bool`

## Error Handling

- Non-E functions call `require.NoError(t, err)` or `t.Fatalf(...)` — they never return errors.
- E functions return the raw error, never call `t.Fatalf` themselves.

## Imports

- Use `context.Background()` for all in-production k8s API calls.
- Use `github.com/stretchr/testify/require` for assertions in non-E wrappers.
- Do **not** duplicate test file imports into module files.

## Test Files (`*_test.go`)

Test files may use stdlib `*testing.T` (they cannot use `testing.TestingT` for `t.Cleanup` etc.) and should alias it as `gotesting` only if the file also imports Terratest testing:

```go
import (
    gotesting "testing"
    "github.com/gruntwork-io/terratest/modules/testing"
)
```

Override client constructor variables in tests using `t.Cleanup` to restore them:

```go
NewClient = func(t testing.TestingT, options *k8s.KubectlOptions) (SomeInterface, error) {
    return fakeClient, nil
}
t.Cleanup(func() { NewClient = newClient })
```

## Adding a New Package

1. Create `pkg/<domain>/` with a `<domain>.go` file containing the `NewClient` constructor.
2. Follow the function pair pattern for every resource type.
3. Import `"github.com/gruntwork-io/terratest/modules/testing"` as `testing`.
4. Use `*k8s.KubectlOptions` from `github.com/gruntwork-io/terratest/modules/k8s`.
5. No `t.Helper()`, `t.Context()`, or stdlib `*testing.T` parameters.

## Adding a New Function to an Existing Package

1. Add the E-variant first (implement real logic, return error).
2. Add the non-E wrapper that calls the E-variant with `require.NoError`.
3. Match signature order: `(t testing.TestingT, options *k8s.KubectlOptions, name string, namespace string, ...)`.

## Dependencies

Key external dependencies and their purpose:

| Module | Purpose |
|---|---|
| `github.com/gruntwork-io/terratest` | `KubectlOptions`, k8s client helpers, retry, testing interface |
| `k8s.io/client-go` | Kubernetes client |
| `k8s.io/apimachinery` | Kubernetes types, `wait.PollUntilContextTimeout` |
| `sigs.k8s.io/controller-runtime` | Used by Flux, ExternalSecrets (controller-runtime client) |
| `github.com/stretchr/testify/require` | Assertions in non-E wrappers |
