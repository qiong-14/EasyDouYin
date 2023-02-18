package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
	"github.com/qiong-14/EasyDouYin/service"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// FavoriteAction change like relation to like_video db
func FavoriteAction(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
	// get user id
	u, _ := c.Get(middleware.IdentityKey)
	userId := u.(*dal.User).Id

	// get video id
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		log.Println("get video_id failed")
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "get video_id failed"})
		return
	}

	// get action type 1-like 2-cancel like
	actionType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		log.Println("get action_type failed")
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "get action_type failed"})
		return
	}

	// find info if exist(err == nil) update else create
	if _, err := dal.FindLikeVideoInfo(ctx, userId, videoId); err == nil {
		// han bing 2023年02月16日22:46:44 传至最上层, 忽略错误
		_ = dal.UpdateLikeInfo(ctx, userId, videoId, int8(actionType))
	} else {
		_ = dal.InsertLikeVideoInfo(ctx, userId, videoId, int8(actionType))
	}
	if actionType == 1 {
		middleware.ActionUserFavVideoRedis(userId, videoId)
	} else {
		middleware.ActionUserUnFavVideoRedis(userId, videoId)
	}
	c.JSON(http.StatusOK, resp.Response{
		StatusCode: 0,
		StatusMsg:  fmt.Sprintf("%s video action success", []string{"like", "dislike"}[actionType-1]),
	},
	)

}

// FavoriteList all users have same favorite video list
func FavoriteList(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
	// get user id
	u, _ := c.Get(middleware.IdentityKey)
	userId := u.(*dal.User).Id

	// han bing 2023年02月16日22:31:58 做join是否会更快一点? like_videos和videos
	videoIdList := service.GetFavVideoList(ctx, userId)
	if len(videoIdList) == 0 {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: resp.Response{
				StatusCode: 0,
				StatusMsg:  "like videos is empty",
			},
		})
		return
	}
	videosList := make([]resp.Video, len(videoIdList))
	// 并发
	var wg sync.WaitGroup
	wg.Add(len(videoIdList))
	for i, id := range videoIdList {
		go func(j int, vid int64) {
			defer wg.Done()
			// han bing 2023年02月16日23:36:04 vid 可能是 0, 可能会在之后的流程中被过滤
			video := service.GetVideoInfo(ctx, vid)
			user := service.GetUserInfo(ctx, userId)
			author := resp.User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   0,
				FollowerCount: 0,
				IsFollow:      true,
			}
			playUrl, coverUrl, _ := middleware.GetUrlOfVideoAndCover(context.Background(),
				video.Title, time.Hour)
			videosList[j] = resp.Video{
				Id:            int64(video.ID),
				Author:        author,
				PlayUrl:       playUrl.String(),
				CoverUrl:      coverUrl.String(),
				FavoriteCount: service.GetVideoFavUserCount(ctx, vid),
				CommentCount:  0,
				IsFavorite:    true,
			}
		}(i, id)
	}
	wg.Wait()
	c.JSON(http.StatusOK, VideoListResponse{
		Response: resp.Response{
			StatusCode: 0,
		},
		VideoList: videosList,
	})
}
