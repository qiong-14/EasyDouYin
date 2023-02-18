package redis

import (
	"fmt"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
	"gorm.io/gorm"
	"testing"
)

func init() {
	middleware.InitRedis()
	fmt.Println("init")
}

func TestSetUserInfoRedis(t *testing.T) {
	err := middleware.SetUserInfoRedis(dal.User{
		Id:   1001,
		Name: "863178540@qq.com",
	})
	if err != nil {
		t.Error(err)
	}
	// Get
	res, err := middleware.GetUserInfoRedis(1001)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n", res)
}

func TestSetVideoInfoRedis(t *testing.T) {
	err := middleware.SetVideoInfoRedis(dal.VideoInfo{
		Model: gorm.Model{
			ID: 3,
		},
		Title:          "123",
		OwnerId:        1001,
		Label:          "",
		LikesCount:     0,
		CommentArchive: 0,
	})
	if err != nil {
		t.Error(err)
	}
	// Get VideoInfo
	res, err := middleware.GetVideoInfoRedis(3)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n", res)
	// Get UserInfo
	resUser, err := middleware.GetUserInfoRedis(res.OwnerId)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n", resUser)
}

// other test
func TestGetUserFavRedis(t *testing.T) {
	var err error

	err = middleware.ActionUserFavVideoRedis(1001, 3)
	if err != nil {
		t.Error(err)
	}

	err = middleware.ActionUserFavVideoRedis(1001, 4)
	if err != nil {
		t.Error(err)
	}

	err = middleware.ActionUserFavVideoRedis(1002, 4)
	if err != nil {
		t.Error(err)
	}

	err = middleware.ActionUserFavVideoRedis(1003, 4)
	if err != nil {
		t.Error(err)
	}

	res, err := middleware.GetUserFavVideosRedis(1001)
	if err != nil {
		t.Error(err)
	}
	if len(res) != 2 {
		t.Error("数量不对")
	}
	t.Log(res)

	res, err = middleware.GetVideosFavRedis(4)
	if err != nil {
		t.Error(err)
	}
	if len(res) != 3 {
		t.Error("数量不对")
	}
	t.Log(res)

	// 最后测一下删除

	err = middleware.ActionUserUnFavVideoRedis(1001, 4)
	if err != nil {
		t.Error(err)
	}
	t.Log("删除 1001, 4")
	res, err = middleware.GetVideosFavRedis(4)
	if err != nil {
		t.Error(err)
	}
	if len(res) != 2 {
		t.Error("数量不对")
	}
	t.Log(res)

}
