package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/utils/hash"
	"context"
	"fmt"
	"io"

	file2 "blogx_server/utils/file"
	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)

func SendFile(file string) (url string, err error) {
	// 创建七牛云的认证凭证
	mac := credentials.NewCredentials(global.Config.QiNiu.AccessKey, global.Config.QiNiu.SecretKey)

	// 计算文件的MD5哈希值
	hashString, err := hash.FileMd5(file)
	if err != nil {
		return "", err
	}

	// 判断文件后缀（假设是图片）
	suffix, _ := file2.ImageSuffixJudgment(file, global.Config.Uploads.WriteList)
	fileName := fmt.Sprintf("%s%s", hashString, suffix) // 修正了格式字符串
	key := fmt.Sprintf("%s/%s", global.Config.QiNiu.Prefix, fileName)

	// 创建上传管理器
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})

	// 上传文件
	err = uploadManager.UploadFile(context.Background(), file, &uploader.ObjectOptions{
		BucketName: global.Config.QiNiu.Bucket,
		ObjectName: &key, // 修正了参数
		FileName:   fileName,
	}, nil)

	if err != nil {
		return "", err
	}

	// 返回完整的文件访问URL
	return fmt.Sprintf("%s/%s", global.Config.QiNiu.Uri, key), nil
}

func SendReader(reader io.Reader) (url string, err error) {
	// 创建七牛云的认证凭证
	mac := credentials.NewCredentials(global.Config.QiNiu.AccessKey, global.Config.QiNiu.SecretKey)

	// 生成唯一的文件名
	uid := uuid.New().String()
	fileName := fmt.Sprintf("%s.png", uid) // 修正了格式字符串
	key := fmt.Sprintf("%s/%s", global.Config.QiNiu.Prefix, fileName)

	// 创建上传管理器
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})

	// 从Reader上传数据
	err = uploadManager.UploadReader(context.Background(), reader, &uploader.ObjectOptions{
		BucketName: global.Config.QiNiu.Bucket,
		ObjectName: &key, // 修正了参数
		FileName:   fileName,
	}, nil)

	if err != nil {
		return "", err
	}

	// 返回完整的文件访问URL
	return fmt.Sprintf("%s/%s", global.Config.QiNiu.Uri, key), nil
}
func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	url, err := SendFile("uploads/images/7e65869e71afed5408ad44125c3641b7.jpg")
	fmt.Println(url, err)
}
