package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type sprintFilter bson.M

func NewSprintFilter() sprintFilter {
	return sprintFilter{}
}

func (f sprintFilter) WithID(id bson.ObjectID) {
	f["_id"] = id
}

func (f sprintFilter) WithProjectID(projectID bson.ObjectID) {
	f["project_id"] = projectID
}

func (f sprintFilter) WithEndDateGreaterThanOrEqualNowOrIsNull() {
	f["$or"] = []bson.M{
		{"end_date": bson.M{"$gte": time.Now()}},
		{"end_date": bson.M{"$eq": nil}},
	}
}

type sprintUpdater bson.M

func NewSprintUpdater() sprintUpdater {
	return sprintUpdater{}
}
