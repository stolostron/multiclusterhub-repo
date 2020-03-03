#!/bin/bash

cd $(dirname $0)
. chart-sync.sh
git add multicloudhub/charts
git commit -m "build"
git push origin travisPush

cd ..
docker build -t $1 .
