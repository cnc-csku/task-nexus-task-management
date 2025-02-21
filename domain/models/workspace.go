package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Workspace struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string        `bson:"name" json:"name"`
	CreatedBy bson.ObjectID `bson:"created_by" json:"createdBy"`
	CreatedAt time.Time     `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updatedAt"`
}
