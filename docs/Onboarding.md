# Onboarding 

This document is intended to be a guide to enable developers who would like to add a chart component to multiclusterhub-repo. Charts deployed via multiclusterhub-repo are expected meet a certain standard. This includes being properly validated and being able to accept certain values and overrides from the multiclusterhub-controller. To properly onboard a component, developers should ensure their associated images, CRDs, and charts are properly integrated.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Prereqs](#prereqs)
  - [Onboarding Images](#onboarding-images)
  - [Onboarding CRDs](#onboarding-crds)
- [Onboarding a Chart](#onboarding-a-chart)
  - [Contributing.md](#contributingmd)
  - [Service Account](#service-account)
  - [Referencing an image in a chart](#referencing-an-image-in-a-chart)
  - [Replica Count](#replica-count)
  - [Affinity](#affinity)
  - [Tolerations](#tolerations)
  - [Node Selector](#node-selector)
  - [Naming ClusterRoles and ClusterRoleBindings](#naming-clusterroles-and-clusterrolebindings)
  - [Security Policies](#security-policies)
  - [Other](#other)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Prereqs

### Onboarding Images

For help onboarding images into the pipeline, please see the following [doc](https://github.com/open-cluster-management/cicd-docs/blob/main/prow/ONBOARDING.md).

### Onboarding CRDs

We do not allow for CRDs to be installed or managed via helm. All CRDs must be removed from the charts before onboarding.

To enable CRDs to be installed, upgraded, and uninstalled properly, ensure your CRD(s) are added to the [hub-crds](https://github.com/open-cluster-management/hub-crds) repository. 

See [Contributing.md](https://github.com/open-cluster-management/hub-crds/blob/main/CONTRIBUTING.md) for help contributing to this repository.

## Onboarding a Chart

In order for a chart to be onboarded properly, ensure your chart meets the specifications below. 

The chart repository must be public and must be in github.com/open-cluster-management organization.

The chart must be valid and helm linted.

### Contributing.md

Before beginning to onboard your chart, please see our [contributind.md](https://github.com/open-cluster-management/multicloudhub-repo/blob/main/CONTRIBUTING.md). Follow the steps here for to ensure that the repo owners are aware of the desired change and can handle the desired change in a standard fashion.

### Service Account

A chart must not use the default service account. Instead if must be sure to create its own service account.

### Referencing an image in a chart

Images must be added under `global.imageOverrides` in the values.yaml file of a chart. Each image referenced in the chart, must have an accompanying image key. A chart must also accept `global.imageOverrides.pullPolicy` value. This should default to a preset of `Always`.

```yaml
# Source: nginx/values.yaml
global:
  imageOverrides:
    ## Image Key/Value pairs
    nginx: "repository.to/nginx-image:latest"
    pullPolicy: Always
```

Images can be referenced in a deployment like so -

```yaml
# Source: nginx/templates/deployment.yaml
containers:
  - name: nginx
    image: {{ .Values.global.imageOverrides.nginx }}
    imagePullPolicy: {{ .Values.global.pullPolicy }}
```

### Replica Count

Unless a replicaCount of 1 is always desirable. It is necessary to add `hubconfig.replicaCount` to the values.yaml file. This will allow the MCH CR to toggle between basic and high availability installation modes.

```yaml
# Source: nginx/values.yaml

hubconfig:
  replicaCount: 1
```

replicaCount can be referenced in a deployment like so -

```yaml
# Source: nginx/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ template "nginx.fullname" . }}-console-v2
spec:
  replicas: {{ .Values.hubconfig.replicaCount }}
...
```

### Affinity

Affinity must be set in a charts deployment(s) like so, Ensure the `ocm-antiaffinity-selector` labels is set.  -

```yaml
# Source: nginx/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  template:
    metadata:
      labels:
        ocm-antiaffinity-selector: <ANTIAFFINITY-SELECTOR>  # Add this label
        ...
    spec:
      ...
      affinity:
        nodeAffinity:
          ...
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 70
            podAffinityTerm:
              topologyKey: topology.kubernetes.io/zone
              labelSelector:
                matchExpressions:
                - key: ocm-antiaffinity-selector
                  operator: In
                  values:
                  - <ANTIAFFINITY-SELECTOR>  # Add this label
                - key: component
                  operator: In
                  values:
                  - <ANTIAFFINITY-SELECTOR>  # Add this label
          - weight: 35
            podAffinityTerm:
              topologyKey: kubernetes.io/hostname
              labelSelector:
                matchExpressions:
                - key: ocm-antiaffinity-selector
                  operator: In
                  values:
                  - <ANTIAFFINITY-SELECTOR>  # Add this label
                - key: component
                  operator: In
                  values:
                  - <ANTIAFFINITY-SELECTOR>  # Add this label
            ...
```

### Tolerations

Tolerations  must be set in a charts deployment(s) like so - 

```yaml
# Source: nginx/templates/deployment.yaml
tolerations:
  - key: dedicated
    operator: Exists
    effect: NoSchedule
  - effect: NoSchedule
    key: node-role.kubernetes.io/infra
    operator: Exists
```

### Node Selector

In the values.yaml file of a chart, there must be a `hubconfig.nodeSelector` key. This should be given a value of `null` to start. Overrides to nodeSelector are applied from the MCH CR and are passed down to each chart through this key.

```yaml
# Source: nginx/values.yaml
org: open-cluster-management

hubconfig:
  nodeSelector: null
```

```yaml
# Source: nginx/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  ...
  template:
    ...
    spec:
      ...
      containers:
        ...
      {{- with .Values.hubconfig.nodeSelector }}
      nodeSelector:
      {{ toYaml . | indent 8 }}
      {{- end }}
```

### Naming ClusterRoles and ClusterRoleBindings

Clusterroles and clusterrolebindings installed as part of open-cluster-management require a prefix specifying their hierarchy/ownership. This should be formatted  in a way that is `<org-name>:<release-name>:<clusterrole/clusterrolebinding-name>`. A full clusterrole name should resemble the following - `open-cluster-management:nginx-chart-72fa7:clusterrole`.

```yaml
# Source: nginx/values.yaml
org: open-cluster-management
```

See [clusterrole.yaml](nginx/templates/clusterrole.yaml) and [clusterrolebinding.yaml](nginx/templates/clusterrolebinding.yaml)

### Security Policies

The following security policies must be specified in the deployments, unless an exemption is approved.

```yaml
# Source: nginx/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  ...
  template:
    ...
    spec:
      hostNetwork: false
      hostPID: false
      hostIPC: false
      securityContext:
        runAsNonRoot: true
      containers:
        - name: nginx
          ...
          securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              drop:
              - ALL
            privileged: false
            readOnlyRootFilesystem: true
          ...
```

### Other

In some specific cases, a chart may require an override to be set in the Multiclusterhub CR to enable it to be passed down to the chart. In these cases, please open an issue describing the desired capability against the installer team for feedback. 