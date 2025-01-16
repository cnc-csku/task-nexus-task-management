package mongo

import "go.mongodb.org/mongo-driver/v2/bson"

type workspaceFilter bson.M

func NewWorkspaceFilter() workspaceFilter {
	return workspaceFilter{}
}

func (f workspaceFilter) WithWorkspaceID(workspaceID bson.ObjectID) {
	f["_id"] = workspaceID
}

func (f workspaceFilter) WithMemberUserID(userID bson.ObjectID) {
	f["members"] = bson.M{
		"$elemMatch": bson.M{
			"user_id": userID,
		},
	}
}
