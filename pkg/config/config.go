package config

type Config struct {
	MerchantKey string
	SecretKey   string
	Endpoint    string
	CheckSumKey string
}

func New(key, secret string, checksum string, endpoint string) *Config {
	return &Config{
		MerchantKey: key,
		SecretKey:   secret,
		Endpoint:    endpoint,
		CheckSumKey: checksum,
	}
}
