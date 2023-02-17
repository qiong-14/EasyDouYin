package dal

import (
	"context"
	"github.com/qiong-14/EasyDouYin/constants"
	"gorm.io/gorm"
)

// Message Table Struct
type Message struct {
	gorm.Model
	Id			int64  `gorm:"primary_key"`
	ToUserId	int64  `json:"to_user_id" gorm:"column:to_user_id"`
	FromUserId	int64  `json:"from_user_id" gorm:"column:from_user_id"`
	Content    	string `json:"content" gorm:"column:content"`
	CreateTime 	int64  `json:"create_time" gorm:"column:create_time"`
}

// TableName Message table name
func (m Message) TableName() string {
	return constants.MessageTableName
}

// Insert message
func InsertMessage(ctx context.Context, m *Message) error {
	if err := DB.
		WithContext(ctx).
		Model(&Message{}).
		Create(m).Error; err != nil {
		return err
	}
	return nil
}

// Get message by to_user_id,from_user_id
func GetMessage(ctx context.Context, toUserId int64, FromUserId int64) ([]Message, error) {
	messages := []Message{}
	if err := DB.WithContext(ctx).Model(&Message{}).
	Where("to_user_id = ? AND from_user_id = ?", toUserId,FromUserId).
	Or("to_user_id = ? AND from_user_id = ?",FromUserId, toUserId).
	Find(&messages).Error; err != nil {
		return messages, err
	}
	return messages, nil
}
