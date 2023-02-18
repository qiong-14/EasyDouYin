package dal

import (
	"context"
	"github.com/qiong-14/EasyDouYin/constants"
	"gorm.io/gorm"
	"log"
)

type Like struct {
	gorm.Model
	UserId  int64 `json:"user_id" gorm:"colum:user_id"`
	VideoId int64 `json:"video_id" gorm:"colum:video_id"`
	Cancel  int8  `json:"cancel" gorm:"colum:cancel"`
}

func (L Like) TableName() string {
	return constants.LikeVideoTableName
}

// InsertLikeVideoInfo insert user like video info
func InsertLikeVideoInfo(ctx context.Context, userId, videoId int64, cancel int8) error {
	if err := DB.
		WithContext(ctx).
		Model(&Like{}).
		Create(&Like{UserId: userId, VideoId: videoId, Cancel: cancel}).Error; err != nil {
		log.Println("insert like video failed")
		return err
	}
	log.Println("insert like video success")
	return nil
}

// GetLikeVideoIdxList return video id lists user likes by update time
func GetLikeVideoIdxList(ctx context.Context, userId int64) ([]int64, error) {
	var likeVideoIdxList []int64
	if err := DB.Table("(?) as u", DB.
		WithContext(ctx).
		Model(&Like{}).
		Where(&Like{UserId: userId, Cancel: 1}).
		Distinct("video_id", "updated_at").
		Order("updated_at desc")).
		Pluck("video_id", &likeVideoIdxList).Error; err != nil {
		log.Println("get like video list failed")
		return nil, err
	}
	//log.Println("get like video list success:", likeVideoIdxList)
	return likeVideoIdxList, nil
}

// GetLikeUserCount get the number of users who like the video
func GetLikeUserCount(ctx context.Context, videoId int64) (int64, error) {
	var cnt int64
	if err := DB.
		WithContext(ctx).
		Model(&Like{}).
		Where(&Like{VideoId: videoId, Cancel: 1}).
		Distinct("user_id").
		Count(&cnt).Error; err != nil {
		log.Printf("no users like video %d", videoId)
		return 0, err
	}
	//log.Printf("get %d users like video %d", cnt, videoId)
	return cnt, nil
}

// GetLikeUserList get all user's id who like this video
func GetLikeUserList(ctx context.Context, videId int64) ([]int64, error) {
	var userId []int64
	if err := DB.
		WithContext(ctx).
		Model(&Like{}).
		Where(&Like{VideoId: videId, Cancel: 1}).
		Distinct("user_id").
		Pluck("user_id", &userId).Error; err != nil {
		return nil, err
	}
	return userId, nil
}

// FindLikeVideoInfo find relation record if not find return nil
func FindLikeVideoInfo(ctx context.Context, userId, videoId int64) (Like, error) {
	var likeVideoInfo Like
	if err := DB.
		WithContext(ctx).
		Model(&Like{}).
		Where(&Like{UserId: userId, VideoId: videoId}).
		First(&likeVideoInfo).Error; err != nil {
		return Like{}, err
	}
	//log.Println("find like video info success")
	return likeVideoInfo, nil
}

// UpdateLikeInfo update like relation 1-like 2-cancel
func UpdateLikeInfo(ctx context.Context, userId, videoId int64, cancel int8) error {
	if err := DB.
		WithContext(ctx).
		Model(&Like{}).
		Where(&Like{UserId: userId, VideoId: videoId}).
		Update("cancel", cancel).Error; err != nil {
		log.Println("update like video info failed:", err.Error())
		return err
	}
	//log.Printf("update  %s video info success", []string{"like", "dislike"}[int8(cancel)-1])
	return nil
}
