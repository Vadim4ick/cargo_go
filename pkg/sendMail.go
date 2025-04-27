package pkg

import (
	"fmt"
	"net/smtp"
	"test-project/config"
)

func SendEmail(to, inviteLink string) error {
	// Настройки SMTP сервера
	smtpHost := config.Envs.SMTP_HOST
	smtpPort := config.Envs.SMTP_PORT
	username := config.Envs.SMTP_USER
	password := config.Envs.SMTP_PASS

	// Адрес отправителя
	from := username

	subject := "Приглашение на регистрацию"

	htmlBody := `<p>Вас приглашают зарегистрироваться.</p><p><a href="` + inviteLink + `">Нажмите здесь</a>, чтобы пройти регистрацию. Время жизни 5 минут.</p>`

	message := []byte(fmt.Sprintf(
		"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s",
		subject,
		htmlBody,
	))

	// Авторизация на сервере
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// Отправляем письмо
	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{to},
		message,
	)

	return err
}
