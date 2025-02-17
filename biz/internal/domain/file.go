package domain

import (
	"github.com/li1553770945/personal-file-service/biz/internal/do"
)

type FileEntity struct {
	do.BaseModel
	Name          string
	Key           string
	OSSPath       string
	MaxDownload   int32
	DownloadCount int32
}
