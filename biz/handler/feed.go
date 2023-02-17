package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/constants"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
	"strconv"
	"sync"
	"time"
)

type FeedResponse struct {
	resp.Response
	VideoList []resp.Video `json:"video_list,omitempty"`
	NextTime  int64        `json:"next_time,omitempty"`
}

func getVideoEntities(ctx context.Context, videoInfos []dal.VideoInfo) []resp.Video {
	videosList := make([]resp.Video, len(videoInfos))

	var wg sync.WaitGroup

	wg.Add(len(videoInfos))
	for idx, info := range videoInfos {
		go func(resPos int, videoInfo dal.VideoInfo) {
			defer wg.Done()

			// 查询视频用户
			user, _ := dal.GetUserById(ctx, videoInfo.OwnerId)

			playUrl, coverUrl, _ := middleware.GetUrlOfVideoAndCover(context.Background(), videoInfo.Title, time.Hour)

			// 增加用户喜欢的查询
			favoriteCount, _ := dal.GetLikeUserCount(ctx, int64(videoInfo.ID))

			videosList[resPos] = resp.Video{
				Id: int64(videoInfo.ID),
				Author: resp.User{
					Id:            user.Id,
					Name:          user.Name,
					FollowCount:   0,
					FollowerCount: 0,
					IsFollow:      true,
				},
				PlayUrl:       playUrl.String(),
				CoverUrl:      coverUrl.String(),
				FavoriteCount: favoriteCount,
				CommentCount:  0,
				IsFavorite:    true,
			}
		}(idx, info)

	}
	wg.Wait()

	return videosList
}

func GetVideoStream(ctx context.Context, lastTime int64, limit int) []resp.Video {
	videoInfos := dal.GetVideoStreamInfo(ctx, lastTime, limit)

	return getVideoEntities(ctx, videoInfos)
}

func Feed(ctx context.Context, c *app.RequestContext) {
	//fmt.Println(c.Query("NextTime"))
	latestTimeStr := c.Query("latest_time")
	latestTime := int(time.Now().Unix())
	if latestTimeStr != "" {
		latestTime, _ = strconv.Atoi(latestTimeStr)
	}

	videoList := GetVideoStream(ctx, int64(latestTime), constants.FeedVideosCount)
	c.JSON(consts.StatusOK, FeedResponse{
		Response:  resp.Response{StatusCode: 0},
		VideoList: videoList,
		// todo: 需要替换成本次视频最小的时间戳
		NextTime: time.Now().Unix(),
	})
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}
