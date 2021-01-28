package model

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Uuid      uuid.UUID          `json:"uuid" bson:"uuid"`
	Question  string             `json:"question" bson:"question"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	CreatedBy string             `json:"createdBy" bon:"createdBy"`
	UpdatedAt time.Time          `json:"updateAt" bson:"updateAt"`
	UpdatedBy string             `json:"updateBy" bson:"updateBy"`
	IsActive  bool               `json:"isActive" bson:"isActive"`
}
