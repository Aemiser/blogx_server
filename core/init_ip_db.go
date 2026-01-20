package core

import (
	utilsIp "blogx_server/utils/ip"
	"fmt"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/sirupsen/logrus"
)

var seacher *xdb.Searcher

func InItIPDB() {
	var dbPath = "init/ip2region.xdb"
	_searcher, err := xdb.NewWithFileOnly(dbPath)
	if err != nil {
		logrus.Fatalf("failed to create searcher: %s\n", err.Error())
		return
	}

	seacher = _searcher
}

func GetIpAddr(ip string) (addr string) {
	if utilsIp.HasLocalIPAddr(ip) {
		return "内网IP"
	}
	region, err := seacher.SearchByStr(ip)
	fmt.Printf("region: %s\n", region)
	if err != nil {
		logrus.Warnf("错误IP: %s", ip)
		return "错误IP"
	}
	addrList := strings.Split(region, "|")
	if len(addrList) != 5 {
		logrus.Warnf("异常IP: %s", ip)
		return "未知IP"
	}

	country := addrList[0]
	province := addrList[2]
	city := addrList[3]

	if province != "0" && city != "0" {
		return fmt.Sprintf("%s-%s", province, city)
	}
	if country != "0" && province != "0" {
		return fmt.Sprintf("%s-%s", country, province)
	}
	if country != "0" {
		return country
	}
	return region
}
