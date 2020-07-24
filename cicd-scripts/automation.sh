# Copyright (c) 2020 Red Hat, Inc.

cd $(dirname $0)

#Pass repo url as arg
update_chart() {
  echo $1
  git clone $1
  cd */
  helm package ./stable/kui-web-terminal/
  find . -type f -name "*tgz" | xargs -I '{}' mv '{}' $CHARTS_PATH
  helm repo index --url http://multiclusterhub-repo:3000/charts ../multiclusterhub/charts
  git add .
  git commit -m "[skip ci] auto-push"
}

update_chart fakekui fakebranch