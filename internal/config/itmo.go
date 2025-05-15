package config

type ITMO struct {
	BaseURL     string `path:"base_url"`
	RedirectURI string `path:"redirect_url"`
	ClientID    string `path:"client_id"`
	ProviderURL string `path:"provider_url"`
}
