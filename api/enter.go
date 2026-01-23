package api

import (
	"blogx_server/api/log_api"
	"blogx_server/api/site_api"
)

type Api struct {
	SiteApi site_api.SiteApi
	LogApi  log_api.LogApi
}

var App = Api{}
