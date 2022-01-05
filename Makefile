# GITHUB_USER containing '@' char must be escaped with '%40'
GITHUB_USER := $(shell echo $(GITHUB_USER) | sed 's/@/%40/g')

.PHONY: default
default:: init;

.PHONY: init\:
init::
ifndef DOCKER_USER
	$(info DOCKER_USER not defined)
	exit -1
endif
ifndef DOCKER_PASS
	$(info DOCKER_PASS not defined)
	exit -1
endif
ifndef GITHUB_TOKEN
	$(info GITHUB_TOKEN not defined)
	exit -1
endif

-include $(shell [ -f ".build-harness-bootstrap" ] || curl -sL -o .build-harness-bootstrap -H "Authorization: token $(GITHUB_TOKEN)" -H "Accept: application/vnd.github.v3.raw" "https://raw.github.com/stolostron/build-harness-extensions/main/templates/Makefile.build-harness-bootstrap"; echo .build-harness-bootstrap)

REGISTRY ?= quay.io/stolostron
IMG ?= multiclusterhub-repo
VERSION ?= latest

image:
	docker build -t "$(REGISTRY)/$(IMG):$(VERSION)" .

push:
	docker push "$(REGISTRY)/$(IMG):$(VERSION)"

unit-test:
	go test -v ./...

update-charts: 
	./cicd-scripts/patchCharts.sh

# local builds a docker image and runs it locally
local:
	docker build -t "$(IMG):$(VERSION)" . && docker run -it --rm --expose 3000 -p 3000:3000  "$(IMG):$(VERSION)"
	
chart-sync-local:
	./cicd-scripts/chart-sync-local.sh