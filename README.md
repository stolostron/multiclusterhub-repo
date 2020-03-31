multicloudhub-repo is intended to serve charts for the [multicloudhub-operator installer](https://github.com/open-cluster-management/multicloudhub-operator)

## Updating Charts in /multiclusterhub/charts
Pulls latest from master branch.

If you want to pin down to specific git commit, go to cicd-scripts/chartSHA.csv and update your respective chart by replacing `latest` with your git commit sha. 
DISCLAIMER: This will update ALL charts in our repo, unless they are pinned to specific version in cicd-scripts/chartSHA.csv

- Update the README with the current time and open a PR on a new branch. During the travis build, a second commit will be pushed to this branch that includes new charts packaged into tgz's. 

- After approval, when you go to Squash and Merge, please delete
this line `[skip ci] add charts` in the Merge Message

Open new PR pulling these changes into release-1.0.0 (might need to contact Installer Squad)

Last updated manually: March 30, 11:20 EST
