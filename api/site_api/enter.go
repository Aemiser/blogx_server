package site_api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type SiteApi struct {
}

func (SiteApi) SIteInfoView(c *gin.Context) {
	fmt.Printf("1")
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"name":    "site",
			"version": "1.0.0",
		},
	})

	return
}
