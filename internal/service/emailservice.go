package service

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	Host        string
	Port        string
	AccName     string
	AccPassword string
}

func NewEmailService() *EmailService {
	return &EmailService{
		Host:        os.Getenv("EMAIL_HOST"),
		Port:        os.Getenv("EMAIL_PORT"),
		AccName:     os.Getenv("EMAIL_ACCOUNT_NAME"),
		AccPassword: os.Getenv("EMAIL_PASSWORD"),
	}
}

func (e *EmailService) SendEmailToMainSupervisor(to string, supervisor_email string, name string, surname string, patronymic string, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.AccName)
	m.SetHeader("To", to)

	m.SetHeader("Subject", "Регистрация супервайзера")

	text := fmt.Sprintf("Необходимо ваше подтверждение для регистрации супервайзера: %s %s %s, c адресом электронной почты: %s.\nСообщите ему следующий код, если вы подтверждаете регистрацию: %s.", surname, name, patronymic, supervisor_email, code)
	m.SetBody("text/plain", text)

	port, err := strconv.Atoi(e.Port)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(e.Host, port, e.AccName, e.AccPassword)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (e *EmailService) SendEmailToCommonSupervisor(to string, agent_email string, name string, surname string, patronymic string, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.AccName)
	m.SetHeader("To", to)

	m.SetHeader("Subject", "Регистрация агента")

	text := fmt.Sprintf("Необходимо ваше подтверждение для регистрации агента: %s %s %s, c адресом электронной почты: %s.\nСообщите ему следующий код, если вы подтверждаете регистрацию: %s.", surname, name, patronymic, agent_email, code)
	m.SetBody("text/plain", text)

	port, err := strconv.Atoi(e.Port)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(e.Host, port, e.AccName, e.AccPassword)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (e *EmailService) SendEmailToSupervisor(to string, email string, name string, surname string, patronymic string, code string, role string) error {
	switch role {
	case "supervisor", "Supervisor":
		if err := e.SendEmailToMainSupervisor(to, email, name, surname, patronymic, code); err != nil {
			return err
		}
	case "agent", "Agent":
		if err := e.SendEmailToCommonSupervisor(to, email, name, surname, patronymic, code); err != nil {
			return err
		}
	}
	return nil
}

func (e *EmailService) SendEmailToRegistratedUser(to string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", e.AccName)
	m.SetHeader("To", to)

	m.SetHeader("Subject", "Успешная регистрация")

	m.SetBody("text/plain", "Вы успешно зарегистрировались в приложении.")

	port, err := strconv.Atoi(e.Port)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(e.Host, port, e.AccName, e.AccPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (e *EmailService) GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000
	return strconv.Itoa(code)
}
