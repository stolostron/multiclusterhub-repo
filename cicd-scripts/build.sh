#!/bin/bash

. chart-sync.sh

docker build -t $1 .
