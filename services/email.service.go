package services

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"devinterface.com/startersaas-go-api/models"
	"github.com/Kamva/mgm/v3"
	"github.com/osteele/liquid"
	mail "github.com/xhit/go-simple-mail/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// EmailService struct
type EmailService struct{ BaseService }

func (emailService *EmailService) getCollection() (collection *mgm.Collection) {
	coll := mgm.CollectionByName("email")
	return coll
}

func loadEmail(code string, lang string) (email *models.Email, err error) {
	email = &models.Email{}
	err = emailService.getCollection().First(bson.M{"code": code, "lang": lang}, email)
	if email.Code == "" {
		err = emailService.getCollection().First(bson.M{"code": code, "lang": os.Getenv("LOCALE")}, email)
	}
	return email, err
}

// SendActivationEmail function
func (emailService *EmailService) SendActivationEmail(q bson.M) (success bool, err error) {
	user := &models.User{}
	var userService = UserService{}
	q["active"] = false
	err = userService.getCollection().First(q, user)
	if err != nil {
		return false, err
	}
	p, _ := rand.Prime(rand.Reader, 20)
	user.ConfirmationToken = fmt.Sprintf("%d", p)
	err = userService.getCollection().Update(user)
	if err != nil {
		return false, err
	}

	engine := liquid.NewEngine()
	emailModel, _ := loadEmail("activationLink", user.Language)

	bindings := map[string]interface{}{
		"email":             user.Email,
		"confirmationToken": user.ConfirmationToken,
	}
	result, err := engine.ParseAndRenderString(string(emailModel.Body), bindings)

	err = SendMail(os.Getenv("DEFAULT_EMAIL_FROM"), emailModel.Subject, result, []string{user.Email})
	return true, err
}

// SendForgotPasswordEmail function
func (emailService *EmailService) SendForgotPasswordEmail(q bson.M) (success bool, err error) {
	user := &models.User{}
	var userService = UserService{}
	err = userService.getCollection().First(q, user)
	if err != nil {
		return false, err
	}
	p, _ := rand.Prime(rand.Reader, 20)
	user.PasswordResetToken = fmt.Sprintf("%d", p)
	user.PasswordResetExpires = time.Now().Add(time.Hour * 3)
	err = userService.getCollection().Update(user)
	if err != nil {
		return false, err
	}

	engine := liquid.NewEngine()
	emailModel, _ := loadEmail("forgotPassword", user.Language)
	bindings := map[string]interface{}{
		"email":              user.Email,
		"passwordResetToken": user.PasswordResetToken,
	}
	result, err := engine.ParseAndRenderString(string(emailModel.Body), bindings)

	err = SendMail(os.Getenv("DEFAULT_EMAIL_FROM"), emailModel.Subject, result, []string{user.Email})
	return true, err
}

// SendActiveEmail function
func (emailService *EmailService) SendActiveEmail(q bson.M) (success bool, err error) {
	user := &models.User{}
	var userService = UserService{}
	err = userService.getCollection().First(q, user)
	if err != nil {
		return false, err
	}
	frontendLoginURL := os.Getenv("FRONTEND_LOGIN_URL")

	engine := liquid.NewEngine()
	emailModel, _ := loadEmail("activate", user.Language)
	bindings := map[string]interface{}{
		"email":            user.Email,
		"frontendLoginURL": frontendLoginURL,
	}
	result, err := engine.ParseAndRenderString(string(emailModel.Body), bindings)
	err = SendMail(os.Getenv("DEFAULT_EMAIL_FROM"), emailModel.Subject, result, []string{user.Email})
	return true, err
}

// SendNotificationEmail function
func (emailService *EmailService) SendNotificationEmail(email string, subject string, message string, lang string) (success bool, err error) {
	engine := liquid.NewEngine()

	emailModel, _ := loadEmail("notification", lang)
	frontendLoginURL := os.Getenv("FRONTEND_LOGIN_URL")
	bindings := map[string]interface{}{
		"subject":          subject,
		"email":            email,
		"message":          message,
		"frontendLoginURL": frontendLoginURL,
	}
	result, err := engine.ParseAndRenderString(string(emailModel.Body), bindings)

	err = SendMail(os.Getenv("DEFAULT_EMAIL_FROM"), subject, result, []string{email})
	if err != nil {
		return false, err
	}

	return true, err
}

// SendMail function
func SendMail(from, subject, body string, to []string) error {
	server := mail.NewSMTPClient()
	// SMTP Server
	server.Host = os.Getenv("MAILER_HOST")
	server.Port, _ = strconv.Atoi(os.Getenv("MAILER_PORT"))
	server.Username = os.Getenv("MAILER_USERNAME")
	server.Password = os.Getenv("MAILER_PASSWORD")

	if os.Getenv("MAILER_SSL") != "false" {
		server.Encryption = mail.EncryptionSSLTLS
		server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// SMTP client
	smtpClient, err := server.Connect()

	if err != nil {
		log.Println(err)
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(strings.Join(to, ",")).
		SetSubject(subject)
	email.SetBody(mail.TextHTML, body)

	err = email.Send(smtpClient)
	return err
}

func (emailService *EmailService) StoreEmails() (err error) {
	activate := []string{"activate", "Welcome"}
	activationLink := []string{"activationLink", "Activation link"}
	forgotPassword := []string{"forgotPassword", "Forgot password"}
	notification := []string{"notification", "Notification"}
	for _, email := range [][]string{activate, activationLink, forgotPassword, notification} {
		code := email[0]
		title := email[1]
		storedEmail, _ := loadEmail(code, "en")
		if storedEmail.Code == "" {
			template, _ := ioutil.ReadFile(fmt.Sprintf("./emails/%s.email.liquid", code))
			email := &models.Email{}
			email.Body = string(template)
			email.Code = code
			email.Lang = "en"
			email.Subject = fmt.Sprintf("[StarterSaaS] %s", title)
			err = emailService.getCollection().Create(email)
		}
	}
	return err
}
