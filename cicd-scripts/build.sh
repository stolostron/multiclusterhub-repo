#!/bin/bash

cd $(dirname $0)
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash 
. chart-sync.sh
git add .
git commit -m "[skip ci] skip travis"
git push origin HEAD:travisPush

cd ..
docker build -t $1 .
