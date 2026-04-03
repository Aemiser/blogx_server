package image_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/utils"
	file2 "blogx_server/utils/file"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
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

	suffix, err := file2.ImageSuffixJudgment(filename, global.Config.Uploads.WriteList)
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

	// 前缀拼接
	//filePath = "http://" + c.Request.Host + "/" + filePath
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
