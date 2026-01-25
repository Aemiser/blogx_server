package conf

type QiNiu struct {
	Enable    bool   `yaml:"enable" json:"enable"`       //	是否启用七牛存储
	AccessKey string `yaml:"accessKey" json:"accessKey"` //	七牛 Access Key
	SecretKey string `yaml:"secretKey" json:"secretKey"` //	七牛 Secret Key
	Bucket    string `yaml:"bucket" json:"bucket"`       //	存储桶名称
	Uri       string `yaml:"uri" json:"uri"`             //	七牛存储访问域名（URL前缀）
	Region    string `yaml:"region" json:"region"`       //	存储区域
	Prefix    string `yaml:"prefix" json:"prefix"`       //	存储路径前缀（可选）
	Size      int    `yaml:"size" json:"size"`           //	文件大小限制（单位：字节）
	Expired   int    `yaml:"expired" json:"expired"`
}
