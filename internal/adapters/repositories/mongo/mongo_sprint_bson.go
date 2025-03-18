package mongo

import (
	"time"

	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
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

func (u sprintUpdater) UpdateStatus(in *repositories.UpdateSprintStatusRequest) {
	u["$set"] = bson.M{
		"status":     in.Status,
		"updated_at": time.Now(),
		"updated_by": in.UpdatedBy,
	}
}
