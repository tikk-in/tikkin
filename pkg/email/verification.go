package email

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/smtp"
	"os"
	"strconv"
	"text/template"
	"tikkin/pkg/config"
	"tikkin/pkg/model"
)

type EmailHandler struct {
	config *config.Config
}

type VerificationEmailProps struct {
	Subject          string
	SiteUrl          string
	SiteName         string
	VerificationCode string
	VerificationUrl  string
}

func NewEmailHandler(cfg *config.Config) EmailHandler {
	return EmailHandler{config: cfg}
}

func (e *EmailHandler) SendVerificationEmail(user model.User) error {
	auth := smtp.PlainAuth("", e.config.Email.SMTP.Username, e.config.Email.SMTP.Password, e.config.Email.SMTP.Host)
	to := []string{user.Email}

	hostPort := e.config.Email.SMTP.Host + ":" + strconv.Itoa(e.config.Email.SMTP.Port)
	//verificationURL := e.config.Site.URL + "/api/v1/auth/verify/" + user.VerificationToken

	// Working Directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get working directory")
		return err
	}
	tmpl, err := template.ParseFiles(wd + "/template/verification.html")
	if err != nil {
		return err
	}

	props := VerificationEmailProps{
		Subject:          fmt.Sprintf("Verify your %s email", e.config.Site.Name),
		SiteUrl:          e.config.Site.URL,
		SiteName:         e.config.Site.Name,
		VerificationCode: *user.VerificationToken,
		VerificationUrl:  e.config.Site.URL + "/api/v1/users/verify/" + *user.VerificationToken,
	}

	buff := new(bytes.Buffer)
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	buff.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", props.Subject, mimeHeaders)))

	err = tmpl.Execute(buff, props)
	if err != nil {
		return err
	}

	err = smtp.SendMail(hostPort, auth, e.config.Email.SMTP.From, to, buff.Bytes())
	if err != nil {
		return err
	}
	log.Info().Str("email", user.Email).Msg("Verification email sent")
	return nil
}
