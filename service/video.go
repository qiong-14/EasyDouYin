package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
)

// GetFavVideoList 通过用户ID获取该用户喜欢的视频ID列表
func GetFavVideoList(ctx context.Context, userId int64) []int64 {
	videoIdList, err := middleware.GetUserFavVideosRedis(userId)
	if videoIdList == nil || err != nil {
		videoIdList, _ = dal.GetLikeVideoIdxList(ctx, userId)
		for _, vId := range videoIdList {
			err := middleware.ActionUserFavVideoRedis(userId, vId)
			if err != nil {
				hlog.Errorf("%#v", err)
			}
		}
	}
	return videoIdList
}

// GetFavVideoCount 通过用户ID获取该用户喜欢的视频数量
func GetFavVideoCount(ctx context.Context, userId int64) int64 {
	cnt, err := middleware.GetUserFavVideosCountRedis(userId)
	if err != nil {
		videoList, err := dal.GetLikeVideoIdxList(ctx, userId)
		if err != nil {
			hlog.CtxErrorf(ctx, "can't get like video list , userId:%d", userId)
			return 0
		}
		for _, id := range videoList {
			err := middleware.ActionUserFavVideoRedis(userId, id)
			if err != nil {
				hlog.Errorf("%#v", err)
			}
		}
		cnt = int64(len(videoList))
	}
	return cnt
}

// GetVideoInfo 通过视频ID获取该视频信息
func GetVideoInfo(ctx context.Context, videoId int64) dal.VideoInfo {
	videoInfo, err := middleware.GetVideoInfoRedis(videoId)
	if err != nil {
		videoInfo = dal.GetVideoInfoById(ctx, videoId)

		err := middleware.SetVideoInfoRedis(videoInfo)
		if err != nil {
			hlog.Errorf("%#v", err)
		}
	}
	return videoInfo
}

// GetVideoFavUserCount 通过视频ID获取喜欢用户数量, 增加err
func GetVideoFavUserCount(ctx context.Context, videoId int64) (int64, error) {
	cnt, err := middleware.GetVideosFavsCountRedis(videoId)
	if err != nil {
		userId, err := dal.GetLikeUserList(ctx, videoId)
		if err != nil {
			hlog.CtxErrorf(ctx, "can't get like user list , videoId:%d", videoId)
			return 0, err
		}
		for _, id := range userId {
			err := middleware.ActionUserFavVideoRedis(id, videoId)
			if err != nil {
				hlog.Errorf("%#v", err)
			}
		}
		cnt = int64(len(userId))
	}
	return cnt, nil
}
