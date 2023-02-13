package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/mw"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	resp.Response
	UserList []resp.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(ctx context.Context, c *app.RequestContext) {
	u, _ := c.Get(mw.IdentityKey)
	id1 := u.(*dal.User).Id
	if _, err := dal.GetUserById(ctx, id1); err != nil {
		c.JSON(http.StatusOK, resp.Response{
			StatusCode: 1,
			StatusMsg:  "用户不存在",
		})
		return
	}

	id := c.Query("to_user_id")
	id2, _ := strconv.ParseInt(id, 10, 64)

	if _, err := dal.GetUserById(ctx, id2); err != nil {
		c.JSON(http.StatusOK, resp.Response{
			StatusCode: 1,
			StatusMsg:  "用户不存在",
		})
		return
	}
	// 1-关注 2-取消关注 => 0-关注 1-取消关注
	ActionType, _ := strconv.Atoi(c.Query("action_type"))
	cancel := int8(ActionType - 1)
	if _, err := dal.FindRelation(ctx, id1, id2); err == nil {
		if err2 := dal.UpdateRelation(ctx, id1, id2, cancel); err2 == nil {
			c.JSON(http.StatusOK, resp.Response{
				StatusCode: 0,
			})
		} else {
			c.JSON(http.StatusOK, resp.Response{
				StatusCode: 1,
				StatusMsg:  "更新关注失败",
			})
		}
	} else {
		if err2 := dal.CreateFollow(ctx, id1, id2, cancel); err2 == nil {
			c.JSON(http.StatusOK, resp.Response{
				StatusCode: 0,
			})
		} else {
			c.JSON(http.StatusOK, resp.Response{
				StatusCode: 1,
				StatusMsg:  "创建关注失败",
			})
		}
	}

	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

// FollowList all users have same follow list
func FollowList(ctx context.Context, c *app.RequestContext) {

	ids := c.Query("user_id")
	id, _ := strconv.ParseInt(ids, 10, 64)
	if _, err := dal.GetUserById(ctx, id); err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: resp.Response{
				StatusCode: 1,
				StatusMsg:  "用户不存在",
			},
		})
		return
	}

	if userIdx, err := dal.GetFollowList(ctx, id); err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: resp.Response{
				StatusCode: 1,
				StatusMsg:  "获取关注列表错误",
			},
		})
		return
	} else {
		var followList []resp.User
		for _, i := range userIdx {
			follow := dal.GetRespUser(ctx, i)
			follow.IsFollow = true
			followList = append(followList, follow)
		}
		c.JSON(http.StatusOK, UserListResponse{
			Response: resp.Response{
				StatusCode: 0,
			},
			UserList: followList,
		})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

// FollowerList all users have same follower list
func FollowerList(ctx context.Context, c *app.RequestContext) {
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

// FriendList all users have same friend list
func FriendList(ctx context.Context, c *app.RequestContext) {
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
