package mongo

import "go.mongodb.org/mongo-driver/v2/bson"

type sprintFilter bson.M

func NewSprintFilter() sprintFilter {
	return sprintFilter{}
}

func (f sprintFilter) WithID(id bson.ObjectID) {
	f["_id"] = id
}

type sprintUpdater bson.M

func NewSprintUpdater() sprintUpdater {
	return sprintUpdater{}
}
