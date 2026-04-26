# ADR 001: Source of Truth Boundary Between Terraform and Git

## Status
Accepted

## Context

This repository manages an AKS-based Kubernetes platform.

We need a clear ownership boundary so the repo does not drift into a mix of
infrastructure provisioning, cluster configuration, and application delivery
all in the same place.

Without an explicit boundary, it becomes easy to:
- create Kubernetes resources in Terraform just because it is possible
- manually change cluster state instead of declaring it in Git
- blur the line between Azure infrastructure and Kubernetes desired state
- let CI/CD update the live cluster directly, bypassing Git history

This project needs a simple rule for where responsibility changes from
cloud infrastructure to cluster configuration.

## Decision

Terraform owns Azure infrastructure and AKS provisioning.

That means Terraform is responsible for things like:
- resource groups
- networking
- AKS cluster creation
- managed identities, node pools, and Azure-side dependencies

Git owns Kubernetes desired state.

That means Kubernetes manifests live in this repository and are applied from Git.
Examples include:
- namespaces
- platform-level Kubernetes objects
- later, application overlays and environment-specific manifests

ArgoCD will later reconcile Git state into the cluster.

ArgoCD will watch the Git-managed overlay structure, not the raw repo root and
not arbitrary chart sources mixed into the same boundary.

CI updates Git, not the cluster.

The pipeline may validate, format, and promote changes, but the source of truth
remains the Git repository. The cluster is a projection of that state.

## Consequences

This gives us a clean separation of responsibilities:
- Terraform provisions and changes cloud infrastructure
- Git defines Kubernetes configuration
- ArgoCD later becomes the reconciler for cluster state

The main benefits are:
- clear review history
- repeatable cluster configuration
- fewer accidental ownership conflicts
- easier onboarding for future contributors
- safer evolution toward GitOps

The tradeoff is that some changes require editing the right layer instead of
using whichever tool is most convenient. That is intentional. The boundary keeps
the platform understandable as it grows.

## Alternatives considered

### Put namespaces in Terraform
Rejected.

Terraform can create namespaces, but that would blur the line between Azure
infrastructure and Kubernetes desired state. It also makes later GitOps adoption
messier because the cluster baseline would already be split across two owners.

### Manage everything manually with kubectl
Rejected.

Manual changes do not scale, are hard to review, and are easy to drift from the
intended state stored in Git.

### Jump straight to ArgoCD before defining the baseline
Rejected.

ArgoCD is useful, but the project needs a minimal, easy-to-understand Git-managed
slice first. Starting with namespaces and ownership rules makes the later ArgoCD
handoff much clearer.

### Store raw Helm charts or repo-root manifests as the initial source of truth
Rejected.

That would make the structure less explicit and harder to reason about. The
baseline should live in a small, intentional path that makes ownership obvious.
