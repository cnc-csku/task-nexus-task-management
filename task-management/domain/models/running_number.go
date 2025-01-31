package models

import "go.mongodb.org/mongo-driver/v2/bson"

type RunningNumber struct {
	ID       bson.ObjectID     `bson:"_id" json:"id"`
	Sequence int               `bson:"sequence" json:"sequence"`
	Type     RunningNumberType `bson:"type" json:"type"`
}

type RunningNumberType string

const (
	RunningNumberTypeTask   RunningNumberType = "TASK"
	RunningNumberTypeSprint RunningNumberType = "SPRINT"
)

func (r RunningNumberType) String() string {
	return string(r)
}

func (r RunningNumberType) IsValid() bool {
	switch r {
	case RunningNumberTypeTask, RunningNumberTypeSprint:
		return true
	}
	return false
}
