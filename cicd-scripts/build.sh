#!/bin/bash

cd $(dirname $0)

. chart-sync.sh

cd ..
docker build -t $1 .

if [ "${TRAVIS_BRANCH}" != "master" ] && [ "${TRAVIS_BRANCH}" != "release-1.0.0" ]; then
    git add .
    git merge master -m "[skip ci] resolve conflicts" -s recursive -X ours
    git push origin "HEAD:${TRAVIS_BRANCH}"
fi


# git add ../multiclusterhub/charts
# git commit -m "[skip ci] skip travis"
# git pull origin master -s recursive -X ours --allow-unrelated-histories
# git push origin "HEAD:${TRAVIS_BRANCH}"
