DOCKER_SCRATCH_REGISTRY := hyc-cloud-private-scratch-docker-local.artifactory.swg-devops.com

.PHONY: default
default:: init;

.PHONY: init\:
init:
ifndef DOCKER_USER
	$(info DOCKER_USER not defined)
	exit -1
endif
ifndef DOCKER_PASS
	$(info DOCKER_PASS not defined)
	exit -1
endif

pull-secret:
	kubectl create secret docker-registry scratch \
	--docker-server=$(DOCKER_SCRATCH_REGISTRY) \
	--docker-username=$(DOCKER_USER) \
	--docker-password=$(DOCKER_PASS) \
	--docker-email=$(DOCKER_USER)

docker-login:
	docker login hyc-cloud-private-scratch-docker-local.artifactory.swg-devops.com/rhacm-installer -u $(DOCKER_USER) -p $(DOCKER_PASS)

docker-build:
	docker build -t rhacm-repo:latest .
	
docker-release:
	docker tag rhacm-repo:latest hyc-cloud-private-scratch-docker-local.artifactory.swg-devops.com/rhacm-installer/rhacm-repo:0.1.0
	docker push hyc-cloud-private-scratch-docker-local.artifactory.swg-devops.com/rhacm-installer/rhacm-repo:0.1.0

docker-run:
	docker run --rm rhacm-repo:latest

docker-run-it:
	docker run --rm -it --entrypoint=/bin/bash rhacm-repo:latest