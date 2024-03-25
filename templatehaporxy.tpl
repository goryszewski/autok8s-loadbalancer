global
    # global settings here

defaults
    # defaults here

{{ range .Ports }}
#
{{ $NodePort :=  .NodePort}}
frontend f_{{ $.Namespace }}_{{ $.Name }}_{{ .Port }}
    bind {{ $.Ip }}:{{ .Port }}
    default_backend b_{{ $.Namespace }}_{{ $.Name }}_{{ .Port }}

backend b_{{ $.Namespace }}_{{ $.Name }}_{{ .Port }}
    {{- range  $value := $.Nodes }}
    server {{  $value.Name }} {{  $value.Private_ip }}:{{ $NodePort }} check 
    {{- end }}
{{ end}}
