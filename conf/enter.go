package conf

type Config struct {
	System  System  `yaml:"server"`
	Jwt     Jwt     `yaml:"jwt"`
	Log     Log     `yaml:"log"`
	Redis   Redis   `yaml:"redis"`
	DB      DB      `yaml:"db"`
	DB1     DB      `yaml:"db1"`
	Site    Site    `yaml:"site"`
	Email   Email   `yaml:"email"`
	QQ      QQ      `yaml:"qq"`
	QiNiu   QiNiu   `yaml:"qiniu"`
	Ai      Ai      `yaml:"ai"`
	Uploads Uploads `yaml:"uploads"`
}
