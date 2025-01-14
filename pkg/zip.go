package pkg

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func PathToZip(targetPath, zipPath string) error {
	// 创建压缩文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 创建zip写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历目录
	return filepath.Walk(targetPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(targetPath, filePath)
		if err != nil {
			return err
		}

		// 如果相对路径为空，则使用文件名
		if relPath == "." {
			relPath = info.Name()
		}

		// 跳过目录本身
		if info.IsDir() {
			return nil
		}

		// 创建zip文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath
		header.Method = zip.Deflate

		// 创建压缩文件
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// 打开源文件
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// 复制文件内容到压缩文件
		_, err = io.Copy(writer, file)
		return err
	})
}
