package services

import (
	"errors"
	"fmt"

	"devinterface.com/startersaas-go-api/models"
	"github.com/Kamva/mgm/v3"
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TeamService struct
type TeamService struct{ BaseService }

func (teamService *TeamService) getCollection() (collection *mgm.Collection) {
	coll := mgm.CollectionByName("team")
	return coll
}

// ByID function
func (teamService *TeamService) ByID(id interface{}, accountID primitive.ObjectID) (team *models.Team, err error) {
	team = &models.Team{}
	primitiveID, err := team.PrepareID(id)
	if err != nil {
		return team, err
	}
	err = teamService.getCollection().First(bson.M{"_id": primitiveID, "accountId": accountID}, team)
	if err != nil {
		return team, err
	}
	return team, err
}

// OneBy function
func (teamService *TeamService) OneBy(q bson.M) (team *models.Team, err error) {
	team = &models.Team{}
	err = teamService.getCollection().First(q, team)
	return team, err
}

// Create function
func (teamService *TeamService) Create(params interface{}, accountID primitive.ObjectID) (createdTeam *models.Team, err error) {
	team := &models.Team{}
	err = mapstructure.Decode(params, &team)
	if err != nil {
		return team, err
	}
	team.Code = strcase.ToKebab(team.Code)
	team.AccountID = accountID
	team.Users = []models.UserInner{}
	existentTeam := &models.Team{}

	teamService.getCollection().First(bson.M{"code": team.Code, "accountId": accountID}, existentTeam)
	if existentTeam.ID != primitive.NilObjectID {
		return existentTeam, errors.New("code is invalid or already taken")
	}

	err = teamService.getCollection().Create(team)
	return team, err
}

// Update function
func (teamService *TeamService) Update(id interface{}, accountID primitive.ObjectID, params interface{}) (updateTeam *models.Team, err error) {
	team := &models.Team{}
	primitiveID, err := team.PrepareID(id)
	if err != nil {
		return team, err
	}

	err = teamService.getCollection().First(bson.M{"_id": primitiveID, "accountId": accountID}, team)
	if err != nil {
		return team, err
	}

	fmt.Print(params)
	err = mapstructure.Decode(params, &team)

	if err != nil {
		return nil, err
	}

	err = teamService.getCollection().Update(team)
	return team, err
}

// Delete function
func (teamService *TeamService) Delete(id interface{}) (success bool, err error) {
	team := &models.Team{}
	err = teamService.getCollection().FindByID(id, team)
	if err != nil {
		return false, err
	}

	for _, user := range team.Users {
		teamService.RemoveUser(id, team.AccountID, user.ID)
	}

	// update all users
	err = teamService.getCollection().Delete(team)
	return err == nil, err
}

// FindBy function
func (teamService *TeamService) FindBy(q bson.M) (teams []models.Team, err error) {
	teams = []models.Team{}
	err = teamService.getCollection().SimpleFind(&teams, q)
	return teams, err
}

// AddUser function
func (teamService *TeamService) AddUser(id interface{}, accountID primitive.ObjectID, userID interface{}) (updatedTeam *models.Team, err error) {

	team := &models.Team{}
	primitiveID, err := team.PrepareID(id)
	if err != nil {
		return team, err
	}
	err = teamService.getCollection().First(bson.M{"_id": primitiveID, "accountId": accountID}, team)
	if err != nil {
		return team, err
	}

	user := &models.User{}
	userPrimitiveID, err := user.PrepareID(userID)
	if err != nil {
		return team, err
	}
	err = userService.getCollection().First(bson.M{"_id": userPrimitiveID, "accountId": accountID}, user)
	if err != nil {
		return nil, err
	}

	users := team.Users

	users = append(users, user.ToUserInner())
	uniqueUsers := funk.UniqBy(users, func(f models.UserInner) primitive.ObjectID {
		return f.ID
	})

	team.Users = uniqueUsers.([]models.UserInner)

	err = teamService.getCollection().Update(team)
	if err != nil {
		return team, err
	}

	teams := user.Teams
	teams = append(teams, team.ToTeamInner())

	uniqueTeams := funk.UniqBy(teams, func(f models.TeamInner) primitive.ObjectID {
		return f.ID
	})

	user.Teams = uniqueTeams.([]models.TeamInner)

	err = userService.getCollection().Update(user)
	return team, err
}

// RemoveUser function
func (teamService *TeamService) RemoveUser(id interface{}, accountID primitive.ObjectID, userID interface{}) (updatedTeam *models.Team, err error) {

	team := &models.Team{}
	primitiveID, err := team.PrepareID(id)
	if err != nil {
		return team, err
	}
	err = teamService.getCollection().First(bson.M{"_id": primitiveID, "accountId": accountID}, team)
	if err != nil {
		return team, err
	}

	user := &models.User{}
	userPrimitiveID, err := user.PrepareID(userID)
	if err != nil {
		return team, err
	}
	err = userService.getCollection().First(bson.M{"_id": userPrimitiveID, "accountId": accountID}, user)
	if err != nil {
		return nil, err
	}

	team.Users = funk.Filter(team.Users, func(element models.UserInner) bool {
		return element.ID != user.ID
	}).([]models.UserInner)

	err = teamService.getCollection().Update(team)
	if err != nil {
		return team, err
	}

	user.Teams = funk.Filter(user.Teams, func(element models.TeamInner) bool {
		return element.ID != team.ID
	}).([]models.TeamInner)

	err = userService.getCollection().Update(user)
	return team, err
}
