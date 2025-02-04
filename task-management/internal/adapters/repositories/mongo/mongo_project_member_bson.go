package mongo

import "go.mongodb.org/mongo-driver/v2/bson"

type projectMemberFilter bson.M

func NewProjectMemberFilter() projectMemberFilter {
	return projectMemberFilter{}
}

func (f projectMemberFilter) WithUserID(userID bson.ObjectID) {
	f["user_id"] = userID
}

func (f projectMemberFilter) WithProjectID(projectID bson.ObjectID) {
	f["project_id"] = projectID
}
