package dal

import (
	"context"
	"log"

	"github.com/qiong-14/EasyDouYin/constants"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserId      int64  `json:"user_id" gorm:"user_id"`
	VideoId     int64  `json:"video_id" gorm:"video_id"`
	CommentText string `json:"commenttext" gorm:"comment_text"`
}

func (C Comment) TableName() string {
	return constants.CommentVideoTableName
}

// InsertCommentVideoInfo insert user comment video info
func InsertCommentVideoInfo(ctx context.Context, userId, videoId int64, commentText string) (int64, error) {
	comment := Comment{UserId: userId, VideoId: videoId, CommentText: commentText}
	if err := DB.
		WithContext(ctx).
		Model(&Comment{}).
		Create(&comment).Error; err != nil {
		log.Println("insert comment failed")
		return 0, err
	}
	log.Println("insert comment success")
	return int64(comment.ID), nil
}

// GetCommentVideoIdxList return comment id lists of video by update time
func GetCommentVideoIdxList(ctx context.Context, videoId int64) ([]int64, error) {
	var commentVideoIdxList []int64
	if err := DB.
		WithContext(ctx).
		Model(&Comment{}).
		Where(&Comment{VideoId: videoId}).
		Order("updated_at desc").
		Pluck("id", &commentVideoIdxList).Error; err != nil {
		log.Println("get comment list failed")
		return nil, err
	}
	//log.Println("get like video list success:", likeVideoIdxList)
	return commentVideoIdxList, nil
}

// GetCommentById get comment by id
func GetCommentById(ctx context.Context, id int64) (*Comment, error) {
	c := &Comment{}
	if err := DB.
		WithContext(ctx).
		Model(&Comment{}).
		Where("id = ? ", id).
		First(c).Error; err != nil {
		log.Printf("get comment failed")
		return c, err
	}
	//log.Printf("get %d users like video %d", cnt, videoId)
	return c, nil
}

// DeleteCommentInfo delete comment by Id
func DeleteCommentInfo(ctx context.Context, id int64) error {
	if err := DB.
		WithContext(ctx).
		Delete(&Comment{}, id).Error; err != nil {
		log.Println("delete comment info failed:", err.Error())
		return err
	}
	//log.Printf("update  %s video info success", []string{"like", "dislike"}[int8(cancel)-1])
	return nil
}
