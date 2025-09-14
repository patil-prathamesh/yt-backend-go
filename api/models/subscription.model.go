package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscription struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Subscriber primitive.ObjectID `json:"subscriber" bson:"subscriber"`
	Channel    primitive.ObjectID `json:"channel" bson:"channel"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
}
