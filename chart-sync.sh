#!/bin/bash
# Copyright 2019 IBM Corporation.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

##
# This script is designed to leverage environment variables to move charts from a source HelmRepository
# To a github hosted target repository
#
# SOURCE_URL - URL to the source Helm registry
# SOURCE_FILE - List of charts with versions to be extracted
# SOURCE_USER - User to use for connecting to a secure chart repository
# SOURCE_TOKEN - Token/APIKEY/Password to connect to the chart repository
#
# TARGET_URL - URL to the Github repository that will be targeted
# TARGET_USER - User to use for connecting to a secure github repository
# TARGET_TOKEN - Token/APIKEY/Password to connect to the chart repository

# Launching from Travis, make sure that all of the above Variables are configured for the Travis job.
cd $(dirname $0)

# Load variables for delevoper mode
# if [ -f ../../env-cp4mcm.sh ]; then
#   . ../../env-cp4mcm.sh
# fi

# if [ ! -f ./$SOURCE_FILE ]; then
#   echo "Could not find chart list file: ${SOURCE_FILE}"
#   exit 1
# fi
# CHART_LIST=`cat ${SOURCE_FILE}`

# # Always refresh the Source Helm Repository
# echo "Refreshing the SOURCE-REPO connection: ${SOURCE_URL}"
# helm repo list | grep SOURCE-REPO > /dev/null 2>&1
# if [ $? -eq 0 ]; then
#   helm repo remove SOURCE-REPO
# fi
# helm repo add --username ${SOURCE_USER} --password ${SOURCE_TOKEN} SOURCE-REPO ${SOURCE_URL}

# echo "Clone the rh-ibm-synergy repository"
#  git clone "https://${TARGET_TOKEN}@${TARGET_URL}"
# if [ $? -ne 0 ]; then
#   echo "Failed to clone Target ${TRAGET_URL}"
#   exit 1
# fi

package_chart() 
{
  cd stable/*/ 
  PACKAGE="$(helm package ./)"
  find . -name '*tgz' | xargs -J% mv % ../../../
  cd ../../../../../
  echo "after" $PWD
}
echo "Fetching charts from SOURCE-REPO"
mkdir -p multicloudhub/charts
git submodule update --init --recursive
SUB_LIST="$(grep path .gitmodules | sed 's/.*= //')"
for submodule in ${SUB_LIST}; do
  echo ${submodule}
  mv ${submodule}  multicloudhub/charts
  cd "multicloudhub/charts/${submodule}"
  package_chart 
  rm -rf multicloudhub/charts/${submodule}
done
helm repo index multicloudhub/charts

# for chart in ${CHART_LIST}; do
#   # chart=(${chartVer//:/ })
#   # git clone ${chart}
#   # echo " > ${chart[0]} ${chart[1]}"
#   # helm fetch --version ${chart[1]} SOURCE-REPO/${chart[0]}
#   # helm fetch --untar ${chart} ${stable/ibm-search-prod}
#   # if [ $? -ne 0 ]; then
#   #   echo "Failed. Could not pull chart: ${chart[0]} ${chart[1]}"
#   #   exit 1
#   # fi
# done
# if [ -f ./index.yaml ]; then
#   echo "Copy original index.yaml"
#   cat ./index.yaml
#   mv index.yaml orig_index.yaml
# else
#   echo "No charts found, create a NEW index.yaml"
#   touch ./orig_index.yaml
# fi

# echo "Generating a new index"
# rawUrl="https://${TARGET_URL/github.com/raw.githubusercontent.com}/master/charts"
# echo "Raw URL: ${rawUrl}"
# helm repo index ./ --url $rawUrl
# if [ $? -ne 0 ]; then
#   echo "Could not generate Helm repository index.yaml"
#   exit 1
# fi
# echo "Comparing indexes to find NEW digests"
# diff index.yaml orig_index.yaml | grep "digest:"
# if [ $? -ne 0 ]; then
#   echo "No new content found, Done!"
#   exit 0
# fi
# echo "New digest found, committing new Helm Repository index.yaml and charts"
# rm ./orig_index.yaml
# echo "Adding NEW items"
# git add .
# echo "Commit NEW items"
# git commit -m "cp4mcm-manifest BoT, daily build"
# echo "Push NEW items"
# git push
# if [ $? -ne 0 ]; then
#   echo "Failed push changes, check the logs"
#   exit 1
# fi
