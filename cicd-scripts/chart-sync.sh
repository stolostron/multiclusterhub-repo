#!/bin/bash
# Copyright (c) 2020 Red Hat, Inc.
# Copyright Contributors to the Open Cluster Management project

CHARTS_PATH="multiclusterhub/charts"
CHART_VERSION="$(cat CHART_VERSION)"
FORMAT=sha
echo "---Setting chart version as ${CHART_VERSION}---"

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
  # filename is the desired chart package name
  filename="${chartname}-${CHART_VERSION}.tgz"

  # if [ $FORMAT == sha ]; then
  #   # Check if chart is using correct sha and chart version
  #   if [ "$currentsha" == "$shaorbranch" ] && [ -f "${CHARTS_PATH}/${filename}" ]; then
  #     echo $"Current sha matches desired sha and chart file exists. Copying chart over."
  #     cp "${CHARTS_PATH}/${filename}" "temp-charts/${filename}"
      
  #     # Add to new CSV of our current shas
  #     echo -en "$url,$chartpath,$chartname,$shaorbranch\n" >> temp-currentSHAs.csv
  #     continue
  #   fi
  # fi

  # if [ $FORMAT == branch ]; then
  #   ## Find the most recent sha in the repository branch

  #   # Get the repo name without the leading 'https://github.com/'
  #   httpsTrimmedURL=${url#*//}
  #   githubTrimmedURL=${httpsTrimmedURL#*/}
  #   lastsha=$(curl -s -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/${githubTrimmedURL}/git/refs/heads/${shaorbranch} | jq -r '.object.sha')
  #   echo "Last sha in the ${shaorbranch} branch of ${githubTrimmedURL} is ${lastsha}"
    
  #   # Check if chart is using the latest sha and correct chart version
  #   if [ "$currentsha" == "$lastsha" ] && [ -f "${CHARTS_PATH}/${filename}" ]; then
  #     echo "Latest SHA matches SHA in branch ${shaorbranch} and file chart version exists. Copying chart over."
  #     cp "${CHARTS_PATH}/${filename}" "temp-charts/${filename}"
      
  #     # Add to new CSV of our current shas
  #     echo -en "$url,$chartpath,$chartname,$lastsha\n" >> temp-currentSHAs.csv
  #     continue
  #   fi
  # fi

  echo "$url either does not have desired sha, doesn't have latest sha, or isn't set to the current chart version. It will need to be repackaged."

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
  # helm package $chartpath --version="${CHART_VERSION}" --destination="../../temp-charts"

  cd ../..
  rm -rf tmp

  # Add to new CSV of our current shas
  echo -en "$url,$chartpath,$chartname,$lastsha\n" >> temp-currentSHAs.csv
done < desiredSHAs.csv

rm -rf ${CHARTS_PATH}
mv temp-charts ${CHARTS_PATH}
mv temp-currentSHAs.csv currentSHAs.csv
git status --porcelain
