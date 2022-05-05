{{/*
Copyright Contributors to the Open Cluster Management project
*/}}

{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "management-ingress.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 24 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 24 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "management-ingress.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 24 | trimSuffix "-" -}}
{{- end -}}


{{/*
Define console route host.
*/}}
{{- define "management-ingress.consoleRouteHost" -}}
{{- printf "%s.%s" .Values.console_route.name .Values.cluster_basedomain | trimSuffix "-" -}}
{{- end -}}


{{/*
Define oauth callback url.
*/}}
{{- define "management-ingress.oauthCallbackUrl" -}}
{{- printf "https://%s%s" (include "management-ingress.consoleRouteHost" .) .Values.oauth_client.callback_url | trimSuffix "-" -}}
{{- end -}}
