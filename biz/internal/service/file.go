package file

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/google/uuid"
	"github.com/li1553770945/personal-file-service/biz/constant"
	"github.com/li1553770945/personal-file-service/biz/internal/assembler"
	"github.com/li1553770945/personal-file-service/kitex_gen/base"
	"github.com/li1553770945/personal-file-service/kitex_gen/file"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"
	"unicode"
)

func (s *FileService) UploadFile(ctx context.Context, req *file.UploadFileReq) (resp *file.UploadFileResp, err error) {
	if req.Key == nil || *req.Key == "" {
		for { // 防止已有相同key
			Key := s.generateKey(constant.DEFAULT_KEY_LENGTH)
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
	} else {
		if !s.isAlphanumeric(*req.Key) { // 用户提交了key，检查是否合法
			resp := &file.UploadFileResp{
				BaseResp: &base.BaseResp{
					Code:    constant.InvalidInput,
					Message: "非法key，只能使用数字或者大小写字母",
				},
			}
			return resp, nil
		}
		entity, err := s.Repo.GetFile(*req.Key)
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
		if entity.ID != 0 {
			resp = &file.UploadFileResp{
				BaseResp: &base.BaseResp{
					Code:    constant.InvalidInput,
					Message: "已存在相同key的文件",
				},
			}
			return resp, nil
		}
	}
	entity := assembler.FileReqToEntity(req)
	entity.DownloadCount = 0
	ossPath := s.generateFilePath(entity.Name)
	entity.OSSPath = ossPath
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

	signedUrl, err := s.GetSignedUrl(ctx, http.MethodPut, entity.OSSPath)
	if err != nil {
		resp = &file.UploadFileResp{
			BaseResp: &base.BaseResp{
				Code:    constant.SystemError,
				Message: "生成预签名url失败",
			},
		}
		klog.CtxInfof(ctx, "生成预签名url失败:%v", err.Error())
		return resp, nil
	}
	resp = &file.UploadFileResp{
		BaseResp: &base.BaseResp{
			Code: constant.Success,
		},
		SignedUrl: signedUrl,
		Key:       entity.Key,
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

	signedUrl, err := s.GetSignedUrl(ctx, http.MethodGet, entity.OSSPath)
	if err != nil {
		resp = &file.DownloadFileResp{
			BaseResp: &base.BaseResp{
				Code:    constant.SystemError,
				Message: "生成预签名url失败",
			},
		}
		klog.CtxInfof(ctx, "生成预签名url失败:%v", err.Error())
		return resp, nil
	}
	if entity.MaxDownload != 0 {
		entity.DownloadCount = entity.DownloadCount + 1
		if entity.DownloadCount == entity.MaxDownload {

			entity.ExpiredTime = gorm.DeletedAt{Time: time.Now(), Valid: true}
			err = s.Repo.SaveFile(entity)
			if err != nil {
				resp = &file.DownloadFileResp{
					BaseResp: &base.BaseResp{
						Code:    constant.SystemError,
						Message: "数据库访问错误",
					},
				}
				klog.CtxInfof(ctx, "保存文件时数据库访问错误：%v", err.Error())
				return nil, err
			}
			klog.CtxInfof(ctx, "文件过期被删除：%v", entity.Key)
			err = s.Repo.RemoveFile(entity.Key)
			if err != nil {
				resp = &file.DownloadFileResp{
					BaseResp: &base.BaseResp{
						Code:    constant.SystemError,
						Message: "删除文件时数据库访问错误",
					},
				}
				return resp, nil
			}
		} else {
			err = s.Repo.SaveFile(entity)
			if err != nil {
				resp = &file.DownloadFileResp{
					BaseResp: &base.BaseResp{
						Code:    constant.SystemError,
						Message: "数据库访问错误",
					},
				}
				klog.CtxInfof(ctx, "保存文件时数据库访问错误：%v", err.Error())
				return nil, err
			}
		}

	}

	resp = &file.DownloadFileResp{
		BaseResp: &base.BaseResp{
			Code: constant.Success,
		},
		SignedUrl: signedUrl,
		Name:      entity.Name,
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
	_, err = s.CosClient.Object.Delete(context.Background(), entity.OSSPath)
	if err != nil {
		resp = &file.DeleteFileResp{
			BaseResp: &base.BaseResp{
				Code:    constant.SystemError,
				Message: "删除oss文件出错",
			},
		}
		klog.CtxErrorf(ctx, "删除oss文件出错:%v", err.Error())
		return resp, nil
	}
	entity.DeleteOnOssTime = gorm.DeletedAt{Time: time.Now(), Valid: true}
	err = s.Repo.SaveFile(entity)
	if err != nil {
		resp = &file.DeleteFileResp{
			BaseResp: &base.BaseResp{
				Code:    constant.SystemError,
				Message: "删除文件时数据库访问错误",
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

func (s *FileService) FileInfo(ctx context.Context, req *file.FileInfoReq) (resp *file.FileInfoResp, err error) {
	entity, err := s.Repo.GetFile(req.GetKey())
	if err != nil {
		resp = &file.FileInfoResp{
			BaseResp: &base.BaseResp{
				Code:    constant.SystemError,
				Message: "数据库访问错误",
			},
		}
		klog.CtxInfof(ctx, "数据库访问错误：%v", err.Error())
		return resp, nil
	}
	if entity.ID == 0 {
		resp = &file.FileInfoResp{
			BaseResp: &base.BaseResp{
				Code:    constant.NotFound,
				Message: "未找到对应的文件，请检查文件key是否正确",
			},
		}
		return resp, nil
	}
	resp = &file.FileInfoResp{
		BaseResp: &base.BaseResp{
			Code: constant.Success,
		},
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return resp, nil
}
func (s *FileService) generateKey(length int32) (key string) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *FileService) isAlphanumeric(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
func (s *FileService) generateFilePath(fileName string) string {
	currentTime := time.Now()

	// 提取当前年月日
	year := currentTime.Year()
	month := currentTime.Month()
	day := currentTime.Day()

	ext := filepath.Ext(fileName)
	newFileName := uuid.New().String() + ext

	// 拼接路径
	filePath := fmt.Sprintf("%d/%d/%d/%s", year, month, day, newFileName)
	return filePath
}
func (s *FileService) GetSignedUrl(ctx context.Context, method, name string) (string, error) {
	signedURL, err := s.CosClient.Object.GetPresignedURL(ctx, method, name, s.Config.CosConfig.Ak, s.Config.CosConfig.Sk, time.Hour, nil)
	if err != nil {
		return "", err
	}
	return signedURL.String(), err
}
