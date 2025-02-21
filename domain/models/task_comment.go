package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskComment struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	Content   string        `bson:"content" json:"content"`
	UserID    bson.ObjectID `bson:"user_id" json:"userId"`
	TaskID    string        `bson:"task_id" json:"taskId"`
	CreatedAt time.Time     `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updatedAt"`
}
