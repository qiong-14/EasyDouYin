package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/qiong-14/EasyDouYin/biz/common"
	"github.com/qiong-14/EasyDouYin/dal"
	"path/filepath"
	"strconv"
)

type VideoListResponse struct {
	common.Response
	VideoList []common.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(ctx context.Context, c *app.RequestContext) {
	id := c.Query("user_id")
	userId, _ := strconv.ParseInt(id, 10, 64)
	user, err := dal.GetUserById(ctx, userId)
	if err != nil {
		c.JSON(consts.StatusOK, common.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(consts.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err = c.SaveUploadedFile(data, saveFile); err == nil {
		c.JSON(consts.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, VideoListResponse{
		Response: common.Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
