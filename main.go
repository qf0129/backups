package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/qf0129/backups/conf"
	"github.com/qf0129/backups/pkg"
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
		fmt.Println(err)
		return
	}
	os.MkdirAll(tmpDir, 0755)
	for _, path := range conf.Conf.Paths {
		if err := doBackup(path); err != nil {
			if !conf.Conf.SkipFailed {
				return
			}
		}
	}
	fmt.Println(">>> All done")
}

func doBackup(path string) error {
	fmt.Println(">>> Backup", path)
	zipPath := pkg.GetLocalZipPath(path)
	remotePath := pkg.GetRemoteFilePath(path)

	// 打包
	err := pkg.PathToZip(path, zipPath)
	if err != nil {
		fmt.Println("pathToZip error:", err)
		return err
	} else {
		fmt.Println("zip success")
	}
	defer os.Remove(zipPath)

	// 上传
	err = pkg.UploadToQiniu(zipPath, remotePath)
	if err != nil {
		fmt.Println("\ruploadToQiniu error:", err)
		return err
	} else {
		fmt.Println("\rUpload success          ")
	}
	return nil
}
