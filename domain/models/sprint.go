package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Sprint struct {
	ID         bson.ObjectID `bson:"_id" json:"id"`
	ProjectID  bson.ObjectID `bson:"project_id" json:"projectId"`
	Title      string        `bson:"title" json:"title"`
	SprintGoal *string       `bson:"sprint_goal" json:"sprintGoal"`
	StartDate  *time.Time    `bson:"start_date" json:"startDate"`
	EndDate    *time.Time    `bson:"end_date" json:"endDate"`
	CreatedAt  time.Time     `bson:"created_at" json:"createdAt"`
	CreatedBy  bson.ObjectID `bson:"created_by" json:"createdBy"`
	UpdatedAt  time.Time     `bson:"updated_at" json:"updatedAt"`
	UpdatedBy  bson.ObjectID `bson:"updated_by" json:"updatedBy"`
}
