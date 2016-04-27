Sorry, we weren’t able to certify {{.Project.Provider.Name}} and cannot offer the PeerProtect™ service guarantee on
your {{.Project.Job.Type}} job.  We only certify the very best service providers out there. But we still want to help.

{{if .Providers}}
We did find these alternative providers in your area that we have already certified.

{{range $provider := .Providers}}
{{ $provider.Name}}
{{ $provider.Address1}}
{{ $provider.City}}, {{ $provider.State}} {{ $provider.Zip}}
Ph: {{ $provider.DisplayablePhone}}
Email: {{ $provider.Email}}
{{end}}

You can contact these providers with the details of your job to gauge their interest and availability. If you
choose to work with them, come back to http://{{.Project.Server}}/#/step1/ to re-submit your details and purchase your
satisfaction guarantee.
{{else}}
If you are considering alternative providers, please come back to http://{{.Project.Server}}/#/step1/ to check
whether they are certified and eligible for the PeerProtect™ guarantee.
{{end}}





