package service

import (
	"errors"
	"time"

	"promptscroll/model"
	"promptscroll/repository"
	"promptscroll/util"
)

type AuthService interface {
	Register(user *model.User) (*model.User, error)
	Login(email, password string) (string, error)
}

type authService struct {
	userRepository repository.UserRepository
	tokenSecret    []byte
}

func NewAuthService(userRepository repository.UserRepository, tokenSecret []byte) AuthService {
	return &authService{
		userRepository: userRepository,
		tokenSecret:    tokenSecret,
	}
}

func (s *authService) Register(user *model.User) (*model.User, error) {
	// 檢查帳號是否已被註冊
	if _, err := s.userRepository.GetByEmail(user.Email); err != nil {
		return nil, err
	}

	// 生成密碼的雜湊值
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// 更新密碼為雜湊值
	user.Password = hashedPassword

	// 設置使用者的註冊時間
	user.CreatedAt = time.Now()

	// 在資料庫中建立新的使用者
	if err := s.userRepository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
	// 檢查帳號是否存在
	user, err := s.userRepository.GetByEmail(email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("invalid email or password")
	}

	// 驗證密碼是否正確
	if err := util.VerifyPassword(password, user.Password); err != nil {
		return "", errors.New("invalid email or password")
	}

	// 生成 JWT Token
	token, err := util.GenerateJWT(user.ID, s.tokenSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
