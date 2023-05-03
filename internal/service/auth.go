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
	NUMBER_ATTEMPTS                  = 3
	CODE_ENTRY_TIME_IN_MINUTES int64 = 5
)

type AuthService struct {
	repo         repository.Authorization
	emailService *EmailService
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo, emailService: NewEmailService()}
}

func (s *AuthService) SetUser(user *model.UserCode) (bool, error) {
	if err := s.repo.CheckForEmail(user.Email); err == nil {
		err = errors.New("почта уже используется")
		return false, err
	}

	flag, err := s.repo.IsTempTableHaveUser(user.Email)
	if err != nil {
		return false, err
	}
	if flag {
		timeFlag, err := s.IsTimeExpired(user.Email)
		if err != nil {
			return false, err
		}
		if timeFlag {
			err = errors.New("время регистрации истекло")
			if err := s.ClearTempTableFromUsersWithExpiredTime(); err != nil {
				return false, err
			}
			return false, err
		}
		attNumber, err := s.repo.GetAttemptNumberByEmail(user.Email)
		if err != nil {
			return false, err
		}
		code, err := s.repo.GetCodeByEmail(user.Email)
		if err != nil {
			return false, err
		}
		dateTime, err := s.repo.GetRegistrationTimeByEmail(user.Email)
		if err != nil {
			return false, err
		}
		user.RegistrationDateTime = dateTime
		user.Code = code
		user.AttemptNumber = attNumber
		password, err := GeneratePasswordHash(user.EncryptedPassword)
		if err != nil {
			return false, err
		}
		user.EncryptedPassword = password
		s.repo.DeleteFromTempTableByEmail(user.Email)
		if err := s.ClearTempTableFromUsersWithExpiredTime(); err != nil {
			return false, err
		}
	} else {
		code := s.emailService.GenerateCode()
		user.Code = code
		user.AttemptNumber = 1
		user.RegistrationDateTime = time.Now()
		password, err := GeneratePasswordHash(user.EncryptedPassword)
		if err != nil {
			return false, err
		}
		user.EncryptedPassword = password
		if err := s.ClearTempTableFromUsersWithExpiredTime(); err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (s *AuthService) CreateAgent(user *model.UserCode) (int, error) {
	user.SupervisorInitials = "NOT SUPERVISOR"
	emailFlag, err := s.SetUser(user)
	if err != nil {
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

	if emailFlag {
		emailSupervisor, err := s.repo.GetSupervisorEmailFromID(sup_id)
		if err != nil {
			return 0, err
		}
		s.emailService.SendEmailToSupervisor(emailSupervisor, user.Email, user.Name, user.Surname, user.Patronymic, user.Code, user.Role)
	}

	return s.repo.CreateUserTempTable(user)
}

func (s *AuthService) CreateSupervisor(user *model.UserCode) (int, error) {
	user.SupervisorID = "0"
	var mainSuper model.Supervisor
	emailFlag, err := s.SetUser(user)
	if err != nil {
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
	if emailFlag {
		emailMainSupervisor, err := s.repo.GetEmailOfMainSupervisor()
		if err != nil {
			return 0, err
		}
		s.emailService.SendEmailToSupervisor(emailMainSupervisor, user.Email, user.Name, user.Surname, user.Patronymic, user.Code, user.Role)
	}

	return s.repo.CreateUserTempTable(user)
}

func (s *AuthService) CompareRegistrationCodes(email string, code string) (int, error) {
	flag, err := s.repo.IsTempTableHaveUser(email)
	if err != nil {
		return 0, err
	}
	if flag {
		timeFlag, err := s.IsTimeExpired(email)
		if err != nil {
			return 0, err
		}
		if timeFlag {
			err = errors.New("время регистрации истекло")
			return 0, err
		}
	} else {
		err = errors.New("время регистрации истекло")
		return 0, err
	}
	result, err := s.repo.IsRegistrationCodeValid(email, code)
	if err != nil {
		return 0, err
	}
	if !result {
		attemptNumber, err := s.repo.GetAttemptNumberByEmail(email)
		if err != nil {
			return 0, err
		}
		if attemptNumber == NUMBER_ATTEMPTS {
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

func TransformTime(dateTime time.Time) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		return time.Now(), err
	}
	dateTime = dateTime.In(loc)
	return dateTime, nil
}

func (s *AuthService) IsTimeExpired(email string) (bool, error) {
	regDateTime, err := s.repo.GetRegistrationTimeByEmail(email)
	if err != nil {
		return false, err
	}
	regTime, err := TransformTime(regDateTime)
	if err != nil {
		return false, err
	}
	timeNow, err := TransformTime(time.Now())
	if err != nil {
		return false, err
	}
	entryTime := CODE_ENTRY_TIME_IN_MINUTES
	timeNowInMinutes := timeNow.Unix() / 60
	regTimeInMinutes := (regTime.Unix() - 25200) / 60
	if timeNowInMinutes-regTimeInMinutes > entryTime {
		s.repo.DeleteFromTempTableByEmail(email)
		return true, nil
	}
	return false, nil
}

func (s *AuthService) ClearTempTableFromUsersWithExpiredTime() error {
	EntryTime := CODE_ENTRY_TIME_IN_MINUTES
	emails, err := s.repo.GetUsersEmailsWithExpiredTime(time.Now(), EntryTime)
	if err != nil {
		return err
	}

	for _, email := range emails {
		s.repo.DeleteFromTempTableByEmail(email)
	}
	return nil
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
