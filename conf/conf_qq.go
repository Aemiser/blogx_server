package conf

import "fmt"

type QQ struct {
	AppID    string `json:"appID" yaml:"appID"`
	AppKey   string `json:"appKey" yaml:"appKey"`
	Redirect string `json:"redirect" yaml:"redirect"`
}

func (q QQ) Url() string {
	return fmt.Sprintf("https://graph.qq.com/oauth2.0/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=blogx",
		q.AppID,
		q.Redirect,
	)
}
