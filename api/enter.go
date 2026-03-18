package api

import (
	"blogx_server/api/article_api"
	"blogx_server/api/banner_api"
	"blogx_server/api/captcha_api"
	"blogx_server/api/comment_api"
	"blogx_server/api/image_api"
	"blogx_server/api/log_api"
	"blogx_server/api/site_api"
	"blogx_server/api/site_msg_api"
	"blogx_server/api/user_api"
)

type Api struct {
	SiteApi    site_api.SiteApi
	LogApi     log_api.LogApi
	Image      image_api.ImageApi
	Banner     banner_api.BannerApi
	CaptchaApi captcha_api.CaptchaApi
	UserApi    user_api.UserApi
	ArticleApi article_api.ArticleApi
	CommentApi comment_api.CommentApi
	SiteMsgAPi site_msg_api.SiteMsgApi
}

var App = Api{}
