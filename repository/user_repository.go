package repository

import (
	"errors"
	"fmt"

	"promptscroll/config"
	"promptscroll/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Create(auth *model.Auth) error
	GetByUserID(userID uint64) (*model.Auth, error)
	DeleteByUserID(userID uint64) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository() AuthRepository {
	return &authRepository{
		db: config.GetDB(),
	}
}

func (r *authRepository) Create(auth *model.Auth) error {
	if err := r.db.Create(auth).Error; err != nil {
		return err
	}

	return nil
}

func (r *authRepository) GetByUserID(userID uint64) (*model.Auth, error) {
	auth := &model.Auth{}
	if err := r.db.Where("user_id = ?", userID).First(auth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get auth by user id error: %v", err)
	}

	return auth, nil
}

func (r *authRepository) DeleteByUserID(userID uint64) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&model.Auth{}).Error; err != nil {
		return fmt.Errorf("delete auth by user id error: %v", err)
	}

	return nil
}
