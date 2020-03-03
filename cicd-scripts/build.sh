#!/bin/bash

cd $(dirname $0)
. chart-sync.sh
git add multicloudhub/charts
git status

cd ..
docker build -t $1 .
