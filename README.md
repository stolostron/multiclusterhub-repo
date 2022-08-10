[comment]: # ( Copyright Contributors to the Open Cluster Management project )

# WORK IN PROGRESS 

We are in the process of enabling this repo for community contribution. See wiki [here](https://open-cluster-management.io/concepts/architecture/).

This repo is intended to serve charts for the [multiclusterhub-operator installer](https://github.com/stolostron/multiclusterhub-operator)

## Onboarding a Component

For help onboarding your component, see the following [doc](docs/Onboarding.md). This guide will help ensure that your chart, images, CRDs, and GitHub repository fit our required specifications.

## Updating Charts

To manually update charts in this repo:
1. Update `desiredSHAs.csv` with the desired SHA or branch you wish to pull from.
2. Run `make update-charts`
3. Open a PR and ask the Installer Squad for review (via the `forum-acm-hub-installer` Slack channel)

## Updating chart in a Kubernetes cluster
Running `make patch-charts-in-cluster` will update the charts in the `multiclusterhub-repo` pod with charts in the local `multiclusterhub/charts` folder.
