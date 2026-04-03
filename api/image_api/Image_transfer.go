package image_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

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
		logrus.Errorf("图片获取失败：%v", err)
		res.FailWithMsg("图片获取失败", c)
		return
	}
	defer response.Body.Close()

	fmt.Println(response.Status)

	// 从 URL 路径中提取文件名
	filename := extractFilenameFromUrl(cr.Url)

	// 从 Content-Disposition 头获取文件名（如果有）
	if contentDisposition := response.Header.Get("Content-Disposition"); contentDisposition != "" {
		if parsedFilename := parseFilenameFromDisposition(contentDisposition); parsedFilename != "" {
			filename = parsedFilename
		}
	}

	// 根据 Content-Type 确定文件后缀
	suffix := getFileSuffix(response.Header.Get("Content-Type"))

	byteData, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("读取图片数据失败：%v", err)
		res.FailWithMsg("读取图片数据失败", c)
		return
	}

	hash := utils.Md5(byteData)
	filePath := fmt.Sprintf("uploads/%s/%s.%s", global.Config.Uploads.ImageDir, hash, suffix)
	err = os.WriteFile(filePath, byteData, 0666)
	if err != nil {
		logrus.Errorf("图片保存失败：%v", err)
		res.FailWithMsg("图片保存失败", c)
		return
	}

	// 检查库中是否存在相同的图片
	err = global.Db.Find(&models.ImageModel{}, "hash = ?", hash).Error
	if err != nil {
		// 入库
		err = global.Db.Create(&models.ImageModel{
			Filename: filename,
			Path:     filePath,
			Size:     int64(len(byteData)),
			Hash:     hash,
		}).Error
		if err != nil {
			logrus.Infof("数据库创建图片失败：%v", err)
			res.FailWithError(err, c)
			return
		}
	}

	res.SuccessWithData("/"+filePath, c)
}

// extractFilenameFromUrl 从 URL 中提取文件名
func extractFilenameFromUrl(urlStr string) string {
	// 获取 URL 路径部分
	lastSlash := strings.LastIndex(urlStr, "/")
	if lastSlash == -1 {
		return "unknown_image"
	}

	filename := urlStr[lastSlash+1:]

	// 如果文件名包含查询参数，去掉查询参数
	if queryIndex := strings.Index(filename, "?"); queryIndex != -1 {
		filename = filename[:queryIndex]
	}

	// 如果文件名为空或只有扩展名，返回默认值
	if filename == "" || !strings.Contains(filename, ".") {
		return "unknown_image"
	}

	return filename
}

// parseFilenameFromDisposition 从 Content-Disposition 头中解析文件名
func parseFilenameFromDisposition(disposition string) string {
	// 匹配 filename="xxx" 或 filename*=utf-8''xxx
	re := regexp.MustCompile(`filename\*?=["']?([^;"'\n]*)`)
	matches := re.FindStringSubmatch(disposition)
	if len(matches) > 1 && matches[1] != "" {
		filename := matches[1]
		// 如果是 RFC 5987 编码（filename*=），需要解码
		if strings.Contains(disposition, "filename*=") {
			// 简单处理 utf-8'' 前缀
			if idx := strings.Index(filename, "''"); idx != -1 {
				filename = filename[idx+2:]
			}
		}
		return filename
	}
	return ""
}

// getFileSuffix 根据 Content-Type 获取文件后缀
func getFileSuffix(contentType string) string {
	suffix := "png" // 默认后缀
	switch contentType {
	case "image/jpeg", "image/jpg":
		suffix = "jpg"
	case "image/png":
		suffix = "png"
	case "image/gif":
		suffix = "gif"
	case "image/webp":
		suffix = "webp"
	case "image/avif":
		suffix = "avif"
	case "image/svg+xml":
		suffix = "svg"
	case "image/bmp":
		suffix = "bmp"
	}
	return suffix
}
