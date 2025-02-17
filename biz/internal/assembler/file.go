package assembler

import (
	"github.com/li1553770945/personal-file-service/biz/internal/domain"
	"github.com/li1553770945/personal-file-service/kitex_gen/file"
)

func FileReqToEntity(req *file.UploadFileReq) (entity *domain.FileEntity) {
	return &domain.FileEntity{
		Name:        req.GetName(),
		Key:         req.GetKey(),
		MaxDownload: req.MaxDownload,
	}
}
