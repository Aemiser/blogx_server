package conf

import (
	"blogx_server/service/river_service/rule"
)

type River struct {
	ServerID uint32        `yaml:"server_id"`
	Flavor   string        `yaml:"flavor"`
	DataDir  string        `yaml:"data_dir"`
	Sources  []RiverSource `yaml:"source"`

	Rules    []*rule.Rule `yaml:"rule"`
	BulkSize int          `yaml:"bulk_size"`
	Enable   bool         `yaml:"enable"`
}

type RiverSource struct {
	Schema string   `yaml:"schema"`
	Tables []string `yaml:"tables"`
}
