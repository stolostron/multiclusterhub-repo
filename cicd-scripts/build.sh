#!/bin/bash
# Copyright (c) 2020 Red Hat, Inc.

cd $(dirname $0)

#if on PR not master or release, update PR with latest charts.Otherwise just build image
if [ "${TRAVIS_BRANCH}" != "master" ] && [[ "${TRAVIS_BRANCH}" != "release-"* ]] && [[ "${TRAVIS_BRANCH}" != "dev-"* ]]; then
    git clone https://github.com/open-cluster-management/multicloudhub-repo.git
    cd multicloudhub-repo
    git checkout "${TRAVIS_BRANCH}"
    cicd-scripts/chart-sync.sh
    docker build -t $1 .
    git add .
    git commit -m "[skip ci] add charts"
    git merge master -m "[skip ci] resolve conflicts" -s recursive -X ours
    git push origin "HEAD:${TRAVIS_BRANCH}"
else 
    cd ..
    docker build -t $1 .
fi
    
