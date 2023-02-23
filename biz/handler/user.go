package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/jwt"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/constants"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
	"github.com/qiong-14/EasyDouYin/service"
	"github.com/qiong-14/EasyDouYin/tools"
)

func Register(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
	username := c.Query("username")
	password := c.Query("password")
	if constants.CheckUserRegisterInfo {
		pass1, pass2 := tools.CheckUserRegisterInfo(username, password)
		info := ""
		if !pass1 {
			info += " username form incorrect "
		}
		if !pass2 {
			info += " weak password"
		}
		if len(info) > 0 {
			c.JSON(consts.StatusOK, resp.UserLoginResponse{
				Response: resp.Response{StatusCode: 1, StatusMsg: info},
			})
			return
		}
	}
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
	var err error
	// 传入user_id时直接获取，未传入时寻找鉴权的user_id
	if len(idStr) == 0 {
		id, err = middleware.GetUserIdRedis(jwt.GetToken(ctx, c))
	} else {
		id, err = strconv.ParseInt(idStr, 10, 64)
	}
	if err != nil {
		return
	}
	if user, _ := service.GetUserInfo(ctx, id); user != dal.InvalidUser {
		FollowCount, _ := dal.FollowCount(ctx, user.Id)
		FollowerCount, _ := dal.FollowerCount(ctx, user.Id)
		favoriteCount := service.GetFavVideoCount(ctx, user.Id)
		workCount, _ := dal.GetPublishListById(ctx, user.Id)
		c.JSON(http.StatusOK, resp.UserResponse{
			Response: resp.Response{StatusCode: 0},
			User: resp.User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   FollowCount,
				FollowerCount: FollowerCount,
				FavoriteCount: favoriteCount,
				WorkCount:     int64(len(workCount)),
			},
		})
	} else {
		c.JSON(http.StatusOK, resp.UserResponse{
			Response: resp.Response{StatusCode: 1, StatusMsg: "user doesn't exist"},
		})
	}

}
