package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 解压
func Unzip(zipFile string, destDir string) ([]string, error) {
	zipReader, err := zip.OpenReader(zipFile)
	var paths []string
	if err != nil {
		return []string{}, err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		if strings.Contains(f.Name, "..") {
			return []string{}, fmt.Errorf("%s 文件名不合法", f.Name)
		}
		fpath := filepath.Join(destDir, f.Name)
		paths = append(paths, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return []string{}, err
			}

			inFile, err := f.Open()
			if err != nil {
				return []string{}, err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return []string{}, err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return []string{}, err
			}
		}
	}
	return paths, nil
}

func isDir(path string) (bool, error) {
	// 尝试获取文件信息
	fileInfo, err := os.Stat(path)

	if err != nil {
		// 如果发生错误（例如文件不存在），返回错误
		return false, err
	}
	// 使用IsDir方法检查是否是目录
	return fileInfo.IsDir(), nil
}

// func dirFiles(dir string) ([]string, error) {
// 	var files []string
// 	//方法一
// 	var walkFunc = func(path string, info os.FileInfo, err error) error {
// 		if !info.IsDir() {
// 			files = append(files, path)
// 		}
// 		return nil
// 	}
// 	err := filepath.Walk(dir, walkFunc)
// 	return files, err
// }

func dirFiles(dir string) ([]string, error) {
	var files []string
	results, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range results {
		files = append(files, dir+"/"+file.Name())
	}
	return files, nil
}

func zipDir(zipWriter *zip.Writer, file string, oldForm, newForm string) error {
	files, err := dirFiles(file)
	if err == nil {
		for _, file := range files {
			zipFile(zipWriter, file, oldForm, newForm)
		}
	}
	return err
}

func zipFile(zipWriter *zip.Writer, file string, oldForm, newForm string) error {
	if b, err := isDir(file); err != nil {
		return err
	} else if b {
		return zipDir(zipWriter, file, oldForm, newForm)
	}

	zipFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	// 获取file的基础信息
	info, err := zipFile.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// 使用上面的FileInforHeader() 就可以把文件保存的路径替换成我们自己想要的了，如下面
	header.Name = strings.Replace(file, oldForm, newForm, -1)

	// 优化压缩
	// 更多参考see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	if _, err = io.Copy(writer, zipFile); err != nil {
		return err
	}
	return nil
}

func ZipFiles(filename string, files []string, oldForm, newForm string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = newZipFile.Close()
	}()

	zipWriter := zip.NewWriter(newZipFile)
	defer func() {
		_ = zipWriter.Close()
	}()

	// 把files添加到zip中
	for _, file := range files {
		zipFile(zipWriter, file, oldForm, newForm)
	}
	return nil
}
