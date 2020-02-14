#!/bin/bash
echo "BUILD GOES HERE!"

echo "<repo>/<component>:<tag> : $1"

/bin/bash ./chart-sync.sh

docker build -t $1 .