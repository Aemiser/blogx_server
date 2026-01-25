package site

type SiteInfo struct {
	Title string `json:"title" yaml:"title"`
	Logo  string `json:"logo" yaml:"logo"`
	Beian string `json:"beian" yaml:"beian"`
	Mode  int8   `json:"mode" yaml:"mode"  binding:"oneof=1 2"` // 1 社区模式 2 博客模式
}
type Project struct {
	Title   string `json:"title" yaml:"title"`
	Icon    string `json:"icon" yaml:"icon"`
	WebPath string `json:"webPath" yaml:"webPath"`
}
type Seo struct {
	Keywords    string `json:"keywords" yaml:"keywords"`
	Description string `json:"description" yaml:"description"`
}
type About struct {
	SiteDate string `json:"siteDate" yaml:"siteDate"`
	QQ       string `json:"qq" yaml:"qq"`
	Version  string `json:"version" yaml:"-"`
	Wechat   string `json:"wechat" yaml:"wechat"`
	Gitee    string `json:"gitee" yaml:"gitee"`
	Github   string `json:"github" yaml:"github"`
	Bilibili string `json:"bilibili" yaml:"bilibili"`
}
type Login struct {
	QQLogin          bool `json:"qqLogin" yaml:"qqLogin"`
	UsernamePwdLogin bool `json:"usernamePwdLogin" yaml:"usernamePwdLogin"`
	EmailLogin       bool `json:"emailLogin" yaml:"emailLogin"`
	Captcha          bool `json:"captcha" yaml:"captcha"`
}

type ComponsetInfo struct {
	Title  string `json:"title" yaml:"title"`
	Enable bool   `json:"enable" yaml:"enable"`
}
type IndexRight struct {
	List []ComponsetInfo `json:"list" yaml:"list"`
}
type Article struct {
	NoExamine bool `json:"noExamine" yaml:"noExamine"`
}
