package core

import (
	"blogx_server/flags"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Server struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
}
type Config struct {
	Server `yaml:"server"`
}

func ReadConf() {
	byteData, err := os.ReadFile(flags.FlagOptions.File)
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(byteData, &config)
	if err != nil {
		panic(fmt.Sprintf("yaml配置文件格式错误 %s", err))
	}
	fmt.Println(config)
}
