package services

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"devinterface.com/startersaas-go-api/models"
	"github.com/Kamva/mgm/v3"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/kataras/i18n"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// AuthService struct
type AuthService struct{ BaseService }

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Login function
func (authService *AuthService) Login(email string, cleanPassword string, refreshToken bool) (response map[string]string, err error) {
	user := &models.User{}
	coll := mgm.CollectionByName("user")
	coll.First(bson.M{"email": strings.ToLower(email), "active": true}, user)

	// check password if we are not refreshing JWT token
	if !refreshToken {
		match := checkPasswordHash(cleanPassword, user.Password)
		if !match {
			return nil, errors.New("username or password invalid")
		}
	}

	message := map[string]string{
		"success": "true",
		"message": "Enjoy your token!",
		"token":   generateToken(*user),
	}

	return message, err
}

// Signup function
func (authService *AuthService) Signup(params map[string]interface{}, signupWithActivate bool) (response map[string]string, err error) {
	account := &models.Account{}
	accountColl := mgm.CollectionByName("account")
	account.Subdomain = strcase.ToKebab(params["subdomain"].(string))

	// check if subdomain unique
	existentAccount := &models.Account{}
	accountColl.First(bson.M{"subdomain": account.Subdomain}, existentAccount)
	if existentAccount.ID != primitive.NilObjectID {
		return nil, errors.New("subdomain is invalid or already taken")
	}

	user := &models.User{}
	userColl := mgm.CollectionByName("user")
	user.Email = strings.TrimSpace(strings.ToLower(params["email"].(string)))

	// check if email unique
	existentUser := &models.User{}
	userColl.First(bson.M{"email": user.Email}, existentUser)
	if existentUser.ID != primitive.NilObjectID {
		return nil, errors.New("email is invalid or already taken")
	}

	// create account
	trialDays, _ := strconv.Atoi(os.Getenv("TRIAL_DAYS"))
	account.TrialPeriodEndsAt = time.Now().AddDate(0, 0, trialDays)
	account.PrivacyAccepted = params["privacyAccepted"].(bool)
	account.MarketingAccepted = params["marketingAccepted"].(bool)
	account.PlanType = os.Getenv("STARTER_PLAN_TYPE")
	err = accountColl.Create(account)
	if err != nil {
		return nil, err
	}

	// create user
	user.Role = models.AdminRole
	user.AccountOwner = true
	user.Active = signupWithActivate
	if params["language"] != nil {
		user.Language = params["language"].(string)
	} else {
		user.Language = os.Getenv("LOCALE")
	}
	ssoUUID, _ := uuid.NewRandom()
	user.Sso = ssoUUID.String()
	hash, _ := hashPassword(params["password"].(string))
	user.Password = hash
	user.AccountID = account.ID
	err = userColl.Create(user)
	if err != nil {
		return nil, err
	}

	go emailService.SendNotificationEmail(os.Getenv("NOTIFIED_ADMIN_EMAIL"), i18n.Tr("en", "authService.signup.subject"), i18n.Tr("en", "authService.signup.messageAdmin", map[string]string{"Subdomain": account.Subdomain, "Email": user.Email}), os.Getenv("LOCALE"))

	if !signupWithActivate {
		go emailService.SendActivationEmail(bson.M{"_id": user.ID})
		return nil, err
	} else {
		message := map[string]string{
			"token": generateToken(*user),
		}
		return message, err
	}
}

// Activate function
func (authService *AuthService) Activate(token string, email string) (success bool, err error) {
	user := &models.User{}
	var userService = UserService{}
	err = userService.getCollection().First(bson.M{"active": false, "confirmationToken": token, "email": email}, user)
	if err != nil {
		return false, err
	}
	user.Active = true
	user.ConfirmationToken = ""
	err = userService.getCollection().Update(user)
	if err != nil {
		return false, err
	}

	go emailService.SendActiveEmail(bson.M{"_id": user.ID})
	return true, err
}

// ResetPassword function
func (authService *AuthService) ResetPassword(token string, password string, email string) (success bool, err error) {
	user := &models.User{}
	var userService = UserService{}
	err = userService.getCollection().First(bson.M{"passwordResetToken": token, "email": email}, user)
	if err != nil {
		return false, err
	}
	if user.PasswordResetExpires.Before(time.Now()) {
		return false, errors.New("PasswordResetToken is expired")
	}
	hash, _ := hashPassword(password)
	user.Password = hash
	user.PasswordResetToken = ""
	err = userService.getCollection().Update(user)
	return true, err
}

// Sso function
func (authService *AuthService) Sso(sso string) (response map[string]string, err error) {
	user := &models.User{}
	coll := mgm.CollectionByName("user")
	coll.First(bson.M{"active": true, "sso": sso}, user)

	message := map[string]string{
		"success": "true",
		"message": "Enjoy your token!",
		"token":   generateToken(*user),
	}

	return message, err
}

func generateToken(user models.User) (tokenString string) {
	expire, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().AddDate(0, 0, expire).Unix(),
	})

	key := os.Getenv("JWT_SECRET")
	tokenString, _ = token.SignedString([]byte(key))
	return tokenString
}
