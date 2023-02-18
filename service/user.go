package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
)

// GetUserInfo 通过ID获取用户信息 带redis缓存, han bing 增加了一个error
func GetUserInfo(ctx context.Context, userId int64) (dal.User, error) {
	user, err := middleware.GetUserInfoRedis(userId)
	if user == dal.InvalidUser || err != nil {
		user, err = dal.GetUserById(ctx, userId)
		if err != nil {
			hlog.CtxErrorf(ctx, "can't get user info by id: %d", userId)
			return dal.InvalidUser, err
		}
		if err = middleware.SetUserInfoRedis(user); err != nil {
			hlog.CtxErrorf(ctx, "can't set user info cache by id: %d", userId)
			return dal.InvalidUser, err
		}
	}
	return user, nil
}
