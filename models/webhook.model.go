package models

import (
	"github.com/Kamva/mgm/v3"
)

// Webhook struct
type Webhook struct {
	mgm.DefaultModel `bson:",inline"`
	Payload          map[string]interface{} `json:"payload" bson:"payload"`
}
