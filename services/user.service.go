package services

import (
	"errors"

	"devinterface.com/startersaas-go-api/models"
	"github.com/Kamva/mgm/v3"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserService struct
type UserService struct{ BaseService }

func (userService *UserService) getCollection() (collection *mgm.Collection) {
	coll := mgm.CollectionByName("user")
	return coll
}

// OneBy function
func (userService *UserService) OneBy(q bson.M) (user *models.User, err error) {
	user = &models.User{}
	err = userService.getCollection().First(q, user)
	return user, err
}

// ByID function
func (accountService *UserService) ByID(id interface{}, accountID primitive.ObjectID) (loadedUser *models.User, err error) {
	user := &models.User{}
	primitiveID, err := user.PrepareID(id)
	if err != nil {
		return user, err
	}
	err = userService.getCollection().First(bson.M{"_id": primitiveID, "accountId": accountID}, user)
	if err != nil {
		return user, err
	}
	return user, err
}

// Update function
func (userService *UserService) Update(id interface{}, accountID primitive.ObjectID, params interface{}) (updatedUser *models.User, err error) {
	user := &models.User{}
	primitiveID, err := user.PrepareID(id)
	if err != nil {
		return user, err
	}
	err = userService.getCollection().First(bson.M{"_id": primitiveID, "accountId": accountID}, user)
	if err != nil {
		return user, err
	}
	err = mapstructure.Decode(params, &user)
	if err != nil {
		return nil, err
	}
	err = userService.getCollection().Update(user)
	return user, err
}

// Create function
func (userService *UserService) Create(params interface{}, accountID primitive.ObjectID) (createdUser *models.User, err error) {
	user := &models.User{}
	err = mapstructure.Decode(params, &user)

	// check if email unique
	existentUser := &models.User{}
	err = userService.getCollection().First(bson.M{"email": user.Email}, existentUser)
	if existentUser.ID != primitive.NilObjectID {
		return existentUser, errors.New("email is invalid or already taken")
	}
	// if user.Password != "" {
	// 	hash, _ := hashPassword(user.Password)
	// 	user.Password = hash
	// } else {
	// 	defer emailService.SendForgotPasswordEmail(bson.M{"email": user.Email})
	// }

	defer emailService.SendActiveEmail(bson.M{"email": user.Email})
	ssoUUID, _ := uuid.NewRandom()
	user.Sso = ssoUUID.String()
	user.AccountID = accountID
	user.Active = true
	user.Teams = []models.TeamInner{}
	err = userService.getCollection().Create(user)
	return user, err
}

// UpdatePassword function
func (userService *UserService) UpdatePassword(id interface{}, password string) (updatedUser *models.User, err error) {
	user := &models.User{}
	err = userService.getCollection().FindByID(id, user)
	if err != nil {
		return user, err
	}
	hash, _ := hashPassword(password)
	user.Password = hash
	err = userService.getCollection().Update(user)
	return user, err
}

// FindBy function
func (userService *UserService) FindBy(q bson.M) (users []models.User, err error) {
	users = []models.User{}
	err = userService.getCollection().SimpleFind(&users, q)
	return users, err
}

// Delete function
func (userService *UserService) Delete(id interface{}, accountID primitive.ObjectID) (success bool, err error) {
	user := &models.User{}
	primitiveID, err := user.PrepareID(id)
	if err != nil {
		return false, err
	}
	err = userService.getCollection().First(bson.M{"_id": primitiveID, "accountId": accountID}, user)
	if err != nil {
		return false, err
	}
	if user.AccountOwner {
		return false, errors.New("Can't delete account owner")
	}
	err = userService.getCollection().Delete(user)
	return err == nil, err
}
