package domain

import (
	"github.com/li1553770945/personal-file-service/biz/internal/do"
	"gorm.io/gorm"
)

type FileEntity struct {
	do.BaseModel
	Name            string
	Key             string
	OSSPath         string
	MaxDownload     int32
	DownloadCount   int32
	ExpiredTime     gorm.DeletedAt
	DeleteOnOssTime gorm.DeletedAt
}
