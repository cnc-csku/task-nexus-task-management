package mongo

import "go.mongodb.org/mongo-driver/v2/bson"

type taskCommentFilter bson.M

func NewTaskCommentFilter() taskCommentFilter {
	return taskCommentFilter{}
}

func (f taskCommentFilter) WithID(id bson.ObjectID) {
	f["_id"] = id
}

func (f taskCommentFilter) WithIDs(ids []bson.ObjectID) {
	f["_id"] = bson.M{
		"$in": ids,
	}
}

func (f taskCommentFilter) WithTaskID(taskID string) {
	f["task_id"] = taskID
}

type taskCommentUpdate bson.M

func NewTaskCommentUpdate() taskCommentUpdate {
	return taskCommentUpdate{}
}
