#!/bin/bash

cd $(dirname $0)

. chart-sync.sh
git add ../multiclusterhub/charts
git commit -m "[skip ci] skip travis"
git pull origin master -s recursive -X ours --allow-unrelated-histories
git push origin "HEAD:${TRAVIS_BRANCH}"

# cd ..
docker build -t $1 .