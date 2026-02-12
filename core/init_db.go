package core

import (
	"blogx_server/global"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

func InitDB() *gorm.DB {
	dc := global.Config.DB   //读
	dc1 := global.Config.DB1 //写

	//TODO pgsql的支持

	db, err := gorm.Open(mysql.Open(dc.GetDSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, //不生成外键约束
	})
	if err != nil {
		logrus.Fatalf("数据库连接失败: %v", err)
	}
	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(global.Config.DB.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(global.Config.DB.MaxOpenConns)
	// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	logrus.Infof("数据库连接成功")

	if !dc.Empty() {
		fmt.Println("读写库存在")
		// 读写库存在，则配置读写分离
		db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(dc1.GetDSN())}, // 写
			Replicas: []gorm.Dialector{mysql.Open(dc.GetDSN())},  // 读
			Policy:   dbresolver.RandomPolicy{},
		}))
	}
	return db
}
