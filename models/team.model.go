package models

import (
	"github.com/Kamva/mgm/v3"
	"github.com/devinterface/structomap"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamInner struct {
	ID   primitive.ObjectID `json:"id" bson:"id"`
	Code string             `json:"code" bson:"code"`
	Name string             `json:"name" bson:"name"`
}

// Team struct
type Team struct {
	mgm.DefaultModel `bson:",inline"`
	Code             string             `json:"code" bson:"code"`
	Name             string             `json:"name" bson:"name"`
	Users            []UserInner        `json:"users" bson:"users"`
	AccountID        primitive.ObjectID `json:"accountId" bson:"accountId"`
}

// TeamSerializer function
type TeamSerializer struct {
	*structomap.Base
}

// ShowTeamSerializer function
func ShowTeamSerializer() *TeamSerializer {
	g := &TeamSerializer{structomap.New()}
	g.UseCamelCase().Pick("ID", "Name", "Code", "Users", "CreatedAt", "UpdatedAt")
	return g
}

func (team *Team) ToTeamInner() TeamInner {
	teamInner := TeamInner{}
	teamInner.ID = team.ID
	teamInner.Code = team.Code
	teamInner.Name = team.Name
	return teamInner
}
