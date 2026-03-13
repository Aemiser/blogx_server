package models

import (
	"blogx_server/models/ctype"
	"blogx_server/models/enum"
	_ "embed"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ArticleModel struct {
	Model
	Title        string             `gorm:"size:32" json:"title"`
	Abstract     string             `gorm:"size:256" json:"abstract"`
	Content      string             `json:"content"`
	CategoryID   *uint              `json:"categoryID"`
	TagList      ctype.List         `gorm:"type:longtext;" json:"tagList"` // 标签列表
	Cover        string             `gorm:"size:256" json:"cover"`
	UserID       uint               `json:"userID"`
	UserModel    UserModel          `gorm:"foreignKey:UserID" json:"-"`
	LookCount    int                `json:"lookCount"`
	DiggCount    int                `json:"diggCount"`
	CommentCount int                `json:"commentCount"`
	CollectCount int                `json:"collectCount"`
	OpenComment  bool               `json:"openComment"` //开启评论
	Status       enum.ArticleStatus `json:"status"`      // 状态 草稿 审核中 已发布
}

//go:embed mappings/article_mappings_json.json
var articleMapping string

func (ArticleModel) Mapping() string {
	return articleMapping
}

func (ArticleModel) Index() string {
	return "article_index"
}

func (a *ArticleModel) BeforeDelete(tx *gorm.DB) (err error) {
	// 使用传入的事务对象 tx，确保所有操作在同一事务中
	tables := []struct {
		name  string
		count int64
	}{
		{"comment_models", 0},
		{"article_digg_models", 0},
		{"user_article_collect_models", 0},
		{"user_top_article_models", 0},
		{"user_article_look_history_models", 0},
	}

	for i := range tables {
		// 先查询数量用于日志
		if err = tx.Table(tables[i].name).Where("article_id = ?", a.ID).Count(&tables[i].count).Error; err != nil {
			return fmt.Errorf("查询 %s 数量失败：%w", tables[i].name, err)
		}

		// 直接删除，避免先 Find 再 Delete 的 N+1 问题
		if err = tx.Table(tables[i].name).Where("article_id = ?", a.ID).Delete(nil).Error; err != nil {
			return fmt.Errorf("删除 %s 失败：%w", tables[i].name, err)
		}
	}

	logrus.Infof("删除文章 ID=%d 的关联数据：\n评论 %d 条，\n点赞 %d 条，\n收藏 %d 条，\n置顶 %d 条，\n浏览 %d 条\n",
		a.ID, tables[0].count, tables[1].count, tables[2].count, tables[3].count, tables[4].count)
	return nil
}
