package flags

import "flag"

type Options struct {
	File    string
	DB      bool
	Version bool
}

var FlagOptions = new(Options)

func Parse() {
	flag.StringVar(&FlagOptions.File, "f", "settings.yaml", "配置文件")
	flag.BoolVar(&FlagOptions.DB, "d", false, "迁移数据库")
	flag.BoolVar(&FlagOptions.Version, "v", false, "查看版本")
	flag.Parse()
}
