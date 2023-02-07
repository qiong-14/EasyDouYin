package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/qiong-14/EasyDouYin/utils"
	"log"
	"net/url"
	"time"
)

//SERVER_IP=124.221.190.158
//MINIO_PORT=9000
//MINIO_ENDPOINT=$SERVER_IP:$MINIO_PORT
//MINIO_ACCESS_KEY=iR1BteCtymOsJ7rY
//MINIO_SECRET_KEY=SMDeqEnnKdOvRt1hgbHzE5PmG4uhAPbf
//MINIO_BUCKET=videos-and-covers
var (
	endpoint        = utils.GetEnvByKey("MINIO_ENDPOINT")
	accessKeyID     = utils.GetEnvByKey("MINIO_ACCESS_KEY")
	secretAccessKey = utils.GetEnvByKey("MINIO_SECRET_KEY")
	bucketName      = utils.GetEnvByKey("MINIO_BUCKET")
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

func ListAllBuckets(ctx context.Context) []minio.BucketInfo {
	buckets, _ := Client.ListBuckets(context.Background())
	for _, bucket := range buckets {
		log.Println(bucket)
	}
	return buckets
}

func UploadFile(ctx context.Context, filePath, contextType, objectName string) {
	object, err := Client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contextType})
	if err != nil {
		log.Println("upload Error：", err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, object.Size)
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

func getFileUrl(ctx context.Context, objName string, timeout time.Duration) (presignedURL *url.URL, err error) {
	// Set request parameters for content-disposition.
	ctx = context.Background()
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+objName)

	if presignedURL, err = Client.PresignedGetObject(ctx, bucketName, objName, timeout, reqParams); err != nil {
		fmt.Println(err)
	}
	return
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
