{{- define "vela-worker.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "vela-worker.fullname" -}}
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

{{- define "vela-worker.labels" -}}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" }}
{{ include "vela-worker.selectorLabels" . }}
app.kubernetes.io/version: {{ default .Chart.AppVersion .Values.image.tag | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "vela-worker.selectorLabels" -}}
app.kubernetes.io/name: {{ include "vela-worker.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "vela-worker.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (include "vela-worker.fullname" .) .Values.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{- define "vela-worker.secretName" -}}
{{- default (include "vela-worker.fullname" .) .Values.server.existingSecret -}}
{{- end -}}

{{- define "vela-worker.runtimeNamespace" -}}
{{- default .Release.Namespace .Values.runtime.namespace -}}
{{- end -}}

{{- define "vela-worker.podsTemplateName" -}}
{{- if .Values.runtime.podsTemplate.name -}}
{{- .Values.runtime.podsTemplate.name -}}
{{- else if .Values.runtime.podsTemplate.create -}}
{{- printf "%s-pipeline-pods" (include "vela-worker.fullname" .) -}}
{{- end -}}
{{- end -}}

{{- define "vela-worker.isKubernetes" -}}
{{- if eq .Values.runtime.driver "kubernetes" }}true{{ end -}}
{{- end -}}
