package conf

type Jwt struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
	Issuer string `yaml:"issuer"`
}
