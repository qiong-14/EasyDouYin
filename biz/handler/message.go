package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/qiong-14/EasyDouYin/biz/common"
	"github.com/qiong-14/EasyDouYin/dal"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var tempChat = map[string][]common.Message{}

var messageIdSequence = int64(1)

type ChatResponse struct {
	common.Response
	MessageList []common.Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(ctx context.Context, c *app.RequestContext) {
	toUserId := c.Query("to_user_id")
	content := c.Query("content")
	userId := c.Query("user_id")
	userIdA, _ := strconv.ParseInt(userId, 10, 64)
	userIdB, _ := strconv.ParseInt(toUserId, 10, 64)
	if user, err := dal.GetUserById(ctx, userIdA); err == nil {
		chatKey := genChatKey(user.Id, userIdB)
		atomic.AddInt64(&messageIdSequence, 1)
		curMessage := common.Message{
			Id:         messageIdSequence,
			Content:    content,
			CreateTime: time.Now().Format(time.Kitchen),
		}

		if messages, exist := tempChat[chatKey]; exist {
			tempChat[chatKey] = append(messages, curMessage)
		} else {
			tempChat[chatKey] = []common.Message{curMessage}
		}
		c.JSON(http.StatusOK, common.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// MessageChat all users have same follow list
func MessageChat(ctx context.Context, c *app.RequestContext) {
	toUserId := c.Query("to_user_id")
	userId := c.Query("user_id")

	userIdA, _ := strconv.ParseInt(userId, 10, 64)
	userIdB, _ := strconv.ParseInt(toUserId, 10, 64)
	if user, err := dal.GetUserById(ctx, userIdA); err != nil {
		chatKey := genChatKey(user.Id, userIdB)
		c.JSON(http.StatusOK, ChatResponse{Response: common.Response{StatusCode: 0}, MessageList: tempChat[chatKey]})
	} else {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
