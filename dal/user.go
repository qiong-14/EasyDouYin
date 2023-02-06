package dal

import (
	"context"
	"github.com/qiong-14/EasyDouYin/pkg/constants"
	"gorm.io/gorm"
)

// User Table Struct
type User struct {
	gorm.Model
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// TableName user table name
func (u User) TableName() string {
	return constants.UserTableName
}

// CreateUser create user info
func CreateUser(ctx context.Context, u *User) error {
	if err := DB.WithContext(ctx).Create(u).Error; err != nil {
		return err
	}
	return nil
}

// GetUserById get user info by id
func GetUserById(ctx context.Context, id int64) (*User, error) {
	u := &User{}
	if err := DB.WithContext(ctx).Model(&User{}).Where("id = ?", id).Find(u).Error; err != nil {
		return u, err
	}
	return u, nil
}

// GetUserByName get user info by name
func GetUserByName(ctx context.Context, name string) (*User, error) {
	u := &User{}
	if err := DB.WithContext(ctx).Model(&User{}).Where("name = ?", name).Find(u).Error; err != nil {
		return u, err
	}
	return u, nil

}
