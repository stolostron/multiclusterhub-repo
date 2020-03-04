#!/bin/bash

cd $(dirname $0)
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash 
. chart-sync.sh

cd ..
docker build -t $1 .
