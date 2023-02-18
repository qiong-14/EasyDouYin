package apitest

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestMessageServer(t *testing.T) {
	e := newExpect(t)
	userIdA, tokenA := getTestUserToken(testUserA, e)
	userIdB, tokenB := getTestUserToken(testUserB, e)

	// testUserA send to testUserB
	for i := 0; i < 3; i++ {
		messageActionResp := e.POST("/douyin/message/action/").
			WithQuery("token", tokenA).
			WithQuery("to_user_id", userIdB).
			WithQuery("action_type", 1).
			WithQuery("content", "testUserA send to testUserB"+strconv.Itoa(i)).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		messageActionResp.Value("status_code").Number().Equal(0)
		time.Sleep(time.Second)
	}
	var preMsgTime int64
	preMsgTime = 0
	// testUserA get chat message with pre_msg_time=0
	messagechatResp := e.GET("/douyin/message/chat/").
		WithQuery("token", tokenA).
		WithQuery("to_user_id", userIdB).
		WithQuery("pre_msg_time", preMsgTime).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	messagechatResp.Value("status_code").Number().Equal(0)
	messageInfo := messagechatResp.Value("message_list").Array()
	for _, element := range messageInfo.Iter() {
		message := element.Object()
		message.Value("id").Number().Gt(0)
		message.Value("to_user_id").Number().Gt(0)
		message.Value("from_user_id").Number().Gt(0)
		message.Value("content").String().Length().Gt(0)
	}
	// testUserB get chat message with pre_msg_time=0
	messagechatResp = e.GET("/douyin/message/chat/").
		WithQuery("token", tokenB).
		WithQuery("to_user_id", userIdA).
		WithQuery("pre_msg_time", preMsgTime).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	messagechatResp.Value("status_code").Number().Equal(0)
	messageInfo = messagechatResp.Value("message_list").Array()
	for _, element := range messageInfo.Iter() {
		message := element.Object()
		message.Value("id").Number().Gt(0)
		message.Value("to_user_id").Number().Gt(0)
		message.Value("from_user_id").Number().Gt(0)
		message.Value("content").String().Length().Gt(0)
	}
	preMsgTime = 1000 * time.Now().Unix()
	time.Sleep(time.Second)

	// testUserB send to testUserA
	for i := 0; i < 3; i++ {
		messageActionResp := e.POST("/douyin/message/action/").
			WithQuery("token", tokenB).
			WithQuery("to_user_id", userIdA).
			WithQuery("action_type", 1).
			WithQuery("content", "testUserB send to testUserA"+strconv.Itoa(i)).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		messageActionResp.Value("status_code").Number().Equal(0)
		time.Sleep(time.Second)
	}

	// testUserA get chat message with pre_msg_time again
	messagechatResp = e.GET("/douyin/message/chat/").
		WithQuery("token", tokenA).
		WithQuery("to_user_id", userIdB).
		WithQuery("pre_msg_time", preMsgTime).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	messagechatResp.Value("status_code").Number().Equal(0)
	messageInfo = messagechatResp.Value("message_list").Array()
	for _, element := range messageInfo.Iter() {
		message := element.Object()
		message.Value("id").Number().Gt(0)
		message.Value("to_user_id").Number().Gt(0)
		message.Value("from_user_id").Number().Gt(0)
		message.Value("content").String().Length().Gt(0)
	}
}
