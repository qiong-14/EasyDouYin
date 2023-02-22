package dal

import (
	"context"
	"github.com/qiong-14/EasyDouYin/constants"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type Follows struct {
	gorm.Model
	Id         int64 `json:"id" gorm:"primary_key"`
	FollowedId int64 `json:"followed_id" gorm:"column:followed_id"`
	FollowerId int64 `json:"follower_id" gorm:"column:follower_id"`
	//actionType: 1-follow; 2-unfollow
	ActionType int `json:"action_type" gorm:"column:cancel"`
}

type UserVo struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	IsFollow bool   `json:"is_follow"`
}

func (f Follows) TableName() string {
	return constants.FollowTableName
}

func UpdateOrCreateRelation(ctx context.Context, f *Follows) error {
	if DB.Model(&f).Where(&Follows{FollowedId: f.FollowedId, FollowerId: f.FollowerId}).Updates(&f).RowsAffected == 0 {
		if err := CreateRelation(ctx, f); err != nil {
			return err
		}
	}
	return nil
}

func CreateRelation(ctx context.Context, f *Follows) error {
	if err := DB.
		WithContext(ctx).
		Model(&Follows{}).
		Create(f).Error; err != nil {
		return err
	}
	return nil
}

func FollowCount(ctx context.Context, userId int64) (int64, error) {
	var total int64
	if err := DB.
		WithContext(ctx).
		Model(&Follows{}).
		Where(&Follows{FollowerId: userId, ActionType: constants.RelationFollow}).
		Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func FollowerCount(ctx context.Context, userId int64) (int64, error) {
	var total int64
	if err := DB.
		WithContext(ctx).
		Model(&Follows{}).
		Where(&Follows{FollowedId: userId, ActionType: constants.RelationFollow}).
		Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// get the userVo list of those who are followed by the user
// 用户的关注列表，返回值为用户关注的人的userVo集合

func FollowUserList(ctx context.Context, userId int64) ([]UserVo, error) {
	followIdList, err := FollowIdList(ctx, userId)
	if err != nil {
		return nil, err
	}
	userVoList, err := getUserVoList(ctx, followIdList, userId, true)
	if err != nil {
		return nil, err
	}
	return userVoList, nil
}

// get the userVo list of those who follows the user
// 用户的粉丝列表，返回值为用户粉丝的userVo集合

func FollowerUserList(ctx context.Context, userId int64) ([]UserVo, error) {
	followerIdList, err := FollowerIdList(ctx, userId)
	if err != nil {
		return nil, err
	}
	userVoList, err := getUserVoList(ctx, followerIdList, userId, false)
	if err != nil {
		return nil, err
	}
	return userVoList, nil
}

// get the userVo list of those who follows the user and are followed by the user
// 用户的好友列表（好友指互相关注），返回值为用户好友的userVo集合
// for A in the follower list of the user, if A.IsFollow = true, A is the friend of the user

func FriendUserList(ctx context.Context, userId int64) ([]UserVo, error) {
	followerUserList, err := FollowerUserList(ctx, userId)
	if err != nil {
		return nil, err
	}
	var friendUserList []UserVo
	for _, user := range followerUserList {
		if user.IsFollow {
			friendUserList = append(friendUserList, user)
		}
	}
	return friendUserList, nil
}

// get the id list of those who are followed by the user
// 用户的关注列表，返回值为用户关注的人的userid数组

func FollowIdList(ctx context.Context, userId int64) ([]int64, error) {
	var followIds []int64

	if err := DB.
		WithContext(ctx).
		Model(&Follows{}).
		Where(&Follows{FollowerId: userId, ActionType: constants.RelationFollow}).
		Pluck("followed_id", &followIds).Error; err != nil {
		log.Println("fail to find follow id")
		return nil, err
	}
	return followIds, nil
}

// get the id list of those who follows the user
// 用户的粉丝列表，返回值为用户粉丝的userid数组

func FollowerIdList(ctx context.Context, userId int64) ([]int64, error) {
	var followerIds []int64

	if err := DB.
		WithContext(ctx).
		Model(&Follows{}).
		Where(&Follows{FollowedId: userId, ActionType: constants.RelationFollow}).
		Pluck("follower_id", &followerIds).Error; err != nil {
		log.Println("fail to find follower id")
		return nil, err
	}
	return followerIds, nil
}

func getUserVoList(ctx context.Context, userIdList []int64, olduserId int64, followed bool) ([]UserVo, error) {
	var userVoList []UserVo
	var userVo *UserVo
	var err error
	for _, val := range userIdList {
		userVo, err = idToUserVo(ctx, val, olduserId, followed)
		if err != nil {
			return nil, err
		}
		userVoList = append(userVoList, *userVo)
	}
	return userVoList, nil
}

//get the userVo given newuserId.
//newUser follows oldUser, or is followed by the oldUser
//we need to find out if oldUser follows newUser.
//if we already know oldUser follows newUser, followed is true

func idToUserVo(ctx context.Context, newuserId int64, olduserId int64, followed bool) (*UserVo, error) {
	var userVo UserVo
	userVo.Id = newuserId
	var userName string
	if err := DB.
		WithContext(ctx).
		Model(&User{}).
		Select("name").
		Where(&User{Id: newuserId}).
		First(&userName).Error; err != nil {
		log.Println("failed to get username of user" + strconv.FormatInt(newuserId, 10))
		return nil, err
	}
	userVo.Name = userName
	if followed {
		userVo.IsFollow = true
		return &userVo, nil
	}

	var count int64
	if err := DB.
		WithContext(ctx).
		Model(&Follows{}).
		Where(&Follows{FollowedId: newuserId, FollowerId: olduserId, ActionType: constants.RelationFollow}).
		Count(&count).Error; err != nil {
		log.Println("failed to get the relation between users")
		return nil, err
	}
	userVo.IsFollow = count == 1

	return &userVo, nil



}
