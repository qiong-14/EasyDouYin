package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/mw"
	"github.com/qiong-14/EasyDouYin/utils"
	"net/http"
	"strconv"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
//var usersLoginInfo = map[string]resp.User{
//	"zhangleidouyin": {
//		Id:            1,
//		Name:          "zhanglei",
//		FollowCount:   10,
//		FollowerCount: 5,
//		IsFollow:      true,
//	},
//}

//var userIdSequence = int64(1)

type UserLoginResponse struct {
	resp.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	resp.Response
	User resp.User `json:"user"`
}

func Register(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")

	// 查找用户名是否已经注册
	u, err := dal.GetUserByName(ctx, username)
	if u.Name == username {
		c.JSON(consts.StatusOK, UserLoginResponse{
			Response: resp.Response{StatusCode: 1, StatusMsg: "user already exits"},
		})
		return
	}
	// 加密储存
	if err = dal.CreateUser(ctx, &dal.User{Name: username, Password: utils.Encoder(password)}); err != nil {
		c.JSON(consts.StatusOK, UserLoginResponse{
			Response: resp.Response{StatusCode: 1, StatusMsg: "user create failed"},
		})
		return
	}
	if u, err = dal.GetUserByName(ctx, username); err != nil {
		c.JSON(consts.StatusOK, UserLoginResponse{
			Response: resp.Response{StatusCode: 1, StatusMsg: "failed to get user id"},
		})
		return
	} else {
		token, _ := mw.GenerateToken(u.Id, u.Name)
		c.JSON(consts.StatusOK, UserLoginResponse{
			Response: resp.Response{StatusCode: 0},
			UserId:   u.Id,
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

	if user, err := dal.GetUserByName(ctx, username); err == nil {
		if utils.Encoder(password) == user.Password {
			token, _ := mw.GenerateToken(user.Id, user.Name)
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: resp.Response{StatusCode: 0},
				UserId:   user.Id,
				Token:    token,
			})
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: resp.Response{StatusCode: 1, StatusMsg: "user password error"},
			})
		}
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: resp.Response{StatusCode: 1, StatusMsg: "user doesn't exist"},
		})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

func UserInfo(ctx context.Context, c *app.RequestContext) {
	id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if user, err := dal.GetUserById(ctx, id); err == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: resp.Response{StatusCode: 0},
			User: resp.User{
				Id:   user.Id,
				Name: user.Name,
			},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: resp.Response{StatusCode: 1, StatusMsg: "user doesn't exist"},
		})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}
