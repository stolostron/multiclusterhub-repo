#!/bin/bash

cd $(dirname $0)

if [ $TRAVIS_BRANCH != "master" ] && [ $TRAVIS_BRANCH != "release-1.0.0" ]; then 
    . chart-sync.sh
    git add ../multicloudhub/charts
    git commit -m "[skip ci] skip travis"
    git pull origin master -s recursive -X ours --allow-unrelated-histories
    git push origin "HEAD:${TRAVIS_BRANCH}"
fi

cd ..
docker build -t $1 .