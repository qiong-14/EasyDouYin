package dal

import (
	"context"
	"github.com/qiong-14/EasyDouYin/pkg/constants"
	"log"
)

type Follow struct {
	Id       int64 `json:"id" gorm:"primary_key"`
	UserId   int64 `json:"user_id" gorm:"colum:user_id"`
	FollowId int64 `json:"follow_id" gorm:"follow_id"`
}

// TableName follow
func (f Follow) TableName() string {
	return constants.FollowTableName
}

// FindRelation 是否有过关注记录
func FindRelation(ctx context.Context, userId, followId int64) (y bool, err error) {
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Where(&Follow{UserId: userId, FollowId: followId}).
		First(&Follow{}).Error; err != nil {
		log.Println(err.Error())
		return false, err
	}
	return true, nil
}

// GetFansCnt 获取粉丝数目
func GetFansCnt(ctx context.Context, followId int64) (cnt int64, err error) {
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Where(&Follow{FollowId: followId}).
		Count(&cnt).Error; nil != err {
		log.Println(err.Error())
		return 0, err
	}
	return
}

// GetFollowCnt 获取关注数目
func GetFollowCnt(ctx context.Context, userId int64) (cnt int64, err error) {
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Where(&Follow{UserId: userId}, userId).
		Count(&cnt).Error; nil != err {
		log.Println(err.Error())
		return 0, err
	}
	return
}

// CreateRelation 创建关注记录
func CreateRelation(ctx context.Context, userId, followId int64) (err error) {
	follow := Follow{UserId: userId, FollowId: followId}
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Create(&follow).Error; err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// DeleteRelation 删除关注记录
func DeleteRelation(ctx context.Context, userId, followId int64) (err error) {
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Where(&Follow{UserId: userId, FollowId: followId}).
		Delete(&Follow{}).Error; err != nil {
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
		Where(&Follow{UserId: userId}).
		Pluck("follow_id", &idx).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		log.Println("关注列表为空：", err.Error())
		return nil, err
	}
	return
}

// GetFansList 得到用户的粉丝列表
func GetFansList(ctx context.Context, followId int64) (idx []int64, err error) {
	if err = DB.
		WithContext(ctx).
		Model(&Follow{}).
		Where(&Follow{FollowId: followId}).
		Pluck("user_id", &idx).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		log.Println("粉丝列表为空", err.Error())
		return nil, err
	}
	return
}
