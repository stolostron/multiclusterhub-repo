# Copyright (c) 2020 Red Hat, Inc.

#Script to be run in travis when charts are updated in master


cd $(dirname $0)
git clone https://github.com/open-cluster-management/multicloudhub-repo.git
cd multicloudhub-repo
git checkout chartAutomation
cicd-scripts/chart-sync.sh
git add .
git commit -m "[skip ci] greetings from fakeMaster "
git push  https://github.com/open-cluster-management/multicloudhub-repo.git "chartAutomation"
