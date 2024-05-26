# 3rd-Party Licenses

{{ range . }}
## {{ .Name }} ({{ .Version }})

[{{ .LicenseName }}]({{ .LicenseURL }})

```raw
{{ .LicenseText }}
```
{{ end }}
