package main

import (
	"context"
	"github.com/li1553770945/personal-file-service/biz/infra/container"
	"github.com/li1553770945/personal-file-service/kitex_gen/file"
)

// FileServiceImpl FilesServiceImpl implements the last service interface defined in the IDL.
type FileServiceImpl struct{}

// UploadFile implements the FileServiceImpl interface.
func (s *FileServiceImpl) UploadFile(ctx context.Context, req *file.UploadFileReq) (resp *file.UploadFileResp, err error) {
	APP := container.GetGlobalContainer()
	resp, err = APP.FileService.UploadFile(ctx, req)
	return
}

// DownloadFile implements the FileServiceImpl interface.
func (s *FileServiceImpl) DownloadFile(ctx context.Context, req *file.DownloadFileReq) (resp *file.DownloadFileResp, err error) {
	APP := container.GetGlobalContainer()
	resp, err = APP.FileService.DownloadFile(ctx, req)
	return
}

// DeleteFile implements the FileServiceImpl interface.
func (s *FileServiceImpl) DeleteFile(ctx context.Context, req *file.DeleteFileReq) (resp *file.DeleteFileResp, err error) {
	APP := container.GetGlobalContainer()
	resp, err = APP.FileService.DeleteFile(ctx, req)
	return
}

// FileInfo implements the FileServiceImpl interface.
func (s *FileServiceImpl) FileInfo(ctx context.Context, req *file.FileInfoReq) (resp *file.FileInfoResp, err error) {
	APP := container.GetGlobalContainer()
	resp, err = APP.FileService.FileInfo(ctx, req)
	return
}
