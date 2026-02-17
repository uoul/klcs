package config

type DatabaseConfig struct {
	Host     string `default:"localhost"`
	Port     uint16 `default:"5432"`
	User     string
	Password string
	Name     string `default:"postgres"`
	SslMode  string
}
