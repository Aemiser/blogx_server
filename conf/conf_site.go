package conf

import (
	"blogx_server/conf/site"
)

type Site struct {
	SiteInfo   site.SiteInfo   `json:"siteInfo" yaml:"siteInfo"`
	Project    site.Project    `json:"project" yaml:"project"`
	Seo        site.Seo        `json:"seo" yaml:"seo"`
	About      site.About      `json:"about" yaml:"about"`
	Login      site.Login      `json:"login" yaml:"login"`
	IndexRight site.IndexRight `json:"indexRight" yaml:"indexRight"`
	Article    site.Article    `json:"article" yaml:"article"`
}
