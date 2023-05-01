package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/soldiii/diplom/internal/model"
	"github.com/soldiii/diplom/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	NumberOfAttemptsForCodeEntry = 3
	CodeEntryTimeInMinutes       = 5
)

type AuthService struct {
	repo         repository.Authorization
	emailService *EmailService
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo, emailService: NewEmailService()}
}

func (s *AuthService) SetAndCheckUser(user *model.UserCode) error {

	if err := s.repo.CheckForEmail(user.Email); err == nil {
		err = errors.New("почта уже используется")
		return err
	}

	password, err := GeneratePasswordHash(user.EncryptedPassword)
	if err != nil {
		return err
	}
	user.EncryptedPassword = password
	user.RegistrationDateTime = time.Now()
	user.AttemptNumber = 1
	return nil
}

func (s *AuthService) CreateAgent(user *model.UserCode) (int, error) {
	if err := s.SetAndCheckUser(user); err != nil {
		return 0, err
	}

	sup_id, err := strconv.Atoi(user.SupervisorID)
	if err != nil {
		return 0, err
	}

	if err := s.repo.CheckForSupervisor(sup_id); err != nil {
		err = errors.New("супервайзер с таким id не существует")
		return 0, err
	}

	code := s.emailService.GenerateCode()
	user.Code = code

	sup_id, err = strconv.Atoi(user.SupervisorID)
	if err != nil {
		return 0, err
	}

	emailSupervisor, err := s.repo.GetSupervisorEmailFromID(sup_id)
	if err != nil {
		return 0, err
	}

	s.emailService.SendEmailToSupervisor(emailSupervisor, user.Email, user.Name, user.Surname, user.Patronymic, code)

	return s.repo.CreateUserTempTable(user)
}

func (s *AuthService) CreateSupervisor(user *model.UserCode) (int, error) {
	var mainSuper model.Supervisor
	if err := s.SetAndCheckUser(user); err != nil {
		return 0, err
	}

	initials := user.Surname + " " + string([]rune(user.Name)[0:1]) + ". " + string([]rune(user.Patronymic)[0:1]) + "."
	user.SupervisorInitials = initials

	flag, err := s.repo.IsDBHaveMainSupervisor()
	if err != nil {
		return 0, err
	}

	if !flag {
		var userT model.User
		mainSuper.SupervisorInitials = initials
		userT.Name = user.Name
		userT.Surname = user.Surname
		userT.Patronymic = user.Patronymic
		userT.Email = user.Email
		userT.EncryptedPassword = user.EncryptedPassword
		userT.RegistrationDateTime = user.RegistrationDateTime
		userT.Role = user.Role

		return s.repo.CreateMainSupervisor(&userT, &mainSuper)
	}

	emailMainSupervisor, err := s.repo.GetEmailOfMainSupervisor()
	if err != nil {
		return 0, err
	}

	user.SupervisorID = "0"
	code := s.emailService.GenerateCode()
	user.Code = code
	s.emailService.SendEmailToMainSupervisor(emailMainSupervisor, user.Email, user.Name, user.Surname, user.Patronymic, code)

	return s.repo.CreateUserTempTable(user)
}

func (s *AuthService) CompareRegistrationCodes(email string, code string) (int, error) {
	timeFlag, err := s.IsTimeExpired(email)
	if err != nil {
		return 0, err
	}
	if timeFlag {
		err = errors.New("время регистрации истекло")
		return 0, err
	}
	result, err := s.repo.IsRegistrationCodeValid(email, code)
	if err != nil {
		return 0, err
	}
	if !result {
		attemptNumber, err := s.repo.GetAttemptNumber(email)
		if err != nil {
			return 0, err
		}
		if attemptNumber == NumberOfAttemptsForCodeEntry {
			s.repo.DeleteFromTempTableByEmail(email)
			err = errors.New("превышен лимит количества попыток")
		} else {
			s.repo.IncrementAttemptNumberByEmail(email)
			err = errors.New("неверный код")
		}
		return attemptNumber, err
	}
	s.emailService.SendEmailToRegistratedUser(email)
	return s.repo.MigrateFromTemporaryTable(email)
}

func (s *AuthService) IsTimeExpired(email string) (bool, error) {
	regDateTime, err := s.repo.GetRegistrationTime(email)
	if err != nil {
		return false, err
	}
	regTimeInMinutes := regDateTime.Unix() / 60
	timeNowInMinutes := time.Now().Unix() / 60
	if timeNowInMinutes-regTimeInMinutes > CodeEntryTimeInMinutes {
		s.repo.DeleteFromTempTableByEmail(email)
		return true, nil
	}
	return false, nil
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
