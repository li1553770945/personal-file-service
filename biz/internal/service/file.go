package file

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/li1553770945/personal-file-service/biz/constant"
	"github.com/li1553770945/personal-file-service/biz/internal/assembler"
	"github.com/li1553770945/personal-file-service/kitex_gen/base"
	"github.com/li1553770945/personal-file-service/kitex_gen/file"
	"math/rand"
	"unicode"
)

func (s *FileService) UploadFile(ctx context.Context, req *file.UploadFileReq) (resp *file.UploadFileResp, err error) {
	if req.Key == nil {
		for { // 防止已有相同key
			Key := s.generateKey()
			req.Key = &Key
			entity, err := s.Repo.GetFile(Key)
			if err != nil {
				resp = &file.UploadFileResp{
					BaseResp: &base.BaseResp{
						Code:    constant.SystemError,
						Message: "数据库访问错误",
					},
				}
				klog.CtxInfof(ctx, "查询已有文件时数据库访问错误：%v", err.Error())
				return resp, nil
			}
			// 如果没有
			if entity.ID == 0 {
				break
			}
		}
	} else if !s.IsAlphanumeric(*req.Key) { // 用户提交了key，检查是否合法
		resp := &file.UploadFileResp{
			BaseResp: &base.BaseResp{
				Code:    constant.InvalidInput,
				Message: "非法key，只能使用数字或者大小写字母",
			},
		}
		return resp, nil
	}
	entity := assembler.FileReqToEntity(req)
	entity.DownloadCount = 0
	err = s.Repo.SaveFile(entity)
	if err != nil {
		resp = &file.UploadFileResp{
			BaseResp: &base.BaseResp{
				Code:    constant.SystemError,
				Message: "数据库访问错误",
			},
		}
		klog.CtxInfof(ctx, "数据库访问错误：%v", err.Error())
		return resp, nil
	}
	ak, sk, err := s.getAkSk(ctx)
	if err != nil {
		resp = &file.UploadFileResp{
			BaseResp: &base.BaseResp{
				Code:    constant.SystemError,
				Message: "生成临时sk失败",
			},
		}
		klog.CtxInfof(ctx, "生成临时sk失败:%v", err.Error())
		return resp, nil
	}
	resp = &file.UploadFileResp{
		BaseResp: &base.BaseResp{
			Code: constant.Success,
		},
		Ak:  ak,
		Sk:  sk,
		Key: entity.Key,
	}
	return resp, nil
}
func (s *FileService) DownloadFile(ctx context.Context, req *file.DownloadFileReq) (resp *file.DownloadFileResp, err error) {
	entity, err := s.Repo.GetFile(req.GetKey())
	if err != nil {
		resp = &file.DownloadFileResp{
			BaseResp: &base.BaseResp{
				Code:    constant.SystemError,
				Message: "数据库访问错误",
			},
		}
		klog.CtxInfof(ctx, "数据库访问错误：%v", err.Error())
		return resp, nil
	}
	if entity.ID == 0 {
		resp = &file.DownloadFileResp{
			BaseResp: &base.BaseResp{
				Code:    constant.NotFound,
				Message: "未找到对应的文件，请检查文件key是否正确",
			},
		}
		return resp, nil
	}
	ak, sk, err := s.getAkSk(ctx)
	if err != nil {
		resp = &file.DownloadFileResp{
			BaseResp: &base.BaseResp{
				Code:    constant.SystemError,
				Message: "生成临时sk失败",
			},
		}
		klog.CtxInfof(ctx, "生成临时sk失败:%v", err.Error())
		return resp, nil
	}
	resp = &file.DownloadFileResp{
		BaseResp: &base.BaseResp{
			Code: constant.Success,
		},
		Ak:      ak,
		Sk:      sk,
		OssPath: entity.OSSPath,
		Name:    entity.Name,
	}
	return resp, nil
}

func (s *FileService) DeleteFile(ctx context.Context, req *file.DeleteFileReq) (resp *file.DeleteFileResp, err error) {
	entity, err := s.Repo.GetFile(req.GetKey())
	if err != nil {
		resp = &file.DeleteFileResp{
			BaseResp: &base.BaseResp{
				Code:    constant.SystemError,
				Message: "数据库访问错误",
			},
		}
		klog.CtxInfof(ctx, "数据库访问错误：%v", err.Error())
		return resp, nil
	}
	if entity.ID == 0 {
		resp = &file.DeleteFileResp{
			BaseResp: &base.BaseResp{
				Code:    constant.NotFound,
				Message: "未找到对应的文件，请检查文件key是否正确",
			},
		}
		return resp, nil
	}
	err = s.Repo.RemoveFile(req.Key)
	if err != nil {
		resp = &file.DeleteFileResp{
			BaseResp: &base.BaseResp{
				Code:    constant.SystemError,
				Message: "删除文件时数据库访问错误",
			},
		}
		return resp, nil
	}
	resp = &file.DeleteFileResp{
		BaseResp: &base.BaseResp{
			Code: constant.Success,
		},
	}
	return resp, nil
}

func (s *FileService) getAkSk(ctx context.Context) (ak string, sk string, err error) {
	return "ak", "sk", nil
}

func (s *FileService) generateKey() (key string) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 4)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *FileService) IsAlphanumeric(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
