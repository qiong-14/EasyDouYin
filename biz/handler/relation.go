package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/jwt"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
)

type UserListResponse struct {
	resp.Response
	UserList []dal.UserVo `json:"user_list"`
}

// RelationAction

func RelationAction(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
	//get userid from token
	userId, err := middleware.GetUserIdRedis(jwt.GetToken(ctx, c))
	if err != nil {
		return
	}
	if _, err := dal.GetUserById(ctx, userId); err == nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	//get followId from request
	followId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		log.Println("get follow_id failed")
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "get follow_id failed"})
		return
	}
	//get actionType from request
	actionType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		log.Println("get action_type failed")
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "get action_type failed"})
		return
	}
	//update in dataset
	err = dal.UpdateOrCreateRelation(ctx, &dal.Follows{FollowerId: userId, FollowedId: followId, ActionType: actionType})
	if err != nil {
		log.Println("update follow relation failed")
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "update follow relation failed"})
		return
	}
	c.JSON(http.StatusOK, resp.Response{StatusCode: 0, StatusMsg: "successfully updated follow relation"})
}

// FollowList

func FollowList(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())

	////test if token is expired
	//_, err := middleware.GetUserIdRedis(jwt.GetToken(ctx, c))

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	log.Println(userId)

	if err != nil {
		return
	}
	if _, err := dal.GetUserById(ctx, userId); err == nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

	followUserList, err := dal.FollowUserList(ctx, userId)

	if err != nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "fail to get follow list"})
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: resp.Response{
			StatusCode: 0,
		},
		UserList: followUserList,
	})

}

// FollowerList

func FollowerList(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())

	////test if token is expired
	//_, err := middleware.GetUserIdRedis(jwt.GetToken(ctx, c))

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	if err != nil {
		return
	}
	if _, err := dal.GetUserById(ctx, userId); err == nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

	followerUserList, err := dal.FollowerUserList(ctx, userId)
	if err != nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "fail to get follower list"})
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: resp.Response{
			StatusCode: 0,
		},
		UserList: followerUserList,
	})
}

// FriendList

func FriendList(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())

	////test if token is expired
	//_, err := middleware.GetUserIdRedis(jwt.GetToken(ctx, c))

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	if err != nil {
		return
	}
	if _, err := dal.GetUserById(ctx, userId); err == nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

	friendUserList, err := dal.FriendUserList(ctx, userId)
	if err != nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "fail to get friend list"})
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: resp.Response{
			StatusCode: 0,
		},
		UserList: friendUserList,
	})
}
