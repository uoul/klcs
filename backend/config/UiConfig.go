package config

type UiConfig struct {
	Mobile struct {
		DefaultPayment      string `default:"CASH"`
		DescriptionRequired bool   `default:"false"`
	}
}
