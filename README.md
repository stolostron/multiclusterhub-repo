multicloudhub-repo is intended to serve charts for the [multicloudhub-operator installer](https://github.com/open-cluster-management/multicloudhub-operator)

## Updating Charts in /multiclusterhub/charts 
DISCLAIMER: This will update ALL charts in our repo, unless they are pinned to specific version in cicd-scripts/chartSHA.csv

If you make changes to any charts we consume and want them built into the next snapshot: 
- We default pull latest commit from master from the repos, if you want to use a specific commit update `latest` in cicd-scripts/chartSHA.csv to the git commit sha. 

- Update the README with the current time and open a PR on a new branch. During the travis build, a second commit will be pushed to this branch that includes updated charts packaged into tgz's. 

- After approval, when you go to Squash and Merge, in the Merge Message, remove any line with `[skip ci]` in it

Contact Installer Squad to pull these into release branch

## Updating chart in cluster
Running `make update-charts` will update the charts in the `multiclusterhub-repo` pod with charts in the local `multiclusterhub/charts` folder

## Updating charts in the `multiclusterhub/charts` folder (as of 3/10/20)


*NOTE: This requires the Helm 3 CLI*

Calling the cloner script will build all the charts from the master branch and generate an index
```console
bash cloner.sh
```

Last updated manually: June 3, 2020 1PM

