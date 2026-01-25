package image_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (ImageApi) ImageUploadView(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	confSize := global.Config.Uploads.Size              // 从配置项读取大小
	confSizeType := global.Config.Uploads.GetSizeType() // 从配置项读取类型
	if fileHeader.Size >= confSize*confSizeType {
		res.FailWithMsgf(c, "文件大于 %d %s", confSize, global.Config.Uploads.GetType())
		return
	}

	err = imageSuffixJudgment(fileHeader.Filename, global.Config.Uploads.WriteList)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	filePath := fmt.Sprintf("/uploads/images/%s", fileHeader.Filename)
	c.SaveUploadedFile(fileHeader, filePath)
	res.Success(filePath, "上传成功", c)
}

func imageSuffixJudgment(fileName string, list []string) error {
	_list := strings.Split(fileName, ".")
	if len(_list) == 1 {
		return errors.New("文件格式错误")
	}
	suffix := _list[len(_list)-1]
	if !utils.InList(list, suffix) {
		return errors.New("文件非法")
	}
	return nil

}
