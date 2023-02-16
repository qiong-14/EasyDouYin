package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/qiong-14/EasyDouYin/tools"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"net/url"
	"path/filepath"
	"time"
)

var (
	endpoint        = tools.GetEnvByKey("MINIO_ENDPOINT")
	accessKeyID     = tools.GetEnvByKey("MINIO_ACCESS_KEY")
	secretAccessKey = tools.GetEnvByKey("MINIO_SECRET_KEY")
	bucketName      = tools.GetEnvByKey("MINIO_BUCKET")
	Client          *minio.Client
)

func Init(ctx context.Context) {
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%#v\n", minioClient)

	if err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"}); err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("1, We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
	Client = minioClient
}

func GetFileList(ctx context.Context, handler func(info minio.ObjectInfo), maxCount int) (res []any) {
	ctx = context.Background()
	idx := 0

	for objInfo := range Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{}) {
		if handler != nil {
			handler(objInfo)
		} else {
			fmt.Println(objInfo.Key)
		}
		idx += 1
		if idx == maxCount {
			break
		}
	}
	return res
}

func GetUrlOfVideoAndCover(ctx context.Context, name string, timeout time.Duration) (videoURL, coverURL *url.URL, err error) {
	videoURL, err = getFileUrl(ctx, name+".mp4", timeout)
	if err != nil {
		return nil, nil, err
	}
	coverURL, err = getFileUrl(ctx, name+".mp4.jpg", timeout)
	if err != nil {
		return nil, nil, err
	}
	return
}

func UploadVideoAndCover(ctx context.Context, fileName string) error {
	baseName := filepath.Base(fileName)

	uploadFile(ctx, fileName, "application/octet-stream", baseName)

	coverName, err := getCover(ctx, fileName)
	if err != nil {
		return err
	}
	uploadFile(ctx, coverName, "application/octet-stream", filepath.Base(coverName))
	return nil
}

func ListAllBuckets(ctx context.Context) []minio.BucketInfo {
	buckets, _ := Client.ListBuckets(context.Background())
	for _, bucket := range buckets {
		log.Println(bucket)
	}
	return buckets
}

func getFileUrl(ctx context.Context, objName string, timeout time.Duration) (preSignedURL *url.URL, err error) {
	// Set request parameters for content-disposition.
	ctx = context.Background()
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+objName)

	if preSignedURL, err = Client.PresignedGetObject(ctx, bucketName, objName, timeout, reqParams); err != nil {
		fmt.Println(err)
	}
	return
}

func uploadFile(ctx context.Context, filePath, contextType, objectName string) {
	object, err := Client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contextType})
	if err != nil {
		log.Println("upload Error：", err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, object.Size)
}

func getCoverImpl(videoPath, snapshotPath string, frameNum int) (snapshotName string, err error) {
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf).
		Run()

	if err != nil {
		log.Fatal("get cover error：", err)
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("get cover error：", err)
		return "", err
	}

	err = imaging.Save(img, snapshotPath)
	if err != nil {
		log.Fatal("get cover error：", err)
		return "", err
	}

	return snapshotName, nil
}

func getCover(ctx context.Context, fileName string) (string, error) {
	coverPath := fileName + ".jpg"
	_, err := getCoverImpl(fileName, coverPath, 0)
	if err != nil {
		return "", err
	}
	return coverPath, nil
}
