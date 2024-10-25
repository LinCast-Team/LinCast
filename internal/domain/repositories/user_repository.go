package repositories

import (
	"lincast/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetById(id uuid.UUID) (*entities.User, error)
	Create(user *entities.User) error
	Update(user *entities.User) error
	Delete(id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db,
	}
}

func (ur *userRepository) GetById(id uuid.UUID) (*entities.User, error) {
	var u entities.User

	if err := ur.db.First(&u, id).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (ur *userRepository) Create(user *entities.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) Update(user *entities.User) error {
	if err := ur.db.Save(user).Error; err != nil {
		return nil
	}

	return nil
}

func (ur *userRepository) Delete(id uuid.UUID) error {
	if err := ur.db.Delete(&entities.User{ID: id}).Error; err != nil {
		return err
	}

	return nil
}
