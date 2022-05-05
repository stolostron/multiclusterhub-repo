#!/bin/bash
# Copyright (c) 2020 Red Hat, Inc.
# Copyright Contributors to the Open Cluster Management project

CHARTS_PATH="multiclusterhub/charts"
FORMAT=sha

# temp-charts will hold the charts until it is ready to replace the current charts dir
mkdir temp-charts

while IFS=, read -r url chartpath chartname shaorbranch
do

  # Determine whether we are pulling a branch or a specific sha
  if [[ $shaorbranch == main ]] || [[ $shaorbranch == master ]] || [[ $shaorbranch == release* ]];then
    # This is a branch
    printf "Github URL: $url\tPath to chart: $chartpath\tDesired branch: $shaorbranch\n"
    FORMAT=branch
  else
    # Must be a sha
    printf "Github URL: $url\tPath to chart: $chartpath\tDesired sha: $shaorbranch\n"
    FORMAT=sha
  fi

  # currentsha is the sha of the chart currently bundled in the charts dir
  currentsha=$(grep --max-count=1 "$url" currentSHAs.csv | cut -d ',' -f4)

  # Work in temporary directory
  mkdir -p tmp
  cd tmp
  # Clone repo
  git clone $url
  # Enter repo directory
  cd */

  # Checkout branch or commit sha from origin
  git checkout $shaorbranch
  lastsha=$(git rev-parse HEAD)
  echo "Packaging at sha ${lastsha}"

  cp -R $chartpath ../../temp-charts

  cd ../..
  rm -rf tmp

  # Add to new CSV of our current shas
  echo -en "$url,$chartpath,$chartname,$lastsha\n" >> temp-currentSHAs.csv
done < desiredSHAs.csv

rm -rf ${CHARTS_PATH}
mv temp-charts ${CHARTS_PATH}
mv temp-currentSHAs.csv currentSHAs.csv
git status --porcelain
