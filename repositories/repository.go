package repositories

import "gorm.io/gorm"

type IRepository[IdType comparable, ModelType any] interface {
	GetById(id IdType) (*ModelType, error)
	Add(model *ModelType) error
	Delete(id IdType) error
	Update(model *ModelType) error
}

type repository[IdType comparable, ModelType any] struct {
	db *gorm.DB
}

func (r *repository[IdType, ModelType]) GetById(id IdType) (*ModelType, error) {
	var p ModelType

	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *repository[IdType, ModelType]) Add(model *ModelType) error {
	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository[IdType, ModelType]) Delete(id IdType) error {
	return nil
}

func (r *repository[IdType, ModelType]) Update(model *ModelType) error {
	return nil
}
