CHARTS_PATH="../../../../../multicloudhub/charts"
echo "Fetching charts from csv"
while IFS=, read -r f1 f2
do
  mkdir -p tmp
  cd tmp
  git clone $f1
  cd */
  git checkout $f2
  cd stable/*/
  PACKAGE="$(helm package ./)"
  find . -name '*tgz' | xargs -J% mv % CHARTS_PATH
  cd ../../../../
   rm -rf tmp
done < chartSHA.csv
helm repo index CHARTS_PATH
curl -L https://charts.bitnami.com/bitnami/nginx-5.1.6.tgz > multicloudhub/charts/nginx-5.1.6.tgz
