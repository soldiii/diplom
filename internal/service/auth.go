package service

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/soldiii/diplom/internal/model"
	"github.com/soldiii/diplom/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	NUMBER_ATTEMPTS                  = 3
	CODE_ENTRY_TIME_IN_MINUTES int64 = 1
)

type AuthService struct {
	repo         repository.Authorization
	infoRepo     repository.Information
	reportRepo   repository.Report
	planRepo     repository.Plan
	emailService *EmailService
}

func NewAuthService(repo repository.Authorization, infoRepo repository.Information, reportRepo repository.Report, planRepo repository.Plan) *AuthService {
	return &AuthService{repo: repo, infoRepo: infoRepo, emailService: NewEmailService(), reportRepo: reportRepo, planRepo: planRepo}
}

const (
	ACCESS_TOKEN_TTL  = 15 * time.Minute
	REFRESH_TOKEN_TTL = 1 * time.Hour
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
			if err := CreateNewUser(user, s); err != nil {
				return false, err
			}
			return true, nil
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
		if err := CreateNewUser(user, s); err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func CreateNewUser(user *model.UserCode, s *AuthService) error {
	code := s.emailService.GenerateCode()
	user.Code = code
	user.AttemptNumber = 1
	user.RegistrationDateTime = time.Now()
	password, err := GeneratePasswordHash(user.EncryptedPassword)
	if err != nil {
		return err
	}
	user.EncryptedPassword = password
	if err := s.ClearTempTableFromUsersWithExpiredTime(); err != nil {
		return err
	}
	return nil
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

	if err := s.infoRepo.CheckForSupervisor(user.SupervisorID); err != nil {
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

func (s *AuthService) SetReportAndPlanTables(id int) error {
	agID := strconv.Itoa(id)
	sup_id, err := s.infoRepo.GetSupervisorIDByAgentID(agID)
	if err != nil {
		return err
	}
	if _, err := s.reportRepo.SetReport(id); err != nil {
		return err
	}
	supID, err := strconv.Atoi(sup_id)
	if err != nil {
		return err
	}

	if _, err := s.planRepo.SetPlan(supID, id); err != nil {
		return err
	}
	return nil
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
	//regTimeInMinutes := (regTime.Unix() - 25200) / 60
	regTimeInMinutes := regTime.Unix() / 60
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

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

func CreateTokens(user *model.User, accessSecretKey, refreshSecretKey []byte) (*Token, error) {
	accessClaims := jwt.MapClaims{
		"sub":      strconv.Itoa(user.ID),
		"exp":      time.Now().Add(ACCESS_TOKEN_TTL).Unix(),
		"iat":      time.Now().Unix(),
		"userID":   user.ID,
		"userRole": user.Role,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(accessSecretKey)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"sub": strconv.Itoa(user.ID),
		"exp": time.Now().Add(REFRESH_TOKEN_TTL).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(refreshSecretKey)
	if err != nil {
		return nil, err
	}
	return &Token{AccessToken: accessString, RefreshToken: refreshString}, nil
}

func (s *AuthService) GenerateTokens(email, password string) (*Token, error) {
	emailFlag, err := s.repo.IsEmailValid(email)
	if err != nil {
		return nil, err
	}
	if !emailFlag {
		err := errors.New("неверный email")
		return nil, err
	}
	encryptedPassword, err := s.repo.GetPassword(email)
	if err != nil {
		return nil, err
	}
	if !CheckPasswordHash(password, encryptedPassword) {
		err := errors.New("неверный пароль")
		return nil, err
	}
	user, err := s.repo.GetUser(email, encryptedPassword)
	if err != nil {
		return nil, err
	}

	accessSecretKey := os.Getenv("ACCESS_SECRET_KEY")
	refreshSecretKey := os.Getenv("REFRESH_SECRET_KEY")

	tokens, err := CreateTokens(user, []byte(accessSecretKey), []byte(refreshSecretKey))
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
