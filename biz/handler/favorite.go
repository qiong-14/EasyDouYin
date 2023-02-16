package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiong-14/EasyDouYin/biz/resp"
	"net/http"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
	//u, _ := c.Get(middleware.IdentityKey)
	//userId := u.(*dal.User).Id
	//videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	//if err != nil {
	//	log.Println("get video_id failed")
	//	c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "get video_id failed"})
	//	return
	//}
	c.JSON(http.StatusOK, resp.Response{StatusCode: 0, StatusMsg: "get video_id failed"})

	//actionType, err := strconv.Atoi(c.Query("action_type"))
	//if err != nil {
	//	log.Println("get action_type failed")
	//	c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "get action_type failed"})
	//	return
	//}

	//if _, err := dal.FindLikeVideoInfo(ctx, userId, videoId); err != nil {
	//
	//}

}

// FavoriteList all users have same favorite video list
func FavoriteList(ctx context.Context, c *app.RequestContext) {
	defer hlog.CtxTracef(ctx, "status=%d method=%s full_path=%s client_ip=%s host=%s",
		c.Response.StatusCode(),
		c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
	c.JSON(http.StatusOK, VideoListResponse{
		Response: resp.Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
