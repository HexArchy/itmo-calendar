package config

type ITMO struct {
	BaseURL            string `yaml:"base_url"`
	RedirectURI        string `yaml:"redirect_url"`
	ClientID           string `yaml:"client_id"`
	ProviderURL        string `yaml:"provider_url"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
}
