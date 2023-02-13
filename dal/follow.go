package dal

import (
	"context"
	"github.com/qiong-14/EasyDouYin/pkg/constants"
	"gorm.io/gorm"
	"log"
)

type Follow struct {
	gorm.Model
	UserId   int64 `gorm:"user_id"`
	FollowId int64 `json:"follow_id"`
	Cancel   int8  `json:"cancel"`
}

func (f Follow) TableName() string {
	return constants.FollowTableName
}

func FindRelation(ctx context.Context, userId, followId int64) (follow *Follow, err error) {
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Where("user_id = ? AND follow_id = ?", userId, followId).
		Take(&follow).Error; err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return
}

// GetFansCnt 获取粉丝
func GetFansCnt(ctx context.Context, userId int64) (cnt int64, err error) {
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Where("follow_id = ? AND cancel = ?", userId, 0).
		Count(&cnt).Error; nil != err {
		log.Println(err.Error())
		return 0, err
	}
	return
}

// GetFollowCnt 获取关注的人的数目
func GetFollowCnt(ctx context.Context, userId int64) (cnt int64, err error) {
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Where("user_id = ? AND cancel = ?", userId, 0).
		Count(&cnt).Error; nil != err {
		log.Println(err.Error())
		return 0, err
	}
	return
}

func CreateFollow(ctx context.Context, userId, followId int64, cancel int8) (err error) {
	follow := Follow{UserId: userId, FollowId: followId, Cancel: cancel}
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Create(&follow).Error; err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func IsFollowed(ctx context.Context, userId, followId int64) (follow *Follow, err error) {
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Where("user_id = ? AND follow_id = ?", userId, followId).
		Take(&follow).Error; nil != err {
		if err.Error() == "record not found" {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}
	return
}

func UpdateRelation(ctx context.Context, userId, followId int64, cancel int8) (err error) {
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Where("user_id = ? AND follow_id = ?", userId, followId).
		Update("cancel", cancel).Error; err != nil {
		log.Println(err.Error())
		return err
	}
	return
}

// GetFollowList 得到用户的关注列表
func GetFollowList(ctx context.Context, userId int64) (idx []int64, err error) {
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Where("user_id = ? AND cancel = ? ", userId, 0).
		Pluck("follow_id", &idx).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		log.Println("关注列表为空", err.Error())
		return nil, err
	}
	return
}

// GetFansList 得到用户的粉丝列表
func GetFansList(ctx context.Context, userId int64) (idx []int64, err error) {
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Where("follow_id = ?", userId).
		Pluck("user_id", &idx).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		log.Println("粉丝列表为空", err.Error())
		return nil, err
	}
	return
}
