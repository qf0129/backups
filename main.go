package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/qf0129/backups/conf"
	"github.com/qf0129/backups/pkg"
	"github.com/qiniu/go-sdk/v7/storagev2/objects"
)

var tmpDir = "tmp"
var configFlag = flag.String("c", "conf.json", "Config file path")
var versionFlag = flag.Bool("v", false, "Show version")

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println("v0.0.1")
		return
	}
	configFile := strings.TrimSpace(*configFlag)
	if configFile == "" {
		configFile = "conf.json"
	}
	if err := conf.LoadConfig(configFile); err != nil {
		slog.Error(err.Error())
		return
	}
	os.MkdirAll(tmpDir, 0755)
	pkg.InitQiniu()
	for _, path := range conf.Conf.Paths {
		if err := doBackup(path); err != nil {
			if !conf.Conf.SkipFailed {
				return
			}
		}
	}
}

func doBackup(path string) error {
	slog.Info(">>> Backup " + path)
	zipPath := pkg.GetLocalZipPath(path)
	remotePath := pkg.GetRemoteFilePath(path)

	// 打包
	err := pkg.PathToZip(path, zipPath)
	if err != nil {
		slog.Error("pathToZip error:" + err.Error())
		return err
	} else {
		slog.Info("Packaged success")
	}
	defer os.Remove(zipPath)

	// 上传
	err = pkg.UploadToQiniu(zipPath, remotePath)
	if err != nil {
		slog.Error("uploadToQiniu error: " + err.Error())
		return err
	} else {
		slog.Info("Upload success")
	}
	if conf.Conf.RotateByDay && conf.Conf.RotateDays > 0 {
		cleanRotateFiles()
	}
	return nil
}

func cleanRotateFiles() error {
	iter := pkg.QiniuListBucket(conf.Conf.Qiniu.Bucket, conf.Conf.Qiniu.BucketDir)
	defer iter.Close()
	var objectInfo objects.ObjectDetails
	for iter.Next(&objectInfo) {
		timeStr := strings.Split(strings.TrimPrefix(objectInfo.Name, conf.Conf.Qiniu.BucketDir+"/"), "/")[0]
		t, err := time.ParseInLocation(pkg.DirTimeLayout, timeStr, time.Local)
		if err != nil {
			slog.Error("time parse error: " + err.Error())
			continue
		}
		if time.Since(t) > time.Duration(conf.Conf.RotateDays)*24*time.Hour {
			slog.Info("Delete rotate file: " + objectInfo.Name)
			if err := pkg.QiniuDelete(conf.Conf.Qiniu.Bucket, objectInfo.Name); err != nil {
				slog.Error("delete error: " + err.Error())
			}
		}
	}
	if err := iter.Error(); err != nil {
		return err
	}
	slog.Info("Clean success")
	return nil
}
