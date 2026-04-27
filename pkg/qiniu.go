package pkg

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/qf0129/backups/conf"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/objects"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)

var credential *credentials.Credentials

func InitQiniu() {
	credential = credentials.NewCredentials(conf.Conf.Qiniu.AccessKey, conf.Conf.Qiniu.SecretKey)
}

func UploadToQiniu(localPath string, remotePath string) error {
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{Credentials: credential},
	})
	if conf.Conf.Qiniu.BucketDir != "" {
		remotePath = filepath.Join(conf.Conf.Qiniu.BucketDir, remotePath)
	}
	return uploadManager.UploadFile(context.Background(), localPath, &uploader.ObjectOptions{
		BucketName: conf.Conf.Qiniu.Bucket,
		ObjectName: &remotePath,
		FileName:   filepath.Base(remotePath),
		OnUploadingProgress: func(progress *uploader.UploadingProgress) {
			if progress.TotalSize == 0 {
				return
			}
			fmt.Printf("\rUploading %.2f%%", float64(progress.Uploaded*100)/float64(progress.TotalSize))
			if progress.Uploaded == progress.TotalSize {
				println()
			}
		},
	}, nil)
}

func QiniuListBucket(bucketName, filePrefix string) objects.Lister {
	manager := objects.NewObjectsManager(&objects.ObjectsManagerOptions{
		Options: http_client.Options{Credentials: credential},
	})
	options := &objects.ListObjectsOptions{Prefix: filePrefix}
	bucket := manager.Bucket(bucketName)
	return bucket.List(context.Background(), options)
}

func QiniuDelete(bucketName, fileName string) error {
	manager := objects.NewObjectsManager(&objects.ObjectsManagerOptions{
		Options: http_client.Options{Credentials: credential},
	})
	bucket := manager.Bucket(bucketName)
	return bucket.Object(fileName).Delete().Call(context.Background())
}
