package mongo

import (
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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

func (f projectMemberFilter) WithProjectIDs(projectIDs []bson.ObjectID) {
	f["project_id"] = bson.M{
		"$in": projectIDs,
	}
}

func (f projectMemberFilter) WithRole(role models.ProjectMemberRole) {
	f["role"] = role
}

func (f projectMemberFilter) WithPositions(positions []string) {
	f["position"] = bson.M{
		"$in": positions,
	}
}
