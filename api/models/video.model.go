package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	VideoFile   string             `json:"videoFile" bson:"videoFile"`
	Thumbnail   string             `json:"thumbnail" bson:"thumbnail"`
	Owner       primitive.ObjectID `json:"owner" bson:"owner"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Duration    int                `json:"duration" bson:"duration"`
	Views       int                `json:"views" bson:"views"`
	IsPublished bool               `json:"isPublished" bson:"isPublished"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
