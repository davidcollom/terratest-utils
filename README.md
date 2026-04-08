# Terratest Utils

This package provides a set of helpers for [Terratest](https://terratest.gruntwork.io/) and general integration testing. It is designed to help validate that platforms and deployments have been successfully provisioned and configured.

## Features

- Utilities for testing Kubernetes resources, including cert-manager, external-secrets, ArgoCD, and Flux
- Functions to check readiness, status, and correctness of deployed resources
- Simplifies writing robust integration tests for cloud-native platforms
- Can be used with Terratest or standalone in Go test suites

## Usage

Import the relevant package(s) in your Terratest or Go integration tests:

```go
import (
    "github.com/davidcollom/terratest-utils/pkg/certmanager"
    "github.com/davidcollom/terratest-utils/pkg/externalsecrets"
    // ...other helpers
)
```

Use the provided functions to validate resources, e.g.:

```go
import (
    "testing"
    "time"

    "github.com/davidcollom/terratest-utils/pkg/certmanager"
    "github.com/davidcollom/terratest-utils/pkg/flux"
    "github.com/gruntwork-io/terratest/modules/k8s"
)

func TestPlatform(t *testing.T) {
    options := k8s.NewKubectlOptions("", "", "default")

    // Wait for a cert-manager Certificate to become Ready
    certmanager.WaitForCertificateReady(t, options, "my-cert", "default", 5*time.Minute)

    // List all Flux HelmReleases in a namespace
    releases := flux.ListHelmReleases(t, options, "flux-system")

    // Use the E variant to handle errors yourself
    cert, err := certmanager.NewClient(t, options)
    // ...
}
```

Every helper is available in two forms:
- `VerbResource(t, options, ...)` — fails the test on error
- `VerbResourceE(t, options, ...)` — returns `(value, error)` for custom handling

## Structure

| Package | Description |
|---|---|
| `pkg/argo/cd` | Helpers for ArgoCD Application, ApplicationSet, and AppProject resources |
| `pkg/argo/events` | Helpers for Argo Events — EventBus, EventSource, Sensor |
| `pkg/argo/rollouts` | Helpers for Argo Rollouts |
| `pkg/argo/workflows` | Helpers for Argo Workflows, CronWorkflows, WorkflowTemplates, and WorkflowPhases |
| `pkg/certmanager` | Helpers for cert-manager Certificate, Issuer, ClusterIssuer, CertificateRequest, Order, and Challenge resources |
| `pkg/externalsecrets` | Helpers for External Secrets Operator — ExternalSecret, ClusterExternalSecret, SecretStore, ClusterSecretStore, PushSecret |
| `pkg/flux` | Helpers for Flux v2 — HelmRelease, HelmRepository, HelmChart, GitRepository, Kustomization, Bucket, OCIRepository |
| `pkg/istio` | Helpers for Istio networking (Gateway, VirtualService, DestinationRule, ServiceEntry, Sidecar, EnvoyFilter, WorkloadEntry, WorkloadGroup) and security (AuthorizationPolicy, PeerAuthentication, RequestAuthentication) |
| `pkg/k8s` | Core Kubernetes helpers — CRD, StatefulSet — plus the `KubectlOptions` alias |
| `pkg/linkerd` | Helpers for Linkerd policy and traffic resources — Server, ServerAuthorization, AuthorizationPolicy, HTTPRoute, MeshTLSAuthentication, NetworkAuthentication, ServiceProfile, TrafficSplit |
| `pkg/utils` | Shared utility — REST config helper used by all domain packages |
| `pkg/velero` | Helpers for Velero Backup, Restore, Schedule, BackupStorageLocation |

## Purpose

These helpers are intended to:

- Accelerate writing integration tests for Kubernetes platforms
- Provide reusable checks for resource readiness and correctness
- Help ensure deployments are successful and meet expected criteria

## License

MIT
