#!/bin/bash

cd $(dirname $0)
. chart-sync.sh

docker build -t $1 .
