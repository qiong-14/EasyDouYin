package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/jwt"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
)

type UserListResponse struct {
	resp.Response
	UserList []resp.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(ctx context.Context, c *app.RequestContext) {
	userId, err := middleware.GetUserIdRedis(jwt.GetToken(ctx, c))
	if err != nil {
		return
	}
	if _, err := dal.GetUserById(ctx, userId); err == nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

// FollowList all users have same follow list
func FollowList(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: resp.Response{
			StatusCode: 0,
		},
		UserList: []resp.User{DemoUser},
	})
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

// FollowerList all users have same follower list
func FollowerList(ctx context.Context, c *app.RequestContext) {
	u, _ := c.Get(middleware.IdentityKey)
	user_id := u.(*dal.User).Id
	list := []dal.User_info{}
	var err error
	followerIds, err := dal.GetFollowerList(ctx, user_id)
	if err != nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	} else {
		if list, err = dal.GetAllUserInfo(ctx, followerIds); err != nil {
			c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		} else {
			respList := []resp.User{}
			for _, info := range list {
				respList = append(respList, resp.User{
					Id:            info.UserId,
					Name:          info.Name,
					FollowCount:   info.FollowCount,
					FollowerCount: info.FollowerCount,
					IsFollow:      info.Isfollow,
					FavoriteCount: info.FavoriteCount,
					WorkCount:     info.WorkCount,
				})
			}
			c.JSON(http.StatusOK, UserListResponse{
				Response: resp.Response{
					StatusCode: 0,
				},
				UserList: respList,
			})
		}
	}

	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

// FriendList all users have same friend list
func FriendList(ctx context.Context, c *app.RequestContext) {
	u, _ := c.Get(middleware.IdentityKey)
	user_id := u.(*dal.User).Id
	list := []dal.User_info{}
	var err error
	friendIds, err := dal.GetFriendList(ctx, user_id)
	if err != nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	} else {
		if list, err = dal.GetAllUserInfo(ctx, friendIds); err != nil {
			c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		} else {
			respList := []resp.User{}
			for _, info := range list {
				respList = append(respList, resp.User{
					Id:            info.UserId,
					Name:          info.Name,
					FollowCount:   info.FollowCount,
					FollowerCount: info.FollowerCount,
					IsFollow:      info.Isfollow,
					FavoriteCount: info.FavoriteCount,
					WorkCount:     info.WorkCount,
				})
			}
			c.JSON(http.StatusOK, UserListResponse{
				Response: resp.Response{
					StatusCode: 0,
				},
				UserList: respList,
			})
		}
	}

	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}
