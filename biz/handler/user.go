package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/mw"
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
			Response: resp.Response{StatusCode: 1, StatusMsg: "用户已经存在"},
		})
		return
	}
	// 加密储存
	if err = dal.CreateUser(ctx, &dal.User{Name: username, Password: tools.Encoder(password)}); err != nil {
		c.JSON(consts.StatusOK, resp.UserLoginResponse{
			Response: resp.Response{StatusCode: 1, StatusMsg: "用户创建失败"},
		})
		return
	}
	mw.JwtMiddleware.LoginHandler(ctx, c)

}

func UserInfo(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())

	id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if user, err := dal.GetUserById(ctx, id); err == nil {
		u := dal.GetRespUser(ctx, user.Id)
		u.IsFollow = true
		c.JSON(http.StatusOK, resp.UserResponse{
			Response: resp.Response{StatusCode: 0},
			User:     u,
		})
	} else {
		c.JSON(http.StatusOK, resp.UserResponse{
			Response: resp.Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
	}
}
