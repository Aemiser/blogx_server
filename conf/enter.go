package conf

type Config struct {
	System System `yaml:"server"`
	Log    Log    `yaml:"log"`
	DB     DB     `yaml:"db"`
	DB1    DB     `yaml:"db1"`
}
