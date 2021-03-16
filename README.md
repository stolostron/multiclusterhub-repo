[comment]: # ( Copyright Contributors to the Open Cluster Management project )

# WORK IN PROGRESS

We are in the process of enabling this repo for community contribution. See wiki [here](https://open-cluster-management.io/concepts/architecture/).

This repo is intended to serve charts for the [multiclusterhub-operator installer](https://github.com/open-cluster-management/multicloudhub-operator)

## Updating Charts

To incorporate new chart changes into this repo:
1. Update the bottom of this README with the current time and open a PR on a new branch. New chart changes will be automatically pulled during the resulting Travis build.
2. Contact the Installer Squad (via the `forum-acm-hub-installer` Slack channel) to approve and merge the PR
3. **IMPORTANT**: When merging, remove any line with `[skip ci]` in it. This will allow Travis merge jobs to run.

NOTE:
- This will update ALL charts in our repo, unless they are pinned to specific version in `cicd-scripts/chartSHA.csv`
- We default pull latest commit from master from the repos, if you want to use a specific commit update `latest` in `cicd-scripts/chartSHA.csv` to the git commit sha. 

## Updating chart in cluster
Running `make update-charts` will update the charts in the `multiclusterhub-repo` pod with charts in the local `multiclusterhub/charts` folder

Last updated manually: March 16, 2021
