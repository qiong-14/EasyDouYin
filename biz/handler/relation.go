package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/qiong-14/EasyDouYin/biz/common"
	"github.com/qiong-14/EasyDouYin/dal"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	common.Response
	UserList []common.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(ctx context.Context, c *app.RequestContext) {
	id := c.Query("user_id")
	userId, _ := strconv.ParseInt(id, 10, 64)
	if _, err := dal.GetUserById(ctx, userId); err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: common.Response{
			StatusCode: 0,
		},
		UserList: []common.User{DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: common.Response{
			StatusCode: 0,
		},
		UserList: []common.User{DemoUser},
	})
}

// FriendList all users have same friend list
func FriendList(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: common.Response{
			StatusCode: 0,
		},
		UserList: []common.User{DemoUser},
	})
}
