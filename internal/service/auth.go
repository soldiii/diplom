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

func SetUser(user *model.User) error {
	password, err := GeneratePasswordHash(user.EncryptedPassword)
	if err != nil {
		return err
	}
	user.EncryptedPassword = password
	user.RegistrationDateTime = time.Now()
	return nil
}

func (s *AuthService) CreateAgent(user *model.User, agent *model.Agent) (int, error) {
	if err := SetUser(user); err != nil {
		return 0, err
	}
	agent.SupervisorID = user.SupervisorID

	return s.repo.CreateAgent(user, agent)
}

func (s *AuthService) CreateSupervisor(user *model.User, supervisor *model.Supervisor) (int, error) {
	if err := SetUser(user); err != nil {
		return 0, err
	}
	initials := user.Surname + " " + string(user.Name[0]) + ". " + string(user.Patronymic[0]) + "."
	supervisor.SupervisorInitials = initials
	user.SupervisorID = 0

	return s.repo.CreateSupervisor(user, supervisor)
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

func (s *AuthService) GetAllSupervisors() ([]*model.Supervisor, error) {
	return s.repo.GetAllSupervisors()
}

func (s *AuthService) GenerateToken(email, password string) (string, error) {
	return "", nil
}
