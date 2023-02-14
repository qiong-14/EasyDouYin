package dal

import (
	"context"
	"github.com/qiong-14/EasyDouYin/pkg/constants"
	"time"
)

type Message struct {
	Id         int64  `gorm:"primary_key"`
	ToUserId   int64  `json:"to_user_id" gorm:"to_user_id"`
	FromUserID int64  `json:"from_user_id" gorm:"from_user_id"`
	Content    string `json:"content" gorm:"content"`
	CreatedAt  time.Time
}

func (m Message) TableName() string {
	return constants.MessageTableName
}

func CreateMessage(ctx context.Context, toUserId, fromUserId int64, content string) error {
	if err := DB.
		WithContext(ctx).
		Model(&Message{}).
		Create(&Message{ToUserId: toUserId, FromUserID: fromUserId, Content: content}).Error; err != nil {
		return err
	}
	return nil
}
