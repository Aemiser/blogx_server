package image_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/utils"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (ImageApi) ImageUploadView(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	filename := fileHeader.Filename
	confSize := global.Config.Uploads.Size              // 从配置项读取大小
	confSizeType := global.Config.Uploads.GetSizeType() // 从配置项读取类型
	if fileHeader.Size >= confSize*confSizeType {
		res.FailWithMsgf(c, "文件大于 %d %s", confSize, global.Config.Uploads.GetType())
		return
	}

	suffix, err := imageSuffixJudgment(filename, global.Config.Uploads.WriteList)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	byteData, _ := io.ReadAll(file)
	hash := utils.Md5(byteData)
	filePath := fmt.Sprintf("uploads/%s/%s.%s", global.Config.Uploads.ImageDir, hash, suffix)

	// 判断hash是否在库中
	var model models.ImageModel
	err = global.Db.Take(&model, "hash = ?", hash).Error
	if err == nil {
		// 说明存在
		logrus.Infof("上传的图片重复了 %s<==>%s", filename, hash)
		res.Success(model.Path, "图片上传成功", c)
		return
	}
	// 入库
	err = global.Db.Create(&models.ImageModel{
		Filename: fileHeader.Filename,
		Path:     filePath,
		Size:     fileHeader.Size,
		Hash:     hash,
	}).Error
	if err != nil {
		logrus.Infof("数据库创建图片失败: %v", err)
		res.FailWithError(err, c)
		return
	}
	c.SaveUploadedFile(fileHeader, filePath)
	res.Success(filePath, "图片上传成功", c)
}

func imageSuffixJudgment(fileName string, list []string) (suffix string, err error) {
	_list := strings.Split(fileName, ".")
	if len(_list) == 1 {
		err = errors.New("文件格式错误")
		return
	}
	suffix = _list[len(_list)-1]
	if !utils.InList(list, suffix) {
		err = errors.New("文件非法")
		return
	}
	return

}
