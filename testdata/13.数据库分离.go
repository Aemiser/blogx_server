package main

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// 用户模型
type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
}

var DB *gorm.DB

func initDB() {
	// 主库（写）
	masterDSN := "root:123456@tcp(192.168.175.129:13306)/test?charset=utf8mb4&parseTime=True"
	// 从库（读）
	replicaDSN := "root:123456@tcp(192.168.175.129:13307)/test?charset=utf8mb4&parseTime=True"

	db, err := gorm.Open(mysql.Open(masterDSN), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("连接主库失败: %v", err)
	}

	// 注册读写分离插件
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(masterDSN)},  // 写
		Replicas: []gorm.Dialector{mysql.Open(replicaDSN)}, // 读
		Policy:   dbresolver.RandomPolicy{},
	}))
	if err != nil {
		logrus.Fatalf("数据库读写分离配置失败: %v", err)
	}

	DB = db
}
func TestReadWriteSplit(t *testing.T) {
	initDB()

	// 自动迁移表（会走主库）
	err := DB.AutoMigrate(&User{})
	if err != nil {
		t.Fatalf("迁移表失败: %v", err)
	}

	// === 测试 1: 写操作 → 应走主库 ===
	user := User{Name: "Alice"}
	result := DB.Create(&user)
	if result.Error != nil {
		t.Fatalf("插入失败: %v", result.Error)
	}
	logrus.Infof("✅ 写入成功，ID=%d", user.ID)

	// === 测试 2: 读操作 → 应走从库 ===
	var found User
	err = DB.First(&found, user.ID).Error
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	logrus.Infof("✅ 读取成功: %+v", found)

	// === 测试 3: 事务中读 → 强制走主库（保证一致性）===
	tx := DB.Begin()
	var inTx User
	err = tx.First(&inTx, user.ID).Error
	if err != nil {
		t.Fatalf("事务中查询失败: %v", err)
	}
	tx.Rollback()
	logrus.Infof("✅ 事务中读取成功（应走主库）: %+v", inTx)

	// === 测试 4: 手动强制读主库 ===
	var forceMaster User
	err = DB.Clauses(dbresolver.Write).First(&forceMaster, user.ID).Error
	if err != nil {
		t.Fatalf("强制读主库失败: %v", err)
	}
	logrus.Infof("✅ 强制读主库成功: %+v", forceMaster)

	// === 测试 5: Raw SELECT → 走从库 ===
	var count int64
	err = DB.Raw("SELECT COUNT(*) FROM users WHERE id = ?", user.ID).Scan(&count).Error
	if err != nil {
		t.Fatalf("Raw 查询失败: %v", err)
	}
	logrus.Infof("✅ Raw SELECT 成功（应走从库），count=%d", count)

	// === 测试 6: Raw UPDATE → 走主库 ===
	err = DB.Exec("UPDATE users SET name = ? WHERE id = ?", "Alice_Updated", user.ID).Error
	if err != nil {
		t.Fatalf("Raw 更新失败: %v", err)
	}
	logrus.Infof("✅ Raw UPDATE 成功（应走主库）")
}
