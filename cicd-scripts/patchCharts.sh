#!/bin/bash
# build image with charts and replace image in cluster
# Use current namespace
REPOPOD=$(kubectl get pods -o=name | grep multiclusterhub-repo | sed "s/^.\{4\}//")
NAMESPACE=$(echo $(oc project) | cut -d " " -f 3)
NAMESPACE=${NAMESPACE//\"}

echo ${NAMESPACE}
echo ${REPOPOD}
echo $PWD

kubectl exec ${REPOPOD} -- sh -c 'rm -rf multiclusterhub/charts' 
kubectl cp $PWD/multiclusterhub/charts ${NAMESPACE}/${REPOPOD}:multiclusterhub/charts 

