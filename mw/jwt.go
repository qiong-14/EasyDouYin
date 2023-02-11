package mw

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
	utils2 "github.com/qiong-14/EasyDouYin/utils"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	
	"net/http"
	"time"
	"strings"
	"encoding/base64"
	"encoding/json"
)

 var (
	 JwtMiddleware *jwt.HertzJWTMiddleware
	 IdentityKey   = "identity"
 )

 
 func InitJwt() {
	 var err error
	 JwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		 Realm:         "test zone",
		 Key:           []byte("secret key"),
		 Timeout:       time.Hour,
		 MaxRefresh:    time.Hour,
		 TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		 TokenHeadName: "Bearer",
		 LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			payloads := resp.Payload{}
			json_payloads, _ := base64.RawURLEncoding.DecodeString(strings.Split(token, ".")[1])
			err=json.Unmarshal(json_payloads, &payloads)
			c.JSON(http.StatusOK, resp.UserLoginResponse{
				Response: resp.Response{StatusCode: 0,StatusMsg:"login success" },
				UserId:   payloads.Identity,
				Token:    token,
			})
			if err!=nil{
				panic(err)
			}
		 },
		 Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			 var loginStruct struct {
				 Username  string `form:"username" json:"username" query:"username" vd:"(len($) > 0 && len($) < 30); msg:'Illegal format'"`
				 Password string `form:"password" json:"password" query:"password" vd:"(len($) > 0 && len($) < 30); msg:'Illegal format'"`
			 }
			 if err := c.BindAndValidate(&loginStruct); err != nil {
				 return nil, err
			 }
			 if user, err := dal.GetUserByName(ctx, loginStruct.Username); err == nil {
				if utils2.Encoder(loginStruct.Password) == user.Password{
					return user,nil
				}
			 }else{
				c.JSON(http.StatusOK, resp.UserLoginResponse{
					Response: resp.Response{StatusCode: 1,StatusMsg:"user does not exist or wrong password" },
				})
				 return nil, err
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
			 c.JSON(http.StatusOK, utils.H{
				 "code":    code,
				 "message": message,
			 })
		 },
	 })
	 if err != nil {
		 panic(err)
	 }
 }
 