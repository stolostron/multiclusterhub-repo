#!/bin/bash
# Copyright (c) 2020 Red Hat, Inc.

set -e

cd $(dirname $0)

#if on PR not master or release, update PR with latest charts.Otherwise just build image
if [ "${TRAVIS_BRANCH}" != "master" ] && [[ "${TRAVIS_BRANCH}" != "release-"* ]] && [[ "${TRAVIS_BRANCH}" != "dev-"* ]]; then
    git clone https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com/open-cluster-management/multiclusterhub-repo.git
    cd multiclusterhub-repo
    git checkout "${TRAVIS_BRANCH}"
    git pull origin "${TRAVIS_BRANCH}"
    lastCommitMsg=$(git log -1 --pretty=%B)

    if [[ "$lastCommitMsg" == *"[skip-chart-sync]"* ]]; then
        echo "Chart rebuild unnecessary charts due to [skip-chart-sync] commit msg."
        docker build -t $1 .
        exit 0
    fi

    cicd-scripts/chart-sync.sh
    docker build -t $1 .
    git add .
    git commit -s -m "[skip-chart-sync] add charts"
    git merge origin/release-2.1 -m "[skip-chart-sync] resolve conflicts" -s recursive -X ours
    git push origin "HEAD:${TRAVIS_BRANCH}"
else 
    cd ..
    docker build -t $1 .
fi
    
