package repo

import "github.com/li1553770945/personal-file-service/biz/internal/domain"

func (Repo *Repository) RemoveFile(fileKey string) error {
	err := Repo.DB.Delete(&domain.FileEntity{}, "file_key = ?", fileKey).Error
	if err != nil {
		return err
	}
	return nil
}

func (Repo *Repository) SaveFile(entity *domain.FileEntity) error {
	if entity.ID == 0 {
		err := Repo.DB.Create(&entity).Error
		return err
	} else {
		err := Repo.DB.Save(&entity).Error
		return err
	}
}
func (Repo *Repository) GetFile(fileKey string) (*domain.FileEntity, error) {
	var entity domain.FileEntity
	err := Repo.DB.Where("file_key = ?", fileKey).Find(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
