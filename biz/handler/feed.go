package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/pkg/constants"
	minioUtils "github.com/qiong-14/EasyDouYin/utils/minio"
	"math/rand"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

func GetVideoStream(ctx context.Context, lastTime int64, limit int) []Video {
	videoInfos := dal.GetVideoStreamInfo(ctx, lastTime, limit)
	var videos []Video
	for _, info := range videoInfos {
		id := int64(info.ID)
		userInfo, _ := dal.GetUserById(context.Background(), id)
		playUrl, coverUrl, _ := minioUtils.GetUrlOfVideoAndCover(context.Background(),
			info.Title, time.Hour)
		fmt.Println(playUrl.String())
		fmt.Println("cover", coverUrl.String())
		video := &Video{
			Id: id,
			Author: User{
				Id:            userInfo.ID,
				Name:          userInfo.Name,
				FollowCount:   int64(rand.Intn(1999)), // 随机给的
				FollowerCount: int64(rand.Intn(1000)),
				IsFollow:      false,
			},
			PlayUrl:       playUrl.String(),
			CoverUrl:      coverUrl.String(),
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
		}
		videos = append(videos, *video)
	}
	return videos
}
func Feed(ctx context.Context, c *app.RequestContext) {
	fmt.Println(c.Query("NextTime"))
	c.JSON(consts.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: GetVideoStream(ctx, 0, constants.FeedVideosCount),
		NextTime:  time.Now().Unix(),
	})
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
			c.Response.StatusCode(),
			 c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}
