#!/bin/bash

cd $(dirname $0)
. chart-sync.sh

git add ../multicloudhub/charts
git commit -m "[skip ci] skip travis"
git pull origin "HEAD:${TRAVIS_BRANCH}" -s recursive -X ours
git push origin "HEAD:${TRAVIS_BRANCH}"

cd ..
docker build -t $1 .
