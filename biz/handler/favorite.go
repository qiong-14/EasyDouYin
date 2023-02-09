package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/qiong-14/EasyDouYin/biz/common"
	"github.com/qiong-14/EasyDouYin/dal"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(ctx context.Context, c *app.RequestContext) {
	id := c.Query("user_id")
	userId, _ := strconv.ParseInt(id, 10, 64)
	if _, err := dal.GetUserById(ctx, userId); err == nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(ctx context.Context, c *app.RequestContext) {
	//c.JSON(http.StatusOK, VideoListResponse{
	//	Response: Response{
	//		StatusCode: 0,
	//	},
	//	VideoList: DemoVideos,
	//})
}
