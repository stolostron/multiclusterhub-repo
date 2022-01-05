#!/bin/bash
# Copyright (c) 2020 Red Hat, Inc.
# Copyright Contributors to the Open Cluster Management project

# RUN THIS AS MAKE COMMAND, NOT DIRECTLY

if ! command -v helm &> /dev/null
then
    echo "helm could not be found"
    curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
    chmod 700 get_helm.sh
    ./get_helm.sh
    rm get_helm.sh
fi

GIT_REPO_BASE_DIR=$PWD
CHARTS_PATH="$GIT_REPO_BASE_DIR/multiclusterhub/charts"
CICD_SCRIPTS_PATH="$GIT_REPO_BASE_DIR/cicd-scripts"
CHART_SHA_PATH="$CICD_SCRIPTS_PATH/chartSHA.csv"
CHART_VERSION="$(cat $GIT_REPO_BASE_DIR/CHART_VERSION)"
TMP_PATH="$GIT_REPO_BASE_DIR/tmp"
echo "Fetching charts from csv"
rm $CHARTS_PATH/* #rm all charts first, in case chart versions are changed

while IFS=, read -r f1 f2
do
  mkdir -p $TMP_PATH
  cd $TMP_PATH
  git clone "https://${GITHUB_TOKEN}@github.com/stolostron/$(echo $f1 | cut -f2 -d/)"
  REPO_NAME=$(echo "$(ls)" | cut -f5 -d/)  #get the repo name
  cd $REPO_NAME
  git checkout $f2
  FIRST_WORD_IN_REPO_NAME=$(echo $REPO_NAME | cut -f1 -d-) #get the first word (ie kui in kui-web-terminal)
  INNER_CHART_DIR=$(find . -type d -name "$FIRST_WORD_IN_REPO_NAME*") #look for folder in repo that starts with ^ rather than stable/*/
  cd $INNER_CHART_DIR
  helm package ./ --version $CHART_VERSION
  find . -type f -name "*tgz" | xargs -I '{}' mv '{}' $CHARTS_PATH
  cd $GIT_REPO_BASE_DIR
  rm -rf $TMP_PATH
done < $CHART_SHA_PATH
helm repo index --url http://multiclusterhub-repo:3000/charts $CHARTS_PATH