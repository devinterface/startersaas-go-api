package services

import (
	"devinterface.com/startersaas-go-api/models"
	"github.com/Kamva/mgm/v3"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
)

// AccountService struct
type AccountService struct{ BaseService }

func (accountService *AccountService) getCollection() (collection *mgm.Collection) {
	coll := mgm.CollectionByName("account")
	return coll
}

// ByID function
func (accountService *AccountService) ByID(id interface{}) (account *models.Account, err error) {
	account = &models.Account{}
	err = accountService.getCollection().FindByID(id, account)
	return account, err
}

// OneBy function
func (accountService *AccountService) OneBy(q bson.M) (account *models.Account, err error) {
	account = &models.Account{}
	err = accountService.getCollection().First(q, account)
	return account, err
}

// Update function
func (accountService *AccountService) Update(id interface{}, params interface{}) (updatedAccount *models.Account, err error) {
	account := &models.Account{}
	err = accountService.getCollection().FindByID(id, account)
	if err != nil {
		return account, err
	}
	err = mapstructure.Decode(params, &account)
	if err != nil {
		return account, err
	}
	err = accountService.getCollection().Update(account)
	return account, err
}

// Create function
func (accountService *AccountService) Create(params interface{}) (createdAccount *models.Account, err error) {
	account := &models.Account{}
	err = mapstructure.Decode(params, &account)
	if err != nil {
		return account, err
	}
	err = accountService.getCollection().Create(account)
	return account, err
}

// Delete function
func (accountService *AccountService) Delete(id interface{}) (success bool, err error) {
	account := &models.Account{}
	err = accountService.getCollection().FindByID(id, account)
	if err != nil {
		return false, err
	}
	err = accountService.getCollection().Delete(account)
	return err == nil, err
}

// FindBy function
func (accountService *AccountService) FindBy(q bson.M) (accounts []models.Account, err error) {
	accounts = []models.Account{}
	err = accountService.getCollection().SimpleFind(&accounts, q)
	return accounts, err
}
