#!/bin/bash

cd $(dirname $0)
git pull origin "HEAD:${TRAVIS_BRANCH}"
. chart-sync.sh

git add ../multicloudhub/charts
git commit -m "[skip ci] skip travis"
git push origin "HEAD:${TRAVIS_BRANCH}"

cd ..
docker build -t $1 .
