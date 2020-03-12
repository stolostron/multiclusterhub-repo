cd $(dirname $0)
CHARTS_PATH="../../../../../multiclusterhub/charts"
CICD_FOLDER="../../../../"
echo "Fetching charts from csv"
while IFS=, read -r f1 f2
do
  mkdir -p tmp
  cd tmp
  git clone $f1
  var1=$(echo "$(ls)" | cut -f5 -d/)  #get the repo name
  cd */ 
  if [ $f2 != "latest" ]
  then
    git checkout $f2
  fi
  var2=$(echo $var1 | cut -f1 -d-) #get the first word (ie kui in kui-web-terminal)
  var3=$(find . -type d -name "$var2*") #look for folder in repo that starts with ^ rather than stable/*/
  cd $var3
  PACKAGE="$(helm package ./)"
  find . -type f -name "*tgz" | xargs -I '{}' mv '{}' $CHARTS_PATH
  cd $CICD_FOLDER
  rm -rf tmp
done < chartSHA.csv
helm repo index --url http://multiclusterhub-repo:3000/charts ../multiclusterhub/charts