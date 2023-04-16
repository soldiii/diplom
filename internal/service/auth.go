package service

import (
	"time"

	"github.com/soldiii/diplom/internal/model"
	"github.com/soldiii/diplom/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user *model.User) (int, error) {
	password, err := GeneratePasswordHash(user.EncryptedPassword)
	if err != nil {
		return 0, err
	}
	user.EncryptedPassword = password
	user.RegistrationDateTime = time.Now()

	return s.repo.CreateUser(user)
}

func GeneratePasswordHash(password string) (string, error) {
	saltedBytes := []byte(password)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hash := string(hashedBytes[:])
	return hash, nil
}
