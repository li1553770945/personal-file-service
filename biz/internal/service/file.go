package file

import (
	"context"
	"github.com/li1553770945/personal-file-service/biz/internal/assembler"
	"github.com/li1553770945/personal-file-service/kitex_gen/file"
	"unicode"
	"math/rand"

)

func (s *FileService) UploadFile(ctx context.Context, req *file.UploadFileReq) (*file.UploadFileResp, error) {
	if req.Key == nil{
		Key := s.generateKey(ctx)
		req.Key = &Key
	} else {
		if !s.IsAlphanumeric(*req.Key) {
			return nil,
		}
	}
	entity := assembler.FileReqToEntity(req)
	entity.DownloadCount = 0
	err := s.Repo.SaveFile(entity)
	if err != nil {
		return nil, err
	}
	ak, sk, err := s.getAkSk(ctx)
	if err != nil {
		return nil, err
	}
	resp := &file.UploadFileResp{
		Ak: ak,
		Sk: sk,
		Key:
	}
	return nil, nil
}
func (s *FileService) DownloadFile(ctx context.Context, req *file.DownloadFileReq) (*file.DownloadFileResp, error) {
	return nil, nil
}

func (s *FileService) DeleteFile(ctx context.Context, req *file.UploadFileReq) (*file.UploadFileResp, error) {

}

func (s *FileService) getAkSk(ctx context.Context) (ak string, sk string,err error) {

}

func (s *FileService) generateKey(ctx context.Context) (key string) {

}

func  (s *FileService)RandSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func  (s *FileService)IsAlphanumeric(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}


