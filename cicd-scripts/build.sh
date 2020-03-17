#!/bin/bash
git checkout master
git checkout "${TRAVIS_BRANCH}"

cd $(dirname $0)
. chart-sync.sh

cd ..
docker build -t $1 .

# if we are merging into master, update the charts and 
# amend the charts to the previous commit
if [ "${TRAVIS_BRANCH}" != "master" ] && [ "${TRAVIS_BRANCH}" != "release-1.0.0" ]; then
    git add .
    git commit -m "[skip ci] update charts"
    git rebase master -s recursive -X theirs
    git rebase --continue
    git pull
    git push origin "${TRAVIS_BRANCH}"
fi
