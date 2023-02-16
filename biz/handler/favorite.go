package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
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
		dal.UpdateLikeInfo(ctx, userId, videoId, int8(actionType))
	} else {
		dal.InsertLikeVideoInfo(ctx, userId, videoId, int8(actionType))
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

	videoIdList, err := dal.GetLikeVideoIdxList(ctx, userId)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: resp.Response{
				StatusCode: 1,
				StatusMsg:  "get like video id list failed",
			},
		})
		return
	}
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
	var wg sync.WaitGroup
	wg.Add(len(videoIdList))
	for i, id := range videoIdList {
		go func(j int, vid int64) {
			defer wg.Done()
			v := dal.GetVideoInfoById(ctx, vid)
			user, _ := dal.GetUserById(ctx, v.OwnerId)
			author := resp.User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   0,
				FollowerCount: 0,
				IsFollow:      true,
			}
			playUrl, coverUrl, _ := middleware.GetUrlOfVideoAndCover(context.Background(),
				v.Title, time.Hour)
			favoriteCount, _ := dal.GetLikeUserCount(ctx, int64(v.ID))
			videosList[j] = resp.Video{
				Id:            int64(v.ID),
				Author:        author,
				PlayUrl:       playUrl.String(),
				CoverUrl:      coverUrl.String(),
				FavoriteCount: favoriteCount,
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
