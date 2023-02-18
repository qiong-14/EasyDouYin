package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/jwt"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"github.com/qiong-14/EasyDouYin/dal"
	"github.com/qiong-14/EasyDouYin/middleware"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type VideoListResponse struct {
	resp.Response
	VideoList []resp.Video `json:"video_list"`
}

// getCoverAndUpload 获取视频封面并上传
func getCoverAndUpload(ctx context.Context, filePath string, userId int64) error {
	err := middleware.UploadVideoAndCover(ctx, filePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "can not upload files to minio")
		return err
	}
	baseName := filepath.Base(filePath)
	titleNoExt := baseName[:len(baseName)-4]
	return dal.CreateVideoInfo(ctx, &dal.VideoInfo{
		Title:          titleNoExt,
		OwnerId:        userId,
		Label:          filePath,
		LikesCount:     0,
		CommentArchive: 0,
	})
}

// Publish check token then save upload file to public directory
func Publish(ctx context.Context, c *app.RequestContext) {
	userId := 0
	if identity, exist := c.Get(jwt.IdentityKey); exist {
		// 获取一下
		if user, exist := identity.(*dal.User); exist {
			userId = int(user.Id)
		}
	}
	user, err := dal.GetUserById(ctx, int64(userId))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%#v", err)
		c.JSON(consts.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	// 保存到本地
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(consts.StatusOK, resp.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s_%d.mp4", user.Id, filename, time.Now().UnixNano())
	saveFile := filepath.Join("./public/", finalName)

	if err = c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(consts.StatusOK, resp.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 上传视频
	err = getCoverAndUpload(ctx, saveFile, int64(userId))
	if err != nil {
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
	userId, _ := strconv.Atoi(c.Query("user_id"))

	videosInfos, err := dal.GetPublishListById(ctx, int64(userId))
	if err != nil {
		hlog.Error(err)
		return
	}

	videosList := getVideoEntities(ctx, videosInfos)

	c.JSON(consts.StatusOK, VideoListResponse{
		Response: resp.Response{
			StatusCode: 0,
		},
		VideoList: videosList,
	})
	hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
}
