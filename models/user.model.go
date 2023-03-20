package models

import (
	"time"

	"github.com/Kamva/mgm/v3"
	"github.com/devinterface/structomap"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInner struct {
	ID    primitive.ObjectID `json:"id" bson:"id"`
	Email string             `json:"email" bson:"email"`
}

// User struct
type User struct {
	mgm.DefaultModel     `bson:",inline"`
	Name                 string             `json:"name" bson:"name"`
	Surname              string             `json:"surname" bson:"surname"`
	Email                string             `json:"email" bson:"email"`
	Language             string             `json:"language" bson:"language"`
	Password             string             `json:"password" bson:"password"`
	Role                 string             `json:"role" bson:"role"`
	ConfirmationToken    string             `json:"confirmationToken" bson:"confirmationToken"`
	PasswordResetToken   string             `json:"passwordResetToken" bson:"passwordResetToken"`
	PasswordResetExpires time.Time          `json:"passwordResetExpires" bson:"passwordResetExpires"`
	Sso                  string             `json:"sso" bson:"sso"`
	Active               bool               `json:"active" bson:"active"`
	AccountOwner         bool               `json:"accountOwner" bson:"accountOwner"`
	AccountID            primitive.ObjectID `json:"accountId" bson:"accountId"`
	Teams                []TeamInner        `json:"teams" bson:"teams"`
}

const (
	// UserRole function
	UserRole = "user"
	// AdminRole function
	AdminRole = "admin"
	// SuperAdminRole function
	SuperAdminRole = "superadmin"
)

// UserSerializer function
type UserSerializer struct {
	*structomap.Base
}

// ShowUserSerializer function
func ShowUserSerializer() *UserSerializer {
	u := &UserSerializer{structomap.New()}
	u.UseCamelCase().Pick("ID", "Name", "Surname", "Email", "Language", "Role", "Active", "AccountID", "AccountOwner", "Teams", "CreatedAt", "UpdatedAt")
	return u
}

func (user *User) ToUserInner() UserInner {
	userInner := UserInner{}
	userInner.ID = user.ID
	userInner.Email = user.Email
	return userInner
}
