package repository

import (
	"gorm.io/gorm"
)

type BaseRepository struct {
	model *gorm.DB
}

func NewRepository(model *gorm.DB) (Repository, error) {
	return &BaseRepository{
		model: model,
	}, nil
}

func (b *BaseRepository) Find(id int64) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	err := b.model.Where("id = ?", id).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *BaseRepository) Select(query map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	err := b.model.Where(query).First(&result).Error
	if err != nil {
		return nil
	}
	return result
}

func (b *BaseRepository) Update(query map[string]interface{}, data interface{}) (count int64, err error) {
	save := b.model.Save(data)
	if save.Error != nil {
		return 0, save.Error
	}
	return save.RowsAffected, nil
}
func (b *BaseRepository) Add(data map[string]interface{}) (int64, error) {
	err := b.model.Create(data).Error
	if err != nil {
		return 0, err
	}
	return b.model.RowsAffected, nil
}
