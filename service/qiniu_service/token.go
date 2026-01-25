package qiniu_service

import (
	"blogx_server/global"
	"context"
	"time"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
)

func GenToken() (token string, err error) {
	// 创建七牛云凭证
	mac := credentials.NewCredentials(global.Config.QiNiu.AccessKey, global.Config.QiNiu.SecretKey)

	// 创建上传策略，设置有效期
	putPolicy, err := uptoken.NewPutPolicy(global.Config.QiNiu.Bucket, time.Now().Add(time.Duration(global.Config.QiNiu.Expired)*time.Second))

	if err != nil {
		return
	}

	// 生成上传令牌
	token, err = uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background()) // 修正方法名

	if err != nil {
		return
	}

	return
}
