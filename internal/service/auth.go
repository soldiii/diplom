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
	CODE_ENTRY_TIME_IN_MINUTES int64 = 5
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
	ACCESS_TOKEN_TTL  = 30 * time.Minute
	REFRESH_TOKEN_TTL = 120 * time.Hour
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

	supID, err := strconv.Atoi(user.SupervisorID)
	if err != nil {
		return 0, err
	}

	if err := s.infoRepo.CheckForSupervisor(supID); err != nil {
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
	sup_id, err := s.infoRepo.GetSupervisorIDByAgentID(id)
	if err != nil {
		return err
	}
	if _, err := s.reportRepo.SetReport(id); err != nil {
		return err
	}
	if err != nil {
		return err
	}

	if _, err := s.planRepo.SetPlan(sup_id, id); err != nil {
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
		"sub":      strconv.Itoa(user.ID),
		"exp":      time.Now().Add(REFRESH_TOKEN_TTL).Unix(),
		"iat":      time.Now().Unix(),
		"userID":   user.ID,
		"userRole": user.Role,
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
		err := errors.New("неверный email или пароль")
		return nil, err
	}
	encryptedPassword, err := s.repo.GetPassword(email)
	if err != nil {
		return nil, err
	}
	if !CheckPasswordHash(password, encryptedPassword) {
		err := errors.New("неверный email или пароль")
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

	rtokenFlag, err := s.repo.IsUserHaveRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}
	if rtokenFlag {
		_, err = s.repo.UpdateRefreshToken(user.ID, tokens.RefreshToken)
		if err != nil {
			return nil, err
		}
		return tokens, nil
	}
	_, err = s.repo.PostRefreshToken(user.ID, tokens.RefreshToken)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID   int    `json:"userID"`
	UserRole string `json:"userRole"`
}

func (s *AuthService) ParseToken(tokenstring string, flag bool) (*TokenClaims, error) {
	var secretKey string
	if flag {
		secretKey = os.Getenv("ACCESS_SECRET_KEY")
	} else {
		secretKey = os.Getenv("REFRESH_SECRET_KEY")
	}

	token, err := jwt.ParseWithClaims(tokenstring, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("token claims are not of type *TokenClaims")
	}

	return claims, nil
}

func (s *AuthService) IsTokenExpired(ExpiredAt *jwt.NumericDate) bool {
	return ExpiredAt.Unix() < time.Now().Unix()
}

func (s *AuthService) CompareRefreshTokens(token string, id int) (bool, error) {
	return s.repo.CompareRefreshTokens(token, id)
}

func CreateTokensByRefresh(role string, id int, accessSecretKey, refreshSecretKey []byte) (*Token, error) {
	accessClaims := jwt.MapClaims{
		"sub":      strconv.Itoa(id),
		"exp":      time.Now().Add(ACCESS_TOKEN_TTL).Unix(),
		"iat":      time.Now().Unix(),
		"userID":   id,
		"userRole": role,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(accessSecretKey)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"sub":      strconv.Itoa(id),
		"exp":      time.Now().Add(REFRESH_TOKEN_TTL).Unix(),
		"iat":      time.Now().Unix(),
		"userID":   id,
		"userRole": role,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(refreshSecretKey)
	if err != nil {
		return nil, err
	}
	return &Token{AccessToken: accessString, RefreshToken: refreshString}, nil
}

func (s *AuthService) GenerateTokensByRefresh(role string, id int) (*Token, error) {
	accessSecretKey := os.Getenv("ACCESS_SECRET_KEY")
	refreshSecretKey := os.Getenv("REFRESH_SECRET_KEY")

	tokens, err := CreateTokensByRefresh(role, id, []byte(accessSecretKey), []byte(refreshSecretKey))
	if err != nil {
		return nil, err
	}
	_, err = s.repo.UpdateRefreshToken(id, tokens.RefreshToken)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *AuthService) RefreshTokens(refresh_token string, role string, id int) (*Token, error) {
	refreshFlag, err := s.CompareRefreshTokens(refresh_token, id)
	if err != nil {
		return nil, err
	}
	if refreshFlag {
		tokens, err := s.GenerateTokensByRefresh(role, id)
		if err != nil {
			return nil, err
		}
		return tokens, nil
	}
	err = errors.New("неверный токен")
	return nil, err
}
