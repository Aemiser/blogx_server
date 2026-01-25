package file

import (
	"blogx_server/utils"
	"errors"
	"strings"
)

func ImageSuffixJudgment(fileName string, list []string) (suffix string, err error) {
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
