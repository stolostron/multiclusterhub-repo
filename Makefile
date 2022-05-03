# Copyright Contributors to the Open Cluster Management project

REGISTRY ?= quay.io/stolostron
IMG ?= multiclusterhub-repo
VERSION ?= latest

build:
	go build -o start-repo ./cmd/repo

unit-test:
	go test -v ./...

image:
	docker build -t "$(REGISTRY)/$(IMG):$(VERSION)" .

push:
	docker push "$(REGISTRY)/$(IMG):$(VERSION)"

# local builds a docker image and runs it locally
local:
	docker build -t "$(IMG):$(VERSION)" . && docker run -it --rm --expose 3000 -p 3000:3000 -e CHART_VERSION=$(CHART_VERSION)  "$(IMG):$(VERSION)"

update-charts:
	bash cicd-scripts/chart-sync.sh

set-copyright:
	./cicd-scripts/set-copyright.sh
	
patch-charts-in-cluster: 
	./cicd-scripts/patchCharts.sh
