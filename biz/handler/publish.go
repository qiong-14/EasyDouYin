package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/mw"
	"path/filepath"
)

type VideoListResponse struct {
	resp.Response
	VideoList []resp.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(ctx context.Context, c *app.RequestContext) {
	u, _ := c.Get(mw.IdentityKey)
	user, err := dal.GetUserById(ctx, u.(*dal.User).Id)
	if err == nil {
		c.JSON(consts.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(consts.StatusOK, resp.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err = c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(consts.StatusOK, resp.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, resp.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}

// PublishList all users have same publish video list
func PublishList(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, VideoListResponse{
		Response: resp.Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}
