package mongo

import "go.mongodb.org/mongo-driver/v2/bson"

type taskFilter bson.M

func NewTaskFilter() taskFilter {
	return taskFilter{}
}

func (f taskFilter) WithID(id bson.ObjectID) {
	f["_id"] = id
}

func (f taskFilter) WithTaskID(taskID string) {
	f["task_id"] = taskID
}

type taskUpdate bson.M

func NewTaskUpdate() taskUpdate {
	return taskUpdate{}
}
