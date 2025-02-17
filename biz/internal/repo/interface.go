package repo

import (
	"github.com/li1553770945/personal-file-service/biz/internal/domain"
	"gorm.io/gorm"
)

type IRepository interface {
	SaveFile(file *domain.FileEntity) error
	GetFile(fileKey string) (*domain.FileEntity, error)
	RemoveFile(fileKey string) error
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	err := db.AutoMigrate(&domain.FileEntity{})
	if err != nil {
		panic("迁移用户模型失败：" + err.Error())
	}
	return &Repository{
		DB: db,
	}
}
