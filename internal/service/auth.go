package service

import (
	"errors"
	"strconv"
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
	if user.SupervisorID == "" {
		err := errors.New("необходимо ввести id супервайзера")
		return 0, err
	}
	if err := SetUser(user); err != nil {
		return 0, err
	}

	sup_id, err := strconv.Atoi(user.SupervisorID)
	if err != nil {
		return 0, err
	}
	agent.SupervisorID = sup_id

	return s.repo.CreateAgent(user, agent)
}

func (s *AuthService) CreateSupervisor(user *model.User, supervisor *model.Supervisor) (int, error) {
	if user.SupervisorID != "" {
		err := errors.New("вводить id супервайзера при регистрации не нужно")
		return 0, err
	}
	if err := SetUser(user); err != nil {
		return 0, err
	}
	initials := user.Surname + " " + string([]rune(user.Name)[0:1]) + ". " + string([]rune(user.Patronymic)[0:1]) + "."
	supervisor.SupervisorInitials = initials

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
