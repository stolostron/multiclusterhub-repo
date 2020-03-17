#!/bin/bash

cd $(dirname $0)
echo "TRAVIS_BRANCH: ${TRAVIS_BRANCH}"

# if we are merging into master, update the charts and 
# amend the charts to the previous commit
if [ "${TRAVIS_BRANCH}" == "chartTest" ]; then
    git checkout "${TRAVIS_BRANCH}"
    . chart-sync.sh
    git add ../multiclusterhub/charts
    git commit -m "[skip ci] update charts"
    git push origin ${TRAVIS_BRANCH}
fi

cd ..
docker build -t $1 .