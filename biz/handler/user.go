package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiong-14/EasyDouYin/biz/common"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/mw"
	"github.com/qiong-14/EasyDouYin/utils"
	"net/http"
	"strconv"
)

type UserRegisterResponse struct {
	common.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}
type UserLoginResponse struct {
	common.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	common.Response
	User common.User `json:"user"`
}

func Register(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")
	fmt.Println(username)
	u, err := dal.GetUserByName(ctx, username)
	fmt.Println(u, err)
	if err == nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: "user already exits"},
		})
		return
	}
	if err = dal.CreateUser(ctx, &dal.User{Name: username, Password: utils.Encoder(password)}); err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: "user create failed"},
		})
		return
	}
	if user, err := dal.GetUserByName(ctx, username); err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: "failed to get user id"},
		})
		return
	} else {
		token, _ := mw.GenerateToken(user.Id, user.Name)
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}
func Login(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")
	pwd := utils.Encoder(password)
	if user, err := dal.GetUserByName(ctx, username); err == nil {
		if pwd == user.Password {
			token, _ := mw.GenerateToken(user.Id, user.Name)
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: common.Response{StatusCode: 0},
				UserId:   user.Id,
				Token:    token,
			})
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: common.Response{
					StatusCode: 1,
					StatusMsg:  "password error",
				},
			})
		}
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: "user doesn't exist"},
		})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

func UserInfo(ctx context.Context, c *app.RequestContext) {
	id := c.Query("user_id")
	userID, _ := strconv.ParseInt(id, 10, 64)
	if user, err := dal.GetUserById(ctx, userID); err == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: common.Response{
				StatusCode: 0,
			},
			User: common.User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   1,
				FollowerCount: 2,
				IsFollow:      true,
			},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "user doesn't exist",
			},
		})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())

}
