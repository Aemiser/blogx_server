package res

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	ImageUploadFailed   Code = 6001
	ImageFileTooLarge   Code = 6002
	ImageTypeNotAllowed Code = 6003
	ImageTransferFailed Code = 6004
	ImageSaveFailed     Code = 6005
	ImageNotFound       Code = 6006
	QiNiuNotEnabled     Code = 6007
	QiNiuTokenFailed    Code = 6008
)

func InitImageCode() {
	RegisterCode(ImageUploadFailed, "图片上传失败")
	RegisterCode(ImageFileTooLarge, "文件过大")
	RegisterCode(ImageTypeNotAllowed, "文件类型不支持")
	RegisterCode(ImageTransferFailed, "图片获取失败")
	RegisterCode(ImageSaveFailed, "图片保存失败")
	RegisterCode(ImageNotFound, "图片不存在")
	RegisterCode(QiNiuNotEnabled, "七牛云未启用")
	RegisterCode(QiNiuTokenFailed, "七牛云token生成失败")
}

func ImageFailWithMsgf(code Code, c *gin.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	Response{Code: code, Data: empty, Msg: msg}.Json(c)
}
