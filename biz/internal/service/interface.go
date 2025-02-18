package file

import (
	"context"
	"github.com/li1553770945/personal-file-service/biz/infra/config"
	"github.com/li1553770945/personal-file-service/biz/internal/repo"
	"github.com/li1553770945/personal-file-service/kitex_gen/file"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type FileService struct {
	Repo      repo.IRepository
	Config    *config.Config
	CosClient *cos.Client
}

type IFileService interface {
	UploadFile(ctx context.Context, req *file.UploadFileReq) (*file.UploadFileResp, error)
	DownloadFile(ctx context.Context, req *file.DownloadFileReq) (*file.DownloadFileResp, error)
	DeleteFile(ctx context.Context, req *file.DeleteFileReq) (*file.DeleteFileResp, error)
}

func NewFileService(repo repo.IRepository, config *config.Config, cosClient *cos.Client) IFileService {
	return &FileService{
		Repo:      repo,
		Config:    config,
		CosClient: cosClient,
	}
}
