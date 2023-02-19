package dal

import (
	"context"

	"github.com/qiong-14/EasyDouYin/constants"
	"gorm.io/gorm"
)

// User_relation Table Struct
type User_relation struct {
	gorm.Model
	UserId        int64  `gorm:"primary_key"`
	Name          string `json:"name" gorm:"name"`
	Isfollow      bool   `json:"is_follow" gorm:"column:is_follow"`
	FollowerId    int64  `json:"follower_id" gorm:"column:follower_id"`
	FollowCount   int64  `json:"follow_count" gorm:"column:follow_count"`
	FollowerCount int64  `json:"follower_count" gorm:"column:follower_count"`
}

// User_info Table Struct
type User_info struct {
	gorm.Model
	UserId          int64  `gorm:"primary_key"`
	Name            string `json:"name" gorm:"name"`
	Isfollow        bool   `json:"is_follow" gorm:"column:is_follow"`
	FollowerId      int64  `json:"follower_id" gorm:"column:follower_id"`
	FollowCount     int64  `json:"follow_count" gorm:"column:follow_count"`
	FollowerCount   int64  `json:"follower_count" gorm:"column:follower_count"`
	Avatar          string `json:"avatar" gorm:"column:avatar"`
	BackgroundImage string `json:"background_image" gorm:"column:background_image"`
	Signature       string `json:"signature" gorm:"column:signature"`
	TotalFavorited  int64  `json:"total_favorited" gorm:"column:total_favorited"`
	WorkCount       int64  `json:"work_count" gorm:"column:work_count"`
	FavoriteCount   int64  `json:"favorite_count" gorm:"column:favorite_count"`
}

// TableName User_relation table name
func (ur User_relation) TableName() string {
	return constants.UserRelationName
}

// TableName User_relation table name
func (ur User_info) TableName() string {
	return constants.UserInfoName
}

// Get follower list by user_id
func GetFollowerList(ctx context.Context, UserId int64) ([]int64, error) {
	user_relations := []User_relation{}
	followerIds := []int64{}
	if err := DB.WithContext(ctx).Model(&User_relation{}).
		Where("user_id = ?", UserId).
		Select("follower").
		Find(&user_relations).Error; err != nil {
		return followerIds, err
	}
	for _, relation := range user_relations {
		followerIds = append(followerIds, relation.FollowerId)
	}
	return followerIds, nil
}

// Get Friend list by user_id
func GetFriendList(ctx context.Context, UserId int64) ([]int64, error) {
	user_relations := []User_relation{}
	friendIds := []int64{}
	if err := DB.WithContext(ctx).Model(&User_relation{}).
		Table("user_relations AS a").
		Select("UserId").
		Joins("user_relations AS b ON a.user_id = b.follower_id AND a.follower_id = b.user_id").
		Where("user_id = ?", UserId).
		Find(&user_relations).Error; err != nil {
		return friendIds, err
	}
	for _, relation := range user_relations {
		friendIds = append(friendIds, relation.UserId)
	}
	return friendIds, nil
}

func GetAllUserInfo(ctx context.Context, UserIds []int64) ([]User_info, error) {
	infos := []User_info{}
	for _, idx := range UserIds {
		info, err := GetUserInfo(ctx, UserIds[idx])
		if err != nil {
			return infos, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// Get user info by user_id
func GetUserInfo(ctx context.Context, UserId int64) (User_info, error) {
	info := User_info{}
	if err := DB.WithContext(ctx).Model(&User_info{}).
		Where("user_id = ?", UserId).
		Find(&info).Error; err != nil {
		return info, err
	}
	return info, nil
}
