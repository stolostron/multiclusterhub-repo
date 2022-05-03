# Management Ingress Chart

## Overview
This management-ingress-chart installs Management Ingress in your cluster.

## Chart Details
This chart:
* Installs Management Ingress chart.
* Adds a daemonset to running Management Ingress pod on master nodes.
* Adds a service management-ingress.
* Adds a configmap to configure ssl-ciphers.
* Adds a Certificate.

## Developing for management ingress chart
### Prerequisites
* Installing a PodDisruptionBudget
* Kubernetes v1.11+
* OpenShift 3.11+
* Tiller v2.12+

### Resources Required
* cpu: 200m
* memory: 256Mi

### Installing the Chart
```
helm install --namespace kube-system --name management-ingress management-ingress
```

### Configuration
The following table lists the configurable parameters of the `Management Ingress` chart and their default values.

| Parameter                              | Description                                                    | Default                       |
|----------------------------------------|----------------------------------------------------------------|-------------------------------|
| org                                    | organization name                                              | open-cluster-management       |
| pullSecret                             | the name of imagePullSecrets defined in the deployment         | null                          |
| global.pullPolicy                      | the imagePullPolicy defined in the deployment                  | Always                        |
| global.imageOverrides.oauth_proxy      | the image name of oauth_proxy defined in the deployment        | ""                            |
| global.imageOverrides.management_ingress | the image name of management_ingress defined in the deployment | ""                          |
| resources.requests.cpu                 | cpu request to run the deployment                              | 200m                          |
| resources.requests.memory              | memory request to run the deployment                           | 256Mi                         |
| config.disable-access-log              | Management Ingress configmap setting of disable-access-log     | true                          |
| config.ssl-ciphers                     | Management Ingress configmap setting of ssl-ciphers            | ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256|
| tolerations.key                        | tolerations key of the deployment                              | dedicated                     |
| tolerations.operator                   | tolerations operator of the deployment                         | Exists                        |
| tolerations.effect                     | tolerations effect of the deployment                           | NoSchedule                    |
| host_headers_check_enabled             | enable/disable host header check                               | false                         |
| allowed_host_headers                   | allowed host headers list, seperated by space                  | 127.0.0.1 localhost           |
| httpPort                               | Listening http port of nginx backend                           | 8080                          |
| httpsPort                              | Listening https port of nginx backend                          | 8443                          |
| hostPort                               | exposing container port as host port                           | true                          |
| hostNetwork                            | using host network                                             | false                         |
| fips_enabled                           | enable/disable fips mode                                       | false                         |
| enable_impersonation                   | enable/disable impersonation setting                           | false                         |
| apiserver_secure_port                  | kubernetes apiserver secure port                               | 8001                          |
| service_name                           | Management Ingress service name                                | management-ingress            |
| oauth_proxy.httpsPort                  | the https port of oauth_proxy container                        | 9443                          |
| oauth_client.id                        | the name of the oauthclient, also the client id                | multicloudingress             |
| oauth_client.secret                    | the unique secret associated with oauthclient, you can choose to use the default value or customize it, should not be empty | multicloudingresssecret       |
| oauth_client.callback_url              | the valid redirection URIs associated with the oauthclient     | /oauth/callback               |
| cluster_basedomain                     | the cluster domain which will be used to generate route's host | ""                            |
| console_route.name                     | the prefix of route's host, also the route resource name       | multicloud-console            |
| hubconfig.nodeSelector                 | the nodeSelector of the deployment                             | null                          |
| hubconfig.replicaCount                 | the replicas of the deployment                                 | 1                             |
| affinity                               | the affinity of the deployment                                 | podAntiAffinity               |

## Contributing
* See [CONTRIBUTING.md](CONTRIBUTING.md) for information about the workflow that we expect, and instructions on the developer certificate of origin that we require.
