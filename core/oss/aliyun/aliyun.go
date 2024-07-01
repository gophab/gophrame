package aliyun

import (
	"mime/multipart"
	"regexp"
	"time"

	"github.com/gophab/gophrame/core/oss/aliyun/config"
	"github.com/gophab/gophrame/core/util"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliyunOSS struct {
	Client *oss.Client
	Bucket *oss.Bucket
}

func CreateAliyunOSS() (*AliyunOSS, error) {
	endpoint := config.Setting.Endpoint
	if config.Setting.UseCname {
		endpoint = config.Setting.BucketUrl
	}
	client, err := oss.New(
		endpoint,
		config.Setting.AccessKeyId,
		config.Setting.AccessKeySecret,
		oss.UseCname(config.Setting.UseCname),
		oss.AuthVersion(oss.AuthV4))
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(config.Setting.Bucket)
	if err != nil {
		return nil, err
	}
	return &AliyunOSS{
		Client: client,
		Bucket: bucket,
	}, nil
}

func (s *AliyunOSS) Upload(file *multipart.FileHeader, prefix string) (string, string, error) {
	// 读取本地文件。
	f, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer f.Close() // 创建文件 defer 关闭
	// 上传阿里云路径 文件名格式 自己可以改 建议保证唯一性
	// yunFileTmpPath := filepath.Join("uploads", time.Now().Format("2006-01-02")) + "/" + file.Filename
	_re := regexp.MustCompile(("^/"))
	re_ := regexp.MustCompile(("/.$"))
	yunFileTmpPath := _re.ReplaceAllString(
		re_.ReplaceAllString(re_.ReplaceAllString(config.Setting.Path, "")+"/"+_re.ReplaceAllString(prefix, ""), "")+
			"/"+
			time.Now().Format("YYYYMMDDhhmmss")+
			util.GenerateRandomString(6),
		"")

	// 上传文件流。
	err = s.Bucket.PutObject(yunFileTmpPath, f)
	if err != nil {
		return "", "", err
	}

	if len(config.Setting.BucketUrl) == 0 {
		return re_.ReplaceAllString(config.Setting.Endpoint, "") + "/" + yunFileTmpPath, yunFileTmpPath, nil
	} else {
		return re_.ReplaceAllString(config.Setting.BucketUrl, "") + "/" + yunFileTmpPath, yunFileTmpPath, nil
	}
}
