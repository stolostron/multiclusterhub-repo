#!/bin/bash

cd $(dirname $0)
. chart-sync.sh

git pull origin "HEAD:${TRAVIS_BRANCH}" -s recursive -X ours
git add ../multicloudhub/charts
git commit -m "[skip ci] skip travis"
git push origin "HEAD:${TRAVIS_BRANCH}"

cd ..
docker build -t $1 .
