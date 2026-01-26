package utils

import (
	"blogx_server/global"
	"crypto/md5"
	"encoding/hex"
	"image/color"

	"github.com/mojocn/base64Captcha"
)

func InList[T comparable](list []T, value T) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func Md5(data []byte) string {
	md5New := md5.New()
	md5New.Write(data)
	return hex.EncodeToString(md5New.Sum(nil))
}

func GetCaptcha() (string, string, error) {
	var driver base64Captcha.Driver
	var driverString base64Captcha.DriverString
	captchaConfig := base64Captcha.DriverString{
		Height:          60,
		Width:           200,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 125,
		},
	}
	driverString = captchaConfig
	driver = driverString.ConvertFonts()
	captcha := base64Captcha.NewCaptcha(driver, global.Stores)
	lid, lb64s, _, lerr := captcha.Generate()
	return lid, lb64s, lerr
}
