{{/* Full name of the chart */}}
{{- define "secureshare.fullname" -}}
{{- .Chart.Name }}-{{ .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/* Common labels for all resources */}}
{{- define "secureshare.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
app.kubernetes.io/name: {{ include "secureshare.fullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/* Selector labels for all resources */}}
{{- define "secureshare.selectorLabels" -}}
app.kubernetes.io/name: {{ include "secureshare.fullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/* Release name for the chart */}}
{{- define "secureshare.releaseName" -}}
{{- .Chart.Name }}-{{ .Chart.AppVersion }}
{{- end }}
