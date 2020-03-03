
git clone git@github.com:open-cluster-management/kui-web-terminal-chart.git
helm package kui-web-terminal-chart/stable/kui-web-terminal -d multicloudhub/charts
rm -rf kui-web-terminal-chart

git clone git@github.com:open-cluster-management/search-chart.git
helm package search-chart/stable/search-prod -d multicloudhub/charts
rm -rf search-chart

git clone git@github.com:open-cluster-management/application-chart.git
helm package application-chart/stable/application-chart -d multicloudhub/charts
rm -rf application-chart

git clone git@github.com:open-cluster-management/grc-chart.git
helm package grc-chart/stable/grc -d multicloudhub/charts
rm -rf grc-chart

git clone git@github.com:open-cluster-management/management-ingress-chart.git
helm package management-ingress-chart/stable/management-ingress -d multicloudhub/charts
rm -rf management-ingress-chart

git clone git@github.com:open-cluster-management/rcm-chart.git
helm package rcm-chart/stable/rcm -d multicloudhub/charts
rm -rf rcm-chart

git clone git@github.com:open-cluster-management/console-chart.git
helm package console-chart/stable/console-chart -d multicloudhub/charts
rm -rf console-chart

# git clone git@github.com:open-cluster-management/multicloud-mongodb-chart.git
# helm package multicloud-mongodb-chart/stable/multicloud-mongodb -d multicloudhub/charts
# rm -rf multicloud-mongodb-chart

git clone git@github.com:open-cluster-management/cert-manager-webhook-chart.git
helm package cert-manager-webhook-chart/stable/cert-manager-webhook -d multicloudhub/charts
rm -rf cert-manager-webhook-chart

git clone git@github.com:open-cluster-management/cert-manager-chart.git
helm package cert-manager-chart/stable/cert-manager -d multicloudhub/charts
rm -rf cert-manager-chart

git clone git@github.com:open-cluster-management/configmap-watcher-chart.git
helm package configmap-watcher-chart/stable/configmap-watcher -d multicloudhub/charts
rm -rf configmap-watcher-chart

git clone git@github.com:open-cluster-management/topology-chart.git
helm package topology-chart/stable/topology -d multicloudhub/charts
rm -rf topology-chart 

helm repo index --url http://multicloudhub-repo:3000/charts ./multicloudhub/charts