package gorm

import (
	"context"
	"errors"

	"goster/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, err := GetTransaction(ctx); err == nil {
		return tx
	}
	return r.db.WithContext(ctx)
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.getDB(ctx).Create(user).Error
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.getDB(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.getDB(ctx).Model(&domain.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
