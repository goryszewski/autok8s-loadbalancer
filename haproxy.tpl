global
    # global settings here

defaults
    # defaults here

{{ range .Ports }}
#
{{ $NodePort :=  .NodePort }}
frontend f_{{ $.Namespace }}_{{ $.Name }}_{{ .Port }}
    bind 0.0.0.0:{{ .Port }}
    default_backend b_{{ $.Namespace }}_{{ $.Name }}_{{ .Port }}

backend b_{{ $.Namespace }}_{{ $.Name }}_{{ .Port }}
    {{- range  $value := $.Nodes }}
    server {{  $value.Name }} {{  $value.IP }}:{{ $NodePort }} check 
    {{- end }}
{{ end}}