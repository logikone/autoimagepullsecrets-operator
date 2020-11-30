package webhooks

const (
	IPSInjectionBase            = "autoimagepullsecrets.io"
	IPSInjectionEnabled         = IPSInjectionBase + "/injection"
	IPSInjectionMatch           = IPSInjectionBase + "/match"
	IPSInjectionMatchNamespaced = IPSInjectionBase + "/match-namespaced"
	IPSInjectionRegistries      = IPSInjectionBase + "/registries"
)
