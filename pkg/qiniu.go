package pkg

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/qf0129/backups/conf"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)

func UploadToQiniu(localPath string, remotePath string) error {
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: credentials.NewCredentials(conf.Conf.Qiniu.AccessKey, conf.Conf.Qiniu.SecretKey),
		},
	})
	if conf.Conf.Qiniu.BucketDir != "" {
		remotePath = filepath.Join(conf.Conf.Qiniu.BucketDir, remotePath)
	}
	return uploadManager.UploadFile(context.Background(), localPath, &uploader.ObjectOptions{
		BucketName: conf.Conf.Qiniu.Bucket,
		ObjectName: &remotePath,
		FileName:   filepath.Base(remotePath),
		OnUploadingProgress: func(progress *uploader.UploadingProgress) {
			fmt.Printf("\rUploading %.2f%%", float64(progress.Uploaded*100)/float64(progress.TotalSize))
		},
	}, nil)
}
