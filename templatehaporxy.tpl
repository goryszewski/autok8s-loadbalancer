global
    # global settings here

defaults
    # defaults here

{{ range . }}

frontend f_{{ .Namespace }}_{{ .Name }}
    # a frontend that accepts requests from clients
    bind {{ .Address }}:{{ .Port }}
    default_backend b_{{ .Namespace }}_{{ .Name }}

backend b_{{ .Namespace }}_{{ .Name }}
    # servers that fulfill the requests
    balance roundrobin
    server {{ .Name }}_{{ .AddressS }}_{{ .PortS }} {{ .AddressS }}:{{ .PortS }} 

{{ end }}