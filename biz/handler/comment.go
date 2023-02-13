package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"net/http"
)

type CommentListResponse struct {
	resp.Response
	CommentList []resp.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	resp.Response
	Comment resp.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(ctx context.Context, c *app.RequestContext) {
	userId, _ := c.Get("user_id")
	userIdA, _ := userId.(int64)

	actionType := c.Query("action_type")

	if user, err := dal.GetUserById(ctx, userIdA); err == nil {
		if actionType == "1" {
			text := c.Query("comment_text")
			c.JSON(http.StatusOK, CommentActionResponse{Response: resp.Response{StatusCode: 0},
				Comment: resp.Comment{
					Id: 1,
					User: resp.User{
						Id:            user.Id,
						Name:          user.Name,
						FollowCount:   100,
						FollowerCount: 101,
						IsFollow:      true,
					},
					Content:    text,
					CreateDate: "05-01",
				}})
			return
		}
		c.JSON(http.StatusOK, resp.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

// CommentList all videos have same demo comment list
func CommentList(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    resp.Response{StatusCode: 0},
		CommentList: DemoComments,
	})
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}
