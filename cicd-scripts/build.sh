#!/bin/bash

$(cd ../ ; sh chart-sync.sh)

 docker build -t $1 .