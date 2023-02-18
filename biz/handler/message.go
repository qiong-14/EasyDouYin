package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
)

type ChatResponse struct {
	resp.Response
	MessageList []resp.Message `json:"message_list"`
}

// 发送消息
func MessageAction(ctx context.Context, c *app.RequestContext) {
	u, _ := c.Get(middleware.IdentityKey)
	toUserId := c.Query("to_user_id")
	userIdA := u.(*dal.User).Id
	userIdB, _ := strconv.ParseInt(toUserId, 10, 64)
	actionType64, _ := strconv.ParseInt(c.Query("action_type"), 10, 32)
	actionType := int32(actionType64)
	content := c.Query("content")
	if actionType != 1 {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "action_type err"})
	}
	if user, err := dal.GetUserById(ctx, userIdA); err != nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	} else {
		curMessage := dal.Message{
			ToUserId:   userIdB,
			FromUserId: user.Id,
			Content:    content,
			CreateTime: 1000 * time.Now().Unix(),
		}
		if err = dal.InsertMessage(ctx, &curMessage); err != nil {
			return
		}
		c.JSON(http.StatusOK, resp.Response{StatusCode: 0})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

// 获取聊天记录
func MessageChat(ctx context.Context, c *app.RequestContext) {
	u, _ := c.Get(middleware.IdentityKey)
	toUserId := c.Query("to_user_id")
	userIdA := u.(*dal.User).Id
	userIdB, _ := strconv.ParseInt(toUserId, 10, 64)
	preMsgTime, _ := strconv.ParseInt(c.Query("pre_msg_time"), 10, 64)
	if user, err := dal.GetUserById(ctx, userIdA); err != nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	} else {
		// 根据双方ID和上次最新消息的时间获取最新的聊天记录
		messages, err := dal.GetMessage(ctx, userIdB, user.Id, preMsgTime)
		if err != nil {
			c.JSON(http.StatusOK, ChatResponse{Response: resp.Response{StatusCode: 0}, MessageList: []resp.Message{}})
			return
		}
		messagesResp := []resp.Message{}
		for _, msg := range messages {
			messagesResp = append(messagesResp, resp.Message{
				Id:         msg.Id,
				ToUserId:   msg.ToUserId,
				FromUserId: msg.FromUserId,
				Content:    msg.Content,
				CreateTime: msg.CreateTime,
			})
		}
		c.JSON(http.StatusOK, ChatResponse{Response: resp.Response{StatusCode: 0}, MessageList: messagesResp})
	}
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}
