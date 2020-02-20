# Bootstrap (pull) the build harness

# GITHUB_USER containing '@' char must be escaped with '%40'
GITHUB_USER := $(shell echo $(GITHUB_USER) | sed 's/@/%40/g')
GITHUB_TOKEN ?=

DOCKER_USER := $(shell echo $(DOCKER_USER))
DOCKER_PASS := $(shell echo $(DOCKER_PASS))

-include $(shell curl -H 'Authorization: token ${GITHUB_TOKEN}' -H 'Accept: application/vnd.github.v4.raw' -L https://api.github.com/repos/open-cluster-management/build-harness-extensions/contents/templates/Makefile.build-harness-bootstrap -o .build-harness-bootstrap; echo .build-harness-bootstrap)

DOCKER_SCRATCH_REGISTRY := hyc-cloud-private-scratch-docker-local.artifactory.swg-devops.com

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

image:
	./cicd-scripts/build.sh "$(REGISTRY)/$(IMG):$(VERSION)"

push:
	./common/scripts/push.sh "$(REGISTRY)/$(IMG):$(VERSION)"

unit-test:
	./common/scripts/unit-test.sh "$(REGISTRY)/$(IMG):$(VERSION)"







pull-secret:
	kubectl create secret docker-registry scratch \
	--docker-server=$(DOCKER_SCRATCH_REGISTRY) \
	--docker-username=$(DOCKER_USER) \
	--docker-password=$(DOCKER_PASS) \
	--docker-email=$(DOCKER_USER)

docker-login:
	docker login hyc-cloud-private-scratch-docker-local.artifactory.swg-devops.com/rhacm-installer -u $(DOCKER_USER) -p $(DOCKER_PASS)

docker-build:
	docker build -t multicloudhub-repo:latest .
	
docker-release:
	docker tag rhacm-repo:latest hyc-cloud-private-scratch-docker-local.artifactory.swg-devops.com/rhacm-installer/rhacm-repo:0.1.0
	docker push hyc-cloud-private-scratch-docker-local.artifactory.swg-devops.com/rhacm-installer/rhacm-repo:0.1.0

docker-run:
	docker run --rm rhacm-repo:latest

docker-run-it:
	docker run --rm -it --entrypoint=/bin/bash rhacm-repo:latest