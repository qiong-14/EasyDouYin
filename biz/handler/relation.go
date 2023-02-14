package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/mw"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type UserListResponse struct {
	resp.Response
	UserList []resp.User `json:"user_list"`
}

type FriendUser struct {
	resp.User
	message string `json:"message,omitempty"`
	msgType int64  `json:msg_type`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())

	u, _ := c.Get(mw.IdentityKey)
	id1 := u.(*dal.User).Id

	if _, err := dal.GetUserById(ctx, id1); err != nil {
		c.JSON(http.StatusOK, resp.Response{
			StatusCode: 1,
			StatusMsg:  "用户不存在",
		})
		return
	}

	id := c.Query("to_user_id")
	id2, _ := strconv.ParseInt(id, 10, 64)

	if _, err := dal.GetUserById(ctx, id2); err != nil {
		c.JSON(http.StatusOK, resp.Response{
			StatusCode: 1,
			StatusMsg:  "用户不存在",
		})
		return
	}

	ActionType, _ := strconv.Atoi(c.Query("action_type"))

	// 	关注操作： 1-关注 2-取消关注
	if ActionType == 1 {
		// 若没有记录创建记录
		if find, _ := dal.FindRelation(ctx, id1, id2); !find {
			dal.CreateRelation(ctx, id1, id2)
		}
		c.JSON(http.StatusOK, resp.Response{
			StatusCode: 0,
			StatusMsg:  "关注成功",
		})
	} else if ActionType == 2 {
		// 若有记录删除记录
		if find, _ := dal.FindRelation(ctx, id1, id2); find {
			dal.DeleteRelation(ctx, id1, id2)
		}
		c.JSON(http.StatusOK, resp.Response{
			StatusCode: 0,
			StatusMsg:  "取消关注成功",
		})
	} else {
		log.Printf("无效操作")
		c.JSON(http.StatusOK, resp.Response{
			StatusCode: 1,
			StatusMsg:  "无效的关注操作",
		})
	}
}

// FollowList 得到用户的关注列表
func FollowList(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())

	ids := c.Query("user_id")
	id, _ := strconv.ParseInt(ids, 10, 64)
	if _, err := dal.GetUserById(ctx, id); err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: resp.Response{
				StatusCode: 1,
				StatusMsg:  "用户不存在",
			},
		})
		return
	}
	userIdx, _ := dal.GetFollowList(ctx, id)
	followList := make([]resp.User, len(userIdx))
	var wg sync.WaitGroup
	wg.Add(len(userIdx))
	// 可以用goroutine获取用户的关注者信息
	for i, f := range userIdx {
		go func(j int, id int64) {
			defer wg.Done()
			follow := dal.GetRespUser(ctx, id)
			follow.IsFollow = true
			followList[j] = follow
		}(i, f)
	}
	wg.Wait()
	c.JSON(http.StatusOK, UserListResponse{
		Response: resp.Response{
			StatusCode: 0,
		},
		UserList: followList,
	})

}

// FollowerList 获取粉丝列表
func FollowerList(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())

	ids := c.Query("user_id")
	id, _ := strconv.ParseInt(ids, 10, 64)
	if _, err := dal.GetUserById(ctx, id); err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: resp.Response{
				StatusCode: 1,
				StatusMsg:  "用户不存在",
			},
		})
		return
	}
	userIdx, _ := dal.GetFansList(ctx, id)
	fansList := make([]resp.User, len(userIdx))
	var wg sync.WaitGroup
	wg.Add(len(userIdx))
	for i, f := range userIdx {
		// 采用goroutine获取粉丝信息
		go func(j int, fansId int64) {
			defer wg.Done()
			fan := dal.GetRespUser(ctx, fansId)
			// 查找该用户是否关注了粉丝 找到为真，默认为假
			if y, _ := dal.FindRelation(ctx, id, fansId); y {
				fan.IsFollow = true
			}
			fansList[j] = fan
		}(i, f)
	}
	wg.Wait()
	c.JSON(http.StatusOK, UserListResponse{
		Response: resp.Response{
			StatusCode: 0,
		},
		UserList: fansList,
	})

}

// FriendList all users have same friend list
func FriendList(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())

	ids := c.Query("user_id")
	id, _ := strconv.ParseInt(ids, 10, 64)
	if _, err := dal.GetUserById(ctx, id); err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: resp.Response{
				StatusCode: 1,
				StatusMsg:  "用户不存在",
			},
		})
		return
	}
	userIdx, _ := dal.GetFansList(ctx, id)
	fansList := make([]resp.User, len(userIdx))
	var wg sync.WaitGroup
	wg.Add(len(userIdx))
	for i, f := range userIdx {
		// 采用goroutine获取粉丝信息
		go func(j int, fansId int64) {
			defer wg.Done()
			fan := dal.GetRespUser(ctx, fansId)
			// 查找该用户是否关注了粉丝 找到为真，默认为假
			if y, _ := dal.FindRelation(ctx, id, fansId); y {
				fan.IsFollow = true
			}
			fansList[j] = fan
		}(i, f)
	}
	wg.Wait()
	c.JSON(http.StatusOK, UserListResponse{
		Response: resp.Response{
			StatusCode: 0,
		},
		UserList: fansList,
	})
}
