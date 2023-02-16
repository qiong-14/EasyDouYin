package tests

import (
	"context"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/henrylee2cn/goutil"
	"github.com/minio/minio-go/v7"
	"github.com/qiong-14/EasyDouYin/dal"
	minioUtils "github.com/qiong-14/EasyDouYin/middleware"
	"math"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func init() {
	fmt.Println("init")
	dal.Init()
	minioUtils.Init(context.Background())
}

func TestDBInsert(t *testing.T) {
	// 测试插入数据
	for i := 0; i < 1000; i++ {
		name := faker.Email()
		passwd := faker.Password()
		// passwd md5
		passwdHashed := goutil.Sha1([]byte(passwd))
		if err := dal.CreateUser(context.Background(),
			&dal.User{Name: name, Password: passwdHashed}); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "insert data error: %d, %s, %s", i, name, passwd)
			if err != nil {
				//  do nothing
				t.Error("无法插入数据")
			}
		}
	}
}

// TestVideoMinio Minio测试
func TestVideoMinio(t *testing.T) {
	ctx := context.Background()
	//minioUtils.ListAllBuckets(ctx)

	minioUtils.GetFileList(ctx, nil, 10)
	videoURL, coverURL, err := minioUtils.GetUrlOfVideoAndCover(ctx, "v_ApplyEyeMakeup_g01_c01", time.Hour)
	if err != nil {
		t.Error("无法获取到预览链接")
	}

	t.Log(videoURL.String())
	t.Log(coverURL.String())
}

func TestVideoInsert(t *testing.T) {
	ctx := context.Background()
	minioUtils.GetFileList(ctx, func(info minio.ObjectInfo) {
		if strings.HasSuffix(info.Key, ".mp4") && strings.HasPrefix(info.Key, "v_") {
			title := info.Key[:len(info.Key)-4]
			fmt.Println(title)
			if err := dal.CreateVideoInfo(ctx, &dal.VideoInfo{
				Title:          title,
				OwnerId:        rand.Int63n(1000) + 1,
				LikesCount:     0,
				CommentArchive: 0,
				Label:          title[strings.Index(title, "_")+1 : 2+strings.Index(title[2:], "_")],
			}); err != nil {
				t.Error(err)
			}
		}
	}, math.MaxUint32)
}
