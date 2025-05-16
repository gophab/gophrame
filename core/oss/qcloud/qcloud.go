package qcloud

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gophab/gophrame/core/oss/qcloud/config"
	"github.com/gophab/gophrame/core/util"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type QcloudOSS struct {
	Client *cos.Client
}

func CreateQcloudOSS() (*QcloudOSS, error) {
	urlStr, err := url.Parse("https://" + config.Setting.Bucket + ".cos." + config.Setting.Region + ".myqcloud.com")
	if err != nil {
		return nil, err
	}

	baseURL := &cos.BaseURL{BucketURL: urlStr}
	client := cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.Setting.AppId,
			SecretKey: config.Setting.AppKey,
		},
	})

	return &QcloudOSS{
		Client: client,
	}, nil
}

var (
	_re = regexp.MustCompile(("^/"))
	re_ = regexp.MustCompile(("/.$"))
)

func (s *QcloudOSS) Upload(file *multipart.FileHeader, prefix string) (string, string, error) {
	f, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer f.Close() // 创建文件 defer 关闭

	ext := filepath.Ext(file.Filename)
	yunFileTmpPath := _re.ReplaceAllString(
		re_.ReplaceAllString(re_.ReplaceAllString(config.Setting.Path, "")+"/"+_re.ReplaceAllString(prefix, ""), "")+
			"/"+
			time.Now().Format("20060102/150405_")+util.GenerateRandomString(6)+ext,
		"")

	_, err = s.Client.Object.Put(context.Background(), yunFileTmpPath, f, nil)
	if err != nil {
		return "", "", err
	}

	return re_.ReplaceAllString(config.Setting.BucketUrl, "") + "/" + yunFileTmpPath, yunFileTmpPath, nil
}

func (s *QcloudOSS) UploadLocal(fileName string, prefix string) (string, string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", "", err
	}
	defer f.Close() // 创建文件 defer 关闭

	ext := filepath.Ext(fileName)
	yunFileTmpPath := _re.ReplaceAllString(
		re_.ReplaceAllString(re_.ReplaceAllString(config.Setting.Path, "")+"/"+_re.ReplaceAllString(prefix, ""), "")+
			"/"+
			time.Now().Format("20060102/150405_")+util.GenerateRandomString(6)+ext,
		"")

	_, err = s.Client.Object.Put(context.Background(), yunFileTmpPath, f, nil)
	if err != nil {
		return "", "", err
	}

	return re_.ReplaceAllString(config.Setting.BucketUrl, "") + "/" + yunFileTmpPath, yunFileTmpPath, nil
}
