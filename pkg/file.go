package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/qf0129/backups/conf"
)

func GetLocalZipPath(path string) string {
	filePrefix := strings.ReplaceAll(filepath.Join("backups", strings.TrimSpace(path)), string(os.PathSeparator), "_")
	return filepath.Join("tmp", fmt.Sprintf("%s_%s.zip", filePrefix, time.Now().Format("200601021504")))
}

func GetRemoteFilePath(path string) string {
	filePrefix := strings.Trim(path, string(os.PathSeparator))
	remotePath := fmt.Sprintf("%s_%s.zip", filePrefix, time.Now().Format("200601021504"))
	if conf.Conf.RotateByDay {
		return filepath.Join(time.Now().Format("20060102"), remotePath)
	}
	return remotePath
}
