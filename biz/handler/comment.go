package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/qiong-14/EasyDouYin/biz/common"
	"github.com/qiong-14/EasyDouYin/dal"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	common.Response
	CommentList []common.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	common.Response
	Comment common.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(ctx context.Context, c *app.RequestContext) {
	userid := c.Query("user_id")
	actionType := c.Query("action_type")
	userId, _ := strconv.ParseInt(userid, 10, 64)
	if user, err := dal.GetUserById(ctx, userId); err == nil {
		if actionType == "1" {
			text := c.Query("comment_text")
			c.JSON(http.StatusOK, CommentActionResponse{Response: common.Response{StatusCode: 0},
				Comment: common.Comment{
					Id: 1,
					User: common.User{
						Name:          user.Name,
						FollowCount:   2,
						FollowerCount: 3,
					},
					Content:    text,
					CreateDate: "05-01",
				}})
			return
		}
		c.JSON(http.StatusOK, common.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    common.Response{StatusCode: 0},
		CommentList: DemoComments,
	})
}
