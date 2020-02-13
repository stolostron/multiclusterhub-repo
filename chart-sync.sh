package_chart() 
{
  cd stable/*/ 
  PACKAGE="$(helm package ./)"
  find . -name '*tgz' | xargs -J% mv % ../../../
  cd ../../../../../
  echo "after" $PWD
}
echo "Fetching charts from SOURCE-REPO"
mkdir -p cp4mcm-install/charts
git submodule update --init --recursive
SUB_LIST="$(grep path .gitmodules | sed 's/.*= //')"
for submodule in ${SUB_LIST}; do
  echo ${submodule}
  mv ${submodule}  cp4mcm-install/charts
  cd "cp4mcm-install/charts/${submodule}"
  package_chart 
  rm -rf cp4mcm-install/charts/${submodule}
done
helm repo index cp4mcm-install/charts