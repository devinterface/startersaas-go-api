package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"
	"time"

	"devinterface.com/startersaas-go-api/models"
	"github.com/Kamva/mgm/v3"
	"github.com/osteele/liquid"
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

	err = SendMail(os.Getenv("MAILER"), os.Getenv("DEFAULT_EMAIL_FROM"), emailModel.Subject, result, []string{user.Email})
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

	err = SendMail(os.Getenv("MAILER"), os.Getenv("DEFAULT_EMAIL_FROM"), emailModel.Subject, result, []string{user.Email})
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
	err = SendMail(os.Getenv("MAILER"), os.Getenv("DEFAULT_EMAIL_FROM"), emailModel.Subject, result, []string{user.Email})
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

	err = SendMail(os.Getenv("MAILER"), os.Getenv("DEFAULT_EMAIL_FROM"), subject, result, []string{email})
	if err != nil {
		return false, err
	}

	return true, err
}

// SendMail function
func SendMail(addr, from, subject, body string, to []string) error {
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(r.Replace(from)); err != nil {
		return err
	}
	for i := range to {
		to[i] = r.Replace(to[i])
		if err = c.Rcpt(to[i]); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func (emailService *EmailService) StoreEmails() (err error) {
	for _, code := range []string{"activate", "activationLink", "forgotPassword", "notification"} {
		storedEmail, _ := loadEmail(code, "en")
		if storedEmail.Code == "" {
			template, _ := ioutil.ReadFile(fmt.Sprintf("./emails/%s.email.liquid", code))
			email := &models.Email{}
			email.Body = string(template)
			email.Code = code
			email.Lang = "en"
			email.Subject = fmt.Sprintf("[Starter SaaS] %s", code)
			err = emailService.getCollection().Create(email)
		}
	}
	return err
}
