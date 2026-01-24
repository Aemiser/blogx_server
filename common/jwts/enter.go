package jwts

import (
	"blogx_server/global"
	"blogx_server/models/enum"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const TokenExpireDuration = time.Hour * 2

// const TokenExpireDuration = time.Second * 60

var Secret = []byte("人生路漫漫")

type Claims struct {
	UserID   uint          `json:"userID"`
	UserName string        `json:"userName"`
	Role     enum.RoleType `json:"role"`
}
type MyClaims struct {
	Claims Claims
	jwt.StandardClaims
}

// get token
func GetToken(Claims Claims) (string, error) {
	cla := MyClaims{
		Claims: Claims,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(global.Config.Jwt.Expire) * time.Hour).Unix(), // 过期时间
			Issuer:    global.Config.Jwt.Issuer,                                                   // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)
	return token.SignedString([]byte(global.Config.Jwt.Secret)) // 进行签名生成对应的token
}

// parse token
func ParseToken(tokenString string) (*MyClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token为空,请先登入")
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.Jwt.Secret), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, errors.New("token已过期")
		}
		if strings.Contains(err.Error(), "signature is invalid") {
			return nil, errors.New("token无效")
		}
		if strings.Contains(err.Error(), " token contains an invalid number of segments") {
			return nil, errors.New("token非法")
		}

		fmt.Println(err, 1)
		return nil, err

	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func ParseTokenByGin(c *gin.Context) (*MyClaims, error) {
	tokenString := c.GetHeader("token")
	if tokenString == "" {
		tokenString = c.Query("token")
	}
	return ParseToken(tokenString)
}
