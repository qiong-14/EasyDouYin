package middleware

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/jwt"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/tools"

	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

var (
	JwtMiddleware *jwt.HertzJWTMiddleware
	IdentityKey   = "identity"
)

func InitJwt() {
	var err error
	JwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:         "apitest zone",
		Key:           []byte("secret key"),
		Timeout:       24 * time.Hour,
		MaxRefresh:    24 * time.Hour,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt, form: token",
		TokenHeadName: "Bearer",
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			payloads := resp.Payload{}
			jsonPayloads, _ := base64.RawURLEncoding.DecodeString(strings.Split(token, ".")[1])
			err = json.Unmarshal(jsonPayloads, &payloads)
			if err != nil {
				panic(err)
			}
			c.JSON(http.StatusOK, resp.UserLoginResponse{
				Response: resp.Response{StatusCode: 0, StatusMsg: "login success"},
				UserId:   payloads.Identity,
				Token:    token,
			})
			// UserId存入redis
			if err = SetUserIdRedis(payloads.Identity, token); err != nil {
				panic(err)
			}
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginStruct struct {
				Username string `form:"username" json:"username" query:"username" vd:"(len($) > 0 && len($) < 30); msg:'Illegal format'"`
				Password string `form:"password" json:"password" query:"password" vd:"(len($) > 0 && len($) < 30); msg:'Illegal format'"`
			}
			if err := c.BindAndValidate(&loginStruct); err != nil {
				return nil, err
			}
			user, err := dal.GetUserByName(ctx, loginStruct.Username)
			if err != nil {
				return nil, err
			}
			if tools.Encoder(loginStruct.Password) == user.Password {
				return user, nil
			} else {
				err = errors.New("password error")
			}
			return nil, err
		},
		IdentityKey: IdentityKey,
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return &dal.User{
				Id: int64(claims[IdentityKey].(float64)),
			}
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*dal.User); ok {
				return jwt.MapClaims{
					IdentityKey: v.Id,
				}
			}
			return jwt.MapClaims{}
		},
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			hlog.CtxErrorf(ctx, "jwt biz err = %+v", e.Error())
			return e.Error()
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(http.StatusOK, resp.UserLoginResponse{
				Response: resp.Response{StatusCode: int32(code), StatusMsg: message},
			})
		},
	})
	if err != nil {
		panic(err)
	}
}
