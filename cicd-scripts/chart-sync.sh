CHARTS_PATH="../../../../../multicloudhub/charts"
echo "Fetching charts from csv"
while IFS=, read -r f1 f2
do
  mkdir -p tmp
  cd tmp
  git clone $f1
  cd */
  git checkout $f2
  var1=$(echo $f1 | cut -f5 -d/)  #get the repo name
  var2=$(echo $var1 | cut -f1 -d-) #get the first word (ie kui in kui-web-terminal)
  var3=$(find . -type d -name "$var2*") #look for folder in repo that starts with ^ rather than stable/*/
  cd $var3
  PACKAGE="$(helm package ./)"
  find . -name '*tgz' | xargs -J% mv % $CHARTS_PATH
  cd ../../../../
  rm -rf tmp
done < chartSHA.csv
helm repo index ../multicloudhub/charts
