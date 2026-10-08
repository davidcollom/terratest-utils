# Terratest Utils

Terratest-compatible helper libraries for testing Kubernetes-ecosystem resources, split into 12 independent Go modules on top of Terratest v2.

## Modules

Each module is a standalone Go module that can be imported independently. All helpers follow the Terratest job-helper pattern (`VerbResource` and `VerbResourceE` pairs).

| Module | Import path | Coverage |
| --- | --- | --- |
| utils | `github.com/davidcollom/terratest-utils/pkg/utils` | Shared REST config helper used by every domain module |
| k8s | `github.com/davidcollom/terratest-utils/pkg/k8s` | Core Kubernetes (CRD, StatefulSet) and a `KubectlOptions` alias over Terratest v2 |
| certmanager | `github.com/davidcollom/terratest-utils/pkg/certmanager` | cert-manager — Certificate, Issuer, ClusterIssuer, CertificateRequest, Order, Challenge |
| externalsecrets | `github.com/davidcollom/terratest-utils/pkg/externalsecrets` | External Secrets Operator — ExternalSecret, ClusterExternalSecret, SecretStore, ClusterSecretStore, PushSecret |
| flux | `github.com/davidcollom/terratest-utils/pkg/flux` | Flux v2 — HelmRelease, HelmRepository, HelmChart, GitRepository, Kustomization, Bucket, OCIRepository |
| istio | `github.com/davidcollom/terratest-utils/pkg/istio` | Istio networking (Gateway, VirtualService, DestinationRule, ServiceEntry, Sidecar, EnvoyFilter, WorkloadEntry, WorkloadGroup) and security (AuthorizationPolicy, PeerAuthentication, RequestAuthentication) |
| linkerd | `github.com/davidcollom/terratest-utils/pkg/linkerd` | Linkerd — Server, ServerAuthorization, AuthorizationPolicy, HTTPRoute, MeshTLSAuthentication, NetworkAuthentication, ServiceProfile, TrafficSplit |
| velero | `github.com/davidcollom/terratest-utils/pkg/velero` | Velero — Backup, Restore, Schedule, BackupStorageLocation |
| argo/cd | `github.com/davidcollom/terratest-utils/pkg/argo/cd` | ArgoCD — Application, ApplicationSet, AppProject |
| argo/events | `github.com/davidcollom/terratest-utils/pkg/argo/events` | Argo Events — EventBus, EventSource, Sensor |
| argo/rollouts | `github.com/davidcollom/terratest-utils/pkg/argo/rollouts` | Argo Rollouts |
| argo/workflows | `github.com/davidcollom/terratest-utils/pkg/argo/workflows` | Argo Workflows — Workflows, CronWorkflows, WorkflowTemplates, WorkflowPhases |

## Usage

Import the modules you need in your Terratest or Go integration tests. All helpers accept `testing.TestingT` from Terratest v2 and `*k8s.KubectlOptions` from Terratest v2.

```go
import (
    "time"

    "github.com/davidcollom/terratest-utils/pkg/certmanager"
    "github.com/davidcollom/terratest-utils/pkg/flux"
    "github.com/gruntwork-io/terratest/modules/core/v2/testing"
    "github.com/gruntwork-io/terratest/modules/k8s/v2"
)

func TestPlatform(t *testing.T) {
    options := k8s.NewKubectlOptions("", "", "default")

    // Wait for a cert-manager Certificate to become Ready
    certmanager.WaitForCertificateReady(t, options, "my-cert", "default", 5*time.Minute)

    // List all Flux HelmReleases in a namespace
    releases := flux.ListHelmReleases(t, options, "flux-system")

    // Use the E variant to handle errors yourself
    client, err := certmanager.NewClient(t, options)
    // ...
}
```

Every helper is available in two forms:

- `VerbResource(t, options, ...)` — fails the test on error
- `VerbResourceE(t, options, ...)` — returns `(value, error)` for custom handling

## Local development

The repo is a monorepo of 12 independent Go modules wired together by the root `go.work` file. Use the workspace to build and test every module at once, or `cd` into a single module to work on it in isolation.

Build, vet, and test everything from the repo root:

```sh
go build ./...
go vet ./...
go test ./...
```

Work on a single module:

```sh
cd pkg/certmanager
go build ./...
go test ./...
```

The root `go.mod` is a marker module for the workspace and contains no Go source. Each `pkg/<dir>` (and each `pkg/argo/<dir>`) has its own `go.mod` and is the unit of versioning and release.

## Release process

Releases are driven by the `Release` GitHub Actions workflow (`.github/workflows/release.yaml`), triggered manually via `workflow_dispatch`.

1. **Trigger** — Run the workflow from the Actions tab and pick a release type (`auto`, `major`, `minor`, or `patch`).
2. **Version computation** — The workflow installs [`svu`](https://github.com/caarlos0/svu) and computes the next SemVer version from the git history.
3. **Version bump** — It runs `./scripts/bump-version.sh <version>`, which updates every submodule's `go.mod` so cross-module `require` directives point at the new version, and drops the local `replace` directives that only exist for development.
4. **Commit and tag** — The workflow commits the version bump, creates the exact version tag (e.g. `v0.0.12`), and creates the alias tags for the minor (`v0.0`) and major (`v0`) series so consumers can pin to a moving tag.
5. **Push and release** — Tags are pushed to `origin` and a GitHub Release is created with auto-generated notes.

The `bump-version.sh` script can also be run locally: `./scripts/bump-version.sh v0.0.12`.

## License

MIT
