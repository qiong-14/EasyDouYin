package mw

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dgrijalva/jwt-go"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/pkg/constants"
	"net/http"
	"time"
)

type Claims struct {
	UserId   int64
	UserName string
	jwt.StandardClaims
}

// GenerateToken 通过username生成token,设置过期时间为
func GenerateToken(userid int64, username string) (string, error) {
	claims := &Claims{
		UserId:   userid,
		UserName: username, // 私有字段
		StandardClaims: jwt.StandardClaims{
			Issuer:    constants.JWTIssuer,                       // 签发人
			ExpiresAt: time.Now().Unix() + constants.JWTDuration, // 过期时间
			Subject:   constants.JWTSubject,                      // 主题
			Audience:  constants.JWTAudience,                     // 受众
			NotBefore: time.Now().Unix(),                         // 生效时间
			IssuedAt:  time.Now().Unix(),                         // 签发时间
		},
	}
	// hash256加密算法产生token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tokenString, err := token.SignedString([]byte(constants.JWTSecret)); err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(constants.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func LoginAuthentication() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		tokenStr := c.Query("token")
		if len(tokenStr) == 0 {
			tokenStr = c.PostForm("token")
		}
		if len(tokenStr) == 0 {
			c.JSON(http.StatusOK, resp.Response{
				StatusCode: 401,
				StatusMsg:  "Token doesn't exist",
			})
			c.Abort()
			return
		}
		token, err := ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusOK, resp.Response{
				StatusCode: 403,
				StatusMsg:  err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("user_id", token.UserId)
		c.Next(ctx)
	}
}
