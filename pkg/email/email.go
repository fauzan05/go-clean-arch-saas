package email

import (
	"bytes"
	"crypto/tls"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/sirupsen/logrus"
)

//go:embed templates/*.html
var templateFS embed.FS

type EmailService struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	Log      *logrus.Logger
}

func NewEmailService(host string, port int, username, password, from string, log *logrus.Logger) *EmailService {
	return &EmailService{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
		Log:      log,
	}
}

// SendVerificationEmail sends email verification link to user
func (s *EmailService) SendVerificationEmail(toEmail, userName, verificationToken, baseURL string) error {
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", baseURL, verificationToken)

	// Load template from embedded file
	tmpl, err := template.ParseFS(templateFS, "templates/verify_email.html")
	if err != nil {
		s.Log.Errorf("Failed to parse email template: %+v", err)
		return fmt.Errorf("failed to load email template")
	}

	// Prepare template data
	data := struct {
		UserName         string
		VerificationLink string
	}{
		UserName:         userName,
		VerificationLink: verificationLink,
	}

	// Execute template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		s.Log.Errorf("Failed to execute email template: %+v", err)
		return fmt.Errorf("failed to render email template")
	}

	subject := "Verify Your Email Address"
	return s.send(toEmail, subject, body.String())
}

// send sends email using SMTP
func (s *EmailService) send(to, subject, body string) error {
	// If email service not configured, log and return nil (development mode)
	if s.Host == "" || s.Username == "" {
		s.Log.Warnf("Email service not configured. Would send to %s: %s", to, subject)
		s.Log.Infof("Email body:\n%s", body)
		return nil
	}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = s.From
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Setup authentication
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	// Setup TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.Host,
	}

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	// For port 465 (SSL/TLS)
	if s.Port == 465 {
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			s.Log.Errorf("Failed to connect to SMTP server: %+v", err)
			return fmt.Errorf("failed to connect to email server")
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, s.Host)
		if err != nil {
			s.Log.Errorf("Failed to create SMTP client: %+v", err)
			return fmt.Errorf("failed to create email client")
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			s.Log.Errorf("Failed to authenticate: %+v", err)
			return fmt.Errorf("failed to authenticate with email server")
		}

		if err = client.Mail(s.From); err != nil {
			s.Log.Errorf("Failed to set sender: %+v", err)
			return fmt.Errorf("failed to set sender")
		}

		if err = client.Rcpt(to); err != nil {
			s.Log.Errorf("Failed to set recipient: %+v", err)
			return fmt.Errorf("failed to set recipient")
		}

		w, err := client.Data()
		if err != nil {
			s.Log.Errorf("Failed to get data writer: %+v", err)
			return fmt.Errorf("failed to prepare email")
		}

		_, err = w.Write([]byte(message))
		if err != nil {
			s.Log.Errorf("Failed to write message: %+v", err)
			return fmt.Errorf("failed to write email")
		}

		err = w.Close()
		if err != nil {
			s.Log.Errorf("Failed to close writer: %+v", err)
			return fmt.Errorf("failed to send email")
		}

		return client.Quit()
	}

	// For port 587 (STARTTLS) or 25
	return smtp.SendMail(addr, auth, s.From, []string{to}, []byte(message))
}
