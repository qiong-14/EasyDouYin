package feed

import (
	"context"
	"fmt"
	handler "github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	minioUtils "github.com/qiong-14/EasyDouYin/middleware"
	"math/rand"
	"testing"
	"time"
)

func init() {
	dal.Init()
	minioUtils.Init(context.Background())
	fmt.Println("init")
}

func TestVideoGet(t *testing.T) {
	videoInfos := dal.GetVideoStreamInfo(context.Background(), 0, 10)
	var videos = make([]handler.Video, 10)
	for i, info := range videoInfos {
		id := int64(info.ID)
		userInfo, _ := dal.GetUserById(context.Background(), id)
		playUrl, coverUrl, _ := minioUtils.GetUrlOfVideoAndCover(context.Background(),
			info.Title, time.Hour)
		videos = append(videos, handler.Video{
			Id: id,
			Author: handler.User{
				Id:            userInfo.ID,
				Name:          userInfo.Name,
				FollowCount:   int64(rand.Intn(1999)),
				FollowerCount: int64(rand.Intn(1000)),
				IsFollow:      false,
			},
			PlayUrl:       playUrl.String(),
			CoverUrl:      coverUrl.String(),
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
		})
		fmt.Printf("%d, %#v\n", i, videos[len(videos)-1])
	}
}
