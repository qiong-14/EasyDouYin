package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	minioUtils "github.com/qiong-14/EasyDouYin/mw"
	"github.com/qiong-14/EasyDouYin/pkg/constants"
	"time"
)

type FeedResponse struct {
	resp.Response
	VideoList []resp.Video `json:"video_list,omitempty"`
	NextTime  int64        `json:"next_time,omitempty"`
}

func GetVideoStream(ctx context.Context, lastTime int64, limit int) []resp.Video {
	videoInfos := dal.GetVideoStreamInfo(ctx, lastTime, limit)
	var videos []resp.Video
	for _, info := range videoInfos {
		id := int64(info.ID)
		userInfo, _ := dal.GetUserById(context.Background(), id)
		playUrl, coverUrl, _ := minioUtils.GetUrlOfVideoAndCover(context.Background(),
			info.Title, time.Hour)
		video := &resp.Video{
			Id:            id,
			Author:        dal.GetRespUser(ctx, userInfo.Id),
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
	//fmt.Println(c.Query("NextTime"))
	c.JSON(consts.StatusOK, FeedResponse{
		Response:  resp.Response{StatusCode: 0},
		VideoList: GetVideoStream(ctx, 0, constants.FeedVideosCount),
		NextTime:  time.Now().Unix(),
	})
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}
