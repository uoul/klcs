package config

type OidcConfig struct {
	JwksUrl   string
	Authority string
	ClientId  string
	Roles     struct {
		SysAdmin       string `default:"ADMIN"`
		AccountManager string `default:"ACCOUNT_MANAGER"`
		NoPrint        string `default:"NO_PRINT"`
	}
}
