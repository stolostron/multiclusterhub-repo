#!/bin/bash

cd $(dirname $0)
. chart-sync.sh

cd ..
docker build -t $1 .
