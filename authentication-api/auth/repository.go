package auth

import (
	"context"

	"gorm.io/gorm"
)

type TursoRepo struct {
	DB *gorm.DB
}

func (r *TursoRepo) Insert(ctx context.Context, user User) error {
	result := r.DB.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *TursoRepo) GetByID(ctx context.Context, id uint64) (User, error) {
	var user User
	result := r.DB.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func (r *TursoRepo) FindByUsername(ctx context.Context, username string) (User, error) {
	var user User
	result := r.DB.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}