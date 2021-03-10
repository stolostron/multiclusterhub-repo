#!/bin/bash
# Copyright (c) 2020 Red Hat, Inc.
# Copyright Contributors to the Open Cluster Management project

set -e

if ! command -v helm &> /dev/null
then
    echo "helm could not be found"
    curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
    chmod 700 get_helm.sh
    ./get_helm.sh
    rm get_helm.sh
fi

cd $(dirname $0)
CHARTS_PATH=$(echo $(pwd)/../multiclusterhub/charts)
CHART_VERSION="$(cat ../CHART_VERSION)"

echo "Fetching charts from csv"
touch ../multiclusterhub/charts/text.txt
rm ../multiclusterhub/charts/* #rm all charts first, in case chart versions are changed

rm -rf tmp-chart-build
mkdir -p tmp-chart-build

workingDirectory=$(pwd)

while IFS=, read -r chartLink folderName chartBranch
do

  chartLink=https://${GITHUB_USER}:${GITHUB_TOKEN}@${chartLink}

  git clone ${chartLink} tmp-chart-build/$folderName && cd tmp-chart-build/$folderName
  git checkout $chartBranch

  chartPath="stable/$folderName"
  if [ ! -d "$chartPath" ]; then
    echo "Unable to determine path to chart for $chartPath"
    exit 1
  fi

  helm package $chartPath --version $CHART_VERSION -d $CHARTS_PATH

  cd $workingDirectory
  echo ""
done < chartSHA.csv
helm repo index --url http://multiclusterhub-repo:3000/charts $CHARTS_PATH