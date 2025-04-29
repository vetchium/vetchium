{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "vetchium-apps-helm.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "vetchium-apps-helm.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "vetchium-apps-helm.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "vetchium-apps-helm.labels" -}}
helm.sh/chart: {{ include "vetchium-apps-helm.chart" . }}
{{ include "vetchium-apps-helm.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "vetchium-apps-helm.selectorLabels" -}}
app.kubernetes.io/name: {{ include "vetchium-apps-helm.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "vetchium-apps-helm.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "vetchium-apps-helm.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the granger service account to use
*/}}
{{- define "vetchium-apps-helm.grangerServiceAccountName" -}}
{{- if .Values.granger.serviceAccount.create -}}
    {{- default (printf "%s-granger-sa" (include "vetchium-apps-helm.fullname" .)) .Values.granger.serviceAccount.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
    {{- default "default" .Values.granger.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{/*
Create the name of the hermione service account to use
*/}}
{{- define "vetchium-apps-helm.hermioneServiceAccountName" -}}
{{- if .Values.hermione.serviceAccount.create -}}
    {{- default (printf "%s-hermione-sa" (include "vetchium-apps-helm.fullname" .)) .Values.hermione.serviceAccount.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
    {{- default "default" .Values.hermione.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{/*
Create the name of the sortinghat service account to use
*/}}
{{- define "vetchium-apps-helm.sortinghatServiceAccountName" -}}
{{- if .Values.sortinghat.serviceAccount.create -}}
    {{- default (printf "%s-sortinghat-sa" (include "vetchium-apps-helm.fullname" .)) .Values.sortinghat.serviceAccount.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
    {{- default "default" .Values.sortinghat.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{/*
Create the name of the dev-seed service account to use
*/}}
{{- define "vetchium-apps-helm.devSeedServiceAccountName" -}}
{{- if .Values.devSeed.serviceAccount.create -}}
    {{- default (printf "%s-dev-seed-sa" (include "vetchium-apps-helm.fullname" .)) .Values.devSeed.serviceAccount.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
    {{- default "default" .Values.devSeed.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{/*
Create the name of the sqitch service account to use
*/}}
{{- define "vetchium-apps-helm.sqitchServiceAccountName" -}}
{{- if .Values.sqitch.serviceAccount.create -}}
    {{- default (printf "%s-sqitch-sa" (include "vetchium-apps-helm.fullname" .)) .Values.sqitch.serviceAccount.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
    {{- default "default" .Values.sqitch.serviceAccount.name -}}
{{- end -}}
{{- end -}}
