package mongo

import "go.mongodb.org/mongo-driver/v2/bson"

type workspaceMemberFilter bson.M

func NewWorkspaceMemberFilter() workspaceMemberFilter {
	return workspaceMemberFilter{}
}

func (f workspaceMemberFilter) WithWorkspaceID(workspaceID bson.ObjectID) {
	f["workspace_id"] = workspaceID
}

func (f workspaceMemberFilter) WithUserID(userID bson.ObjectID) {
	f["user_id"] = userID
}
