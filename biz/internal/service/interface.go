package file

import (
	"context"
	"github.com/li1553770945/personal-file-service/biz/internal/repo"
	"github.com/li1553770945/personal-file-service/kitex_gen/file"
)

type FileService struct {
	Repo repo.IRepository
}

type IFileService interface {
	UploadFile(ctx context.Context, req *file.UploadFileReq) (*file.UploadFileResp, error)
	DownloadFile(ctx context.Context, req *file.DownloadFileReq) (*file.DownloadFileResp, error)
	DeleteFile(ctx context.Context, req *file.DeleteFileReq) (*file.DeleteFileResp, error)
}

func NewFileService(repo repo.IRepository) IFileService {
	return &FileService{
		Repo: repo,
	}
}
