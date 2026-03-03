package image_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/utils"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ImageTransferRequest struct {
	Url string `json:"url" binding:"required"`
}

func (ImageApi) ImageTransferView(c *gin.Context) {
	cr := middlerware.GetBind[ImageTransferRequest](c)

	response, err := http.Get(cr.Url)
	if err != nil {
		logrus.Errorf("图片获取失败: %v", err)
		res.FailWithMsg("图片获取失败", c)
		return
	}
	fmt.Println(response.Status)
	suffix := "png"
	switch response.Header.Get("Content-Type") {
	case "/image/avif":
		suffix = "avif"
	}
	byteData, err := io.ReadAll(response.Body)
	hash := utils.Md5(byteData)
	filePath := fmt.Sprintf("uploads/%s/%s.%s", global.Config.Uploads.ImageDir, hash, suffix)
	err = os.WriteFile(filePath, byteData, 0666)
	if err != nil {
		logrus.Errorf("图片保存失败: %v", err)
		res.FailWithMsg("图片保存失败", c)
		return
	}

	res.SuccessWithData("/"+filePath, c)

}
