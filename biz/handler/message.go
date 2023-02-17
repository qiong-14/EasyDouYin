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
	"time"
)

var tempChat map[string]int = make(map[string]int)

type ChatResponse struct {
	resp.Response
	MessageList []dal.Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(ctx context.Context, c *app.RequestContext) {
	u, _ := c.Get(middleware.IdentityKey)
	toUserId := c.Query("to_user_id")
	userIdA := u.(*dal.User).Id
	userIdB, _ := strconv.ParseInt(toUserId, 10, 64)
	content := c.Query("content")

	if user, err := dal.GetUserById(ctx, userIdA); err != nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	} else {
		curMessage := dal.Message{
			ToUserId:    userIdB,
			FromUserId:  user.Id,
			Content:     content,
			CreateTime:  time.Now().Unix(),
		}
		if err=dal.InsertMessage(ctx,&curMessage);err !=nil{
			return
		}
		c.JSON(http.StatusOK, resp.Response{StatusCode: 0})
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

	if user, err := dal.GetUserById(ctx, userIdA); err != nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	} else {
		messages,err:=dal.GetMessage(ctx,userIdB,user.Id)
		if err!=nil{
			c.JSON(http.StatusOK, ChatResponse{Response: resp.Response{StatusCode: 0}, MessageList:[]dal.Message{}})
			return
		}
		chatKey := genChatKey(user.Id, userIdB)
		if _,exist:=tempChat[chatKey];exist==false{
			tempChat[chatKey]=len(messages)
			c.JSON(http.StatusOK, ChatResponse{Response: resp.Response{StatusCode: 0}, MessageList:messages})
		}else if tempChat[chatKey]==len(messages){
			c.JSON(http.StatusOK, ChatResponse{Response: resp.Response{StatusCode: 0}, MessageList:[]dal.Message{}})
		}else{
			addedMessages:=[]dal.Message{}
			for i := tempChat[chatKey]; i < len(messages); i++ {
				if messages[i].FromUserId==userIdB{
					addedMessages=append(addedMessages,messages[i])
				}
			}
			c.JSON(http.StatusOK, ChatResponse{Response: resp.Response{StatusCode: 0}, MessageList:addedMessages})
			tempChat[chatKey]=len(messages)
		}	
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
