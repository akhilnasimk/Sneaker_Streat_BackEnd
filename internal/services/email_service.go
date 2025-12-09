package services

import (
	"fmt"
	"net/smtp"

	"github.com/akhilnasimk/SS_backend/internal/config"
)

type EmailService interface {
	SendEmail(to string, subject string, body string) error
	SendMailOTP(to string, otp string) error
}

type emailService struct{}

func NewEmailService() EmailService {
	return &emailService{}
}


func (s *emailService) SendEmail(to, subject, body string) error {
	from := config.AppConfig.SMTPEmail
	password := config.AppConfig.SMTPPass
	host := config.AppConfig.SMTPHost
	port := config.AppConfig.SMTPPort

	fmt.Println("DEBUG: SMTP Config ->", "Host:", host, "Port:", port, "From:", from)

	auth := smtp.PlainAuth("", from, password, host)

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		body

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", host, port),
		auth,
		from,
		[]string{to},
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}


// SendOTP sends an OTP email
func (s *emailService) SendMailOTP(to, otp string) error {
	// You can customize the subject and HTML body here
	subject := "Your OTP Code"
	body := fmt.Sprintf(`
		<html>
		<body>
			<p>Hello,</p>
			<p>Your OTP code is: <b>%s</b></p>
			<p>Please do not share this code with anyone.</p>
		</body>
		</html>
	`, otp)

	// Use the generic SendEmail function
	return s.SendEmail(to, subject, body)
}
