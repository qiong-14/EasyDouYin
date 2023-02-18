package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
	"github.com/qiong-14/EasyDouYin/service"
	"github.com/qiong-14/EasyDouYin/tools"
	"net/http"
	"strconv"
)

func Register(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())

	username := c.Query("username")
	password := c.Query("password")

	// 查找用户名是否已经注册
	u, err := dal.GetUserByName(ctx, username)
	if u.Name == username {
		c.JSON(consts.StatusOK, resp.UserLoginResponse{
			Response: resp.Response{StatusCode: 1, StatusMsg: "user already exits"},
		})
		return
	}
	// 加密储存
	if err = dal.CreateUser(ctx, &dal.User{Name: username, Password: tools.Encoder(password)}); err != nil {
		c.JSON(consts.StatusOK, resp.UserLoginResponse{
			Response: resp.Response{StatusCode: 1, StatusMsg: "user create failed"},
		})
		return
	}
	middleware.JwtMiddleware.LoginHandler(ctx, c)

}

func UserInfo(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
	idStr := c.Query("user_id")

	var id int64
	// 传入user_id时直接获取，未传入时寻找鉴权的user_id
	if len(idStr) == 0 {
		user, _ := c.Get(middleware.IdentityKey)
		id = user.(*dal.User).Id
	} else {
		id, _ = strconv.ParseInt(idStr, 10, 64)
	}

	if user, _ := service.GetUserInfo(ctx, id); user != dal.InvalidUser {
		favoriteCount := service.GetFavVideoCount(ctx, user.Id)
		c.JSON(http.StatusOK, resp.UserResponse{
			Response: resp.Response{StatusCode: 0},
			User: resp.User{
				Id:            user.Id,
				Name:          user.Name,
				FavoriteCount: favoriteCount,
			},
		})
	} else {
		c.JSON(http.StatusOK, resp.UserResponse{
			Response: resp.Response{StatusCode: 1, StatusMsg: "user doesn't exist"},
		})
	}

}
