package models

import (
	"github.com/Kamva/mgm/v3"
)

// User struct
type Email struct {
	mgm.DefaultModel `bson:",inline"`
	Code             string `json:"code" bson:"code"`
	Lang             string `json:"lang" bson:"lang"`
	Subject          string `json:"subject" bson:"subject"`
	Body             string `json:"body" bson:"body"`
}
