package conf

import "strings"

type Uploads struct {
	Size         int64    `yaml:"size"`
	Type         string   `yaml:"type"`
	WriteList    []string `yaml:"writeList"`
	ImageDir     string   `yaml:"imageDir"`
	ArticleCover string   `yaml:"articleCover"`
}

const (
	KB int64 = 1024
	MB int64 = 1024 * 1024
	GB int64 = 1024 * 1024 * 1024
	TB int64 = 1024 * 1024 * 1024 * 1024
)

func (u Uploads) GetType() string {
	return strings.ToUpper(u.Type)
}
func (u Uploads) GetSizeType() int64 {
	switch strings.ToUpper(u.Type) {
	case "KB":
		return KB
	case "MB":
		return MB
	case "GB":
		return GB
	case "TB":
		return TB
	default:
		return KB
	}
}
