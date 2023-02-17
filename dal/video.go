package dal

import (
	"context"
	"fmt"
	"github.com/qiong-14/EasyDouYin/constants"
	"gorm.io/gorm"
	"os"
	"time"
)

// VideoInfo Table Struct
type VideoInfo struct {
	gorm.Model
	Title          string    `json:"title" gorm:"column:title"`
	OwnerId        int64     `json:"ownerId" gorm:"column:owner_id"`
	Label          string    `json:"label" gorm:"column:label"`
	LikesCount     int64     `json:"likesCount" gorm:"column:likes_count"`
	CommentArchive int64     `json:"commentArchive" gorm:"column:comment_archive"`
	CreatedAt      time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdateAt       time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP on update current_timestamp" json:"update_at,omitempty"`
}

var InvalidVideo = VideoInfo{
	Model: gorm.Model{ID: -1},
}

func (v *VideoInfo) VideoIsValid() bool {
	return v.ID == InvalidVideo.ID
}

// TableName videos table name
func (v *VideoInfo) TableName() string {
	return constants.VideoTableName
}

// CreateVideoInfo create user info
func CreateVideoInfo(ctx context.Context, v *VideoInfo) error {
	if err := DB.WithContext(ctx).Create(v).Error; err != nil {
		return err
	}
	return nil
}

// GetVideoStreamInfo 先返回时间倒序
func GetVideoStreamInfo(ctx context.Context, lastTime int64, limit int) (videos []VideoInfo) {
	videos = make([]VideoInfo, limit)
	tx := DB.WithContext(ctx).
		Model(&VideoInfo{}).
		Where("unix_timestamp(created_at) <= ?", lastTime*1000).
		Order("created_at desc").
		Limit(limit).Find(&videos)
	if err := tx.Error; err != nil {
		_, _ = fmt.Fprint(os.Stderr, "获取视频流错误")
	}
	return videos
}

// GetVideoInfoById GetVideoInfo get video info by id
func GetVideoInfoById(ctx context.Context, videoId int64) (video VideoInfo) {
	if err := DB.WithContext(ctx).
		Model(&VideoInfo{}).
		Where("id = ?", videoId).
		First(&video).Error; err != nil {
		_, _ = fmt.Fprint(os.Stderr, "获取视频错误")
	}
	return video
}

func GetPublishListById(ctx context.Context, userId int64) (videos []VideoInfo, err error) {
	videos = make([]VideoInfo, 100)
	if err := DB.WithContext(ctx).
		Model(&VideoInfo{}).
		Where("owner_id = ?", userId).
		Limit(100).Find(&videos).Error; err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "获取用户上传视频列表失败: ", err)
		return nil, err
	}
	return videos, nil
}
