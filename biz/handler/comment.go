package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/jwt"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
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
	// get user id
	userId, err := middleware.GetUserIdRedis(jwt.GetToken(ctx, c))
	if err != nil {
		return
	}
	u, err := dal.GetUserById(ctx, userId)
	if err != nil {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	// get video id
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		log.Println("get video_id failed")
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "get video_id failed"})
		return
	}
	//get action
	//1-发布评论，2-删除评论
	actionType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		log.Println("get action_type failed")
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "get action_type failed"})
		return
	}
	var commentId int64
	var commentText string
	if actionType == 1 {
		//get commentText
		commentText = c.Query("comment_text")
		dal.InsertCommentVideoInfo(ctx, userId, videoId, commentText)
	} else if actionType == 2 {

		//get commentId
		commentId, err = strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			log.Println("get comment_id failed")
			c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "get comment_id failed"})
			return
		}
		dal.DeleteCommentInfo(ctx, commentId)
	}
	// get userById
	user := resp.User{
		Id:            u.Id,
		Name:          u.Name,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      true,
	}

	currentTime := time.Now()
	createDate := currentTime.Format("01-02")
	c.JSON(http.StatusOK, CommentActionResponse{
		Response: resp.Response{
			StatusCode: 0,
			StatusMsg:  fmt.Sprintf("%s action success", []string{"comment", "delete "}[actionType-1]),
		},
		Comment: resp.Comment{
			Id:         commentId,
			User:       user,
			Content:    commentText,
			CreateDate: createDate,
		},
	})
}

// CommentList all videos have same demo comment list
func CommentList(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	commentIdList, err := dal.GetCommentVideoIdxList(ctx, videoId)
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: resp.Response{
				StatusCode: 1,
				StatusMsg:  "get  video comment list failed",
			},
		})
		return
	}
	if len(commentIdList) == 0 {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: resp.Response{
				StatusCode: 0,
				StatusMsg:  "video comment is empty",
			},
		})
		return
	}
	commentsList := make([]resp.Comment, len(commentIdList))
	var wg sync.WaitGroup
	wg.Add(len(commentIdList))
	for i, id := range commentIdList {
		go func(j int, cid int64) {
			defer wg.Done()
			c, _ := dal.GetCommentById(ctx, cid)
			u, _ := dal.GetUserById(ctx, c.UserId)
			user := resp.User{
				Id:            u.Id,
				Name:          u.Name,
				FollowCount:   0,
				FollowerCount: 0,
				IsFollow:      true,
			}
			commentsList[j] = resp.Comment{
				Id:         int64(c.ID),
				User:       user,
				Content:    c.CommentText,
				CreateDate: c.CreatedAt.Format("01-02"),
			}
		}(i, id)

	}
	wg.Wait()
	c.JSON(http.StatusOK, CommentListResponse{
		Response: resp.Response{
			StatusCode: 0,
		},
		CommentList: commentsList,
	})
}
