package utils

import (
	"blogx_server/global"
	"crypto/md5"
	"encoding/hex"
	"image/color"

	"github.com/google/uuid"
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

func GetUUID() string {
	return uuid.New().String()
}

// GetRandomStringInDigital 获取随机位数的数字
func GetRandomInDigital(length int) string {
	return base64Captcha.RandText(length, "0123456789")
}

func GetRandomWord(length int) string {
	return base64Captcha.RandText(length, "0123456789qwertyuiopasdfghjklzxcvbnm")
}

// 切片去重升级版 泛型参数 利用map的key不能重复的特性+append函数  一次for循环搞定
func Unique[T comparable](ss []T) []T {
	size := len(ss)
	if size == 0 {
		return []T{}
	}
	newSlices := make([]T, 0) //这里新建一个切片,大于为0, 因为我们不知道有几个非重复数据,后面都使用append来动态增加并扩容
	m1 := make(map[T]byte)
	for _, v := range ss {
		if _, ok := m1[v]; !ok { //如果数据不在map中,放入
			m1[v] = 1                        // 保存到map中,用于下次判断
			newSlices = append(newSlices, v) // 将数据放入新的切片中
		}
	}
	return newSlices
}
