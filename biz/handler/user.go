package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")
	u, err := dal.GetUserByName(ctx, username)
	if u.Name == username {
		c.JSON(consts.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "user already exits"},
		})
		return
	}
	if err = dal.CreateUser(ctx, &dal.User{Name: username, Password: password}); err != nil {
		c.JSON(consts.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "user create failed"},
		})
		return
	}
	if u, err = dal.GetUserByName(ctx, username); err != nil {
		c.JSON(consts.StatusOK, UserListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "failed to get user id"},
		})
		return
	} else {
		c.JSON(consts.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   u.ID,
		})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
			c.Response.StatusCode(),
			 c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

func Login(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
			c.Response.StatusCode(),
			 c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

func UserInfo(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
