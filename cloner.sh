
git clone git@github.com:open-cluster-management/kui-web-terminal-chart.git
helm3 package kui-web-terminal-chart/stable/kui-web-terminal -d multicloudhub/charts
rm -rf kui-web-terminal-chart

git clone git@github.com:open-cluster-management/search-chart.git
helm3 package search-chart/stable/search-prod -d multicloudhub/charts
rm -rf search-chart

git clone git@github.com:open-cluster-management/application-chart.git
helm3 package application-chart/stable/application-chart -d multicloudhub/charts
rm -rf application-chart

git clone git@github.com:open-cluster-management/grc-chart.git
helm3 package grc-chart/stable/grc-chart -d multicloudhub/charts
rm -rf grc-chart

git clone git@github.com:open-cluster-management/management-ingress-chart.git
helm3 package management-ingress-chart/stable/management-ingress -d multicloudhub/charts
rm -rf management-ingress-chart

git clone git@github.com:open-cluster-management/rcm-chart.git
helm3 package rcm-chart/stable/rcm -d multicloudhub/charts
rm -rf rcm-chart

git clone git@github.com:open-cluster-management/console-chart.git
helm3 package console-chart/stable/console-chart -d multicloudhub/charts
rm -rf console-chart

# git clone git@github.com:open-cluster-management/multicloud-mongodb-chart.git
# helm3 package multicloud-mongodb-chart/stable/multicloud-mongodb -d multicloudhub/charts
# rm -rf multicloud-mongodb-chart

git clone git@github.com:open-cluster-management/cert-manager-webhook-chart.git
helm3 package cert-manager-webhook-chart/stable/cert-manager-webhook -d multicloudhub/charts
rm -rf cert-manager-webhook-chart

git clone git@github.com:open-cluster-management/cert-manager-chart.git
helm3 package cert-manager-chart/stable/cert-manager -d multicloudhub/charts
rm -rf cert-manager-chart

git clone git@github.com:open-cluster-management/configmap-watcher-chart.git
helm3 package configmap-watcher-chart/stable/configmap-watcher -d multicloudhub/charts
rm -rf configmap-watcher-chart

git clone git@github.com:open-cluster-management/topology-chart.git
helm3 package topology-chart/stable/topology-chart -d multicloudhub/charts
rm -rf topology-chart 

curl -L https://charts.bitnami.com/bitnami/nginx-5.1.6.tgz > ./multicloudhub/charts/nginx-5.1.6.tgz

helm3 repo index --url http://multicloudhub-repo:3000/charts ./multicloudhub/charts