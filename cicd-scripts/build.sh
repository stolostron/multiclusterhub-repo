#!/bin/bash

cd $(dirname $0)
echo "TRAVIS_BRANCH: ${TRAVIS_BRANCH}"

# if we are merging into master, update the charts and 
# amend the charts to the previous commit
if [ "${TRAVIS_BRANCH}" == "chartTest" ]; then
    git checkout "${TRAVIS_BRANCH}"
    . chart-sync.sh
    git add ../multiclusterhub/charts
    git commit --amend --no-edit
    git push origin +${TRAVIS_BRANCH}:${TRAVIS_BRANCH}

    #now merge the new master branch into release-1.0.0

    git checkout chartAutomation
    git merge chartTest 
fi
# . chart-sync.sh
# git add ../multiclusterhub/charts
# git commit -m "[skip ci] skip travis"
# git pull origin master -s recursive -X ours --allow-unrelated-histories
# git push origin "HEAD:${TRAVIS_BRANCH}"

cd ..
docker build -t $1 .