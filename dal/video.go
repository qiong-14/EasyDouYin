package dal

import (
	"context"
	"fmt"
	"github.com/qiong-14/EasyDouYin/pkg/constants"
	"gorm.io/gorm"
	"os"
	"time"
)

// VideoInfo Table Struct
type VideoInfo struct {
	gorm.Model
	Title          string `json:"title" gorm:"column:title"`
	OwnerId        int64  `json:"ownerId" gorm:"column:owner_id"`
	Label          string `json:"label" gorm:"column:label"`
	LikesCount     int64  `json:"likesCount" gorm:"column:likes_count"`
	CommentArchive int64  `json:"commentArchive" gorm:"column:comment_archive"`
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
	if err := DB.WithContext(ctx).
		Model(&VideoInfo{}).
		Where("created_at > ?", time.Duration(lastTime)*time.Second).
		Order("created_at").
		Limit(limit).Find(&videos).Error; err != nil {
		_, _ = fmt.Fprint(os.Stderr, "获取视频流错误")
	}
	return videos
}
