#!/bin/bash

cd $(dirname $0)
. chart-sync.sh

git add ../multicloudhub/charts
git commit -m "[skip ci] skip travis"
git pull origin master -s recursive -X ours --../multicloudhub/charts
git push origin "HEAD:${TRAVIS_BRANCH}"

cd ..
docker build -t $1 .
