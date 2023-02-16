package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var tempChat = map[string][]resp.Message{}

var messageIdSequence = int64(1)

type ChatResponse struct {
	resp.Response
	MessageList []resp.Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(ctx context.Context, c *app.RequestContext) {
	u, _ := c.Get(middleware.IdentityKey)
	toUserId := c.Query("to_user_id")
	userIdA := u.(*dal.User).Id
	userIdB, _ := strconv.ParseInt(toUserId, 10, 64)
	content := c.Query("content")

	if user, err := dal.GetUserById(ctx, userIdA); err == nil {
		chatKey := genChatKey(user.Id, userIdB)
		atomic.AddInt64(&messageIdSequence, 1)
		curMessage := resp.Message{
			Id:         messageIdSequence,
			Content:    content,
			CreateTime: time.Now().Format(time.Kitchen),
		}

		if messages, exist := tempChat[chatKey]; exist {
			tempChat[chatKey] = append(messages, curMessage)
		} else {
			tempChat[chatKey] = []resp.Message{curMessage}
		}
		c.JSON(http.StatusOK, resp.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

// MessageChat all users have same follow list
func MessageChat(ctx context.Context, c *app.RequestContext) {
	u, _ := c.Get(middleware.IdentityKey)
	toUserId := c.Query("to_user_id")
	userIdA := u.(*dal.User).Id
	userIdB, _ := strconv.ParseInt(toUserId, 10, 64)

	if user, err := dal.GetUserById(ctx, userIdA); err == nil {
		chatKey := genChatKey(user.Id, userIdB)
		c.JSON(http.StatusOK, ChatResponse{Response: resp.Response{StatusCode: 0}, MessageList: tempChat[chatKey]})
	} else {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
