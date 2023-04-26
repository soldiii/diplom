package service

import (
	"gopkg.in/gomail.v2"
)

func SendEmailAboutRegistration(to string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "t3sttesttest1234@yandex.ru")
	m.SetHeader("To", to)

	m.SetHeader("Subject", "Успешная регистрация")

	m.SetBody("text/plain", "Вы успешно зарегистрировались")

	d := gomail.NewDialer("smtp.yandex.ru", 465, "t3sttesttest1234@yandex.ru", "Narodnaya4719!")

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
