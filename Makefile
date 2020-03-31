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

USE_VENDORIZED_BUILD_HARNESS ?=

ifndef USE_VENDORIZED_BUILD_HARNESS
-include $(shell curl -s -H 'Authorization: token ${GITHUB_TOKEN}' -H 'Accept: application/vnd.github.v4.raw' -L https://api.github.com/repos/open-cluster-management/build-harness-extensions/contents/templates/Makefile.build-harness-bootstrap -o .build-harness-bootstrap; echo .build-harness-bootstrap)
else
-include vbh/.build-harness-vendorized
endif

REGISTRY ?= quay.io/rhibmcollab
IMG ?= multiclusterhub-repo
VERSION ?= latest

image:
	docker build -t "$(REGISTRY)/$(IMG):$(VERSION)" .

push:
	docker push "$(REGISTRY)/$(IMG):$(VERSION)"

unit-test:
	go test -v

update-charts: 
	./cicd-scripts/patchCharts.sh

# local builds a docker image and runs it locally
local:
	docker build -t "$(IMG):$(VERSION)" . && docker run -it --rm --expose 3000 -p 3000:3000  "$(IMG):$(VERSION)"
	