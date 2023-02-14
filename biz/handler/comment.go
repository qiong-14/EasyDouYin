package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/mw"
	"net/http"
	"time"
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
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())

	u, _ := c.Get(mw.IdentityKey)
	actionType := c.Query("action_type")
	if user, err := dal.GetUserById(ctx, u.(*dal.User).Id); err == nil {
		if actionType == "1" {
			text := c.Query("comment_text")
			c.JSON(http.StatusOK, CommentActionResponse{Response: resp.Response{StatusCode: 0},
				Comment: resp.Comment{
					Id:         1,
					User:       dal.GetRespUser(ctx, user.Id),
					Content:    text,
					CreateDate: time.Now().String(),
				}})
			return
		}
		c.JSON(http.StatusOK, resp.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

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
