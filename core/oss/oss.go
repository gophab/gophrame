package oss

import "mime/multipart"

type OSS interface {
	Upload(file *multipart.FileHeader, prefix string) (string, string, error)
	UploadLocal(fileName string, prefix string) (string, string, error)
}
