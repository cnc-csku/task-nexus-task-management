package mongo

import (
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type projectFilter bson.M

func NewProjectFilter() projectFilter {
	return projectFilter{}
}

func (f projectFilter) WithID(id bson.ObjectID) {
	f["_id"] = id
}

func (f projectFilter) WithIDs(ids []bson.ObjectID) {
	f["_id"] = bson.M{
		"$in": ids,
	}
}

func (f projectFilter) WithWorkspaceID(workspaceID bson.ObjectID) {
	f["workspace_id"] = workspaceID
}

func (f projectFilter) WithName(name string) {
	f["name"] = name
}

func (f projectFilter) WithProjectPrefix(projectPrefix string) {
	f["project_prefix"] = projectPrefix
}

func (f projectFilter) WithUserID(userID bson.ObjectID) {
	f["members"] = bson.M{
		"$elemMatch": bson.M{
			"user_id": userID,
		},
	}
}

type projectUpdate bson.M

func NewProjectUpdate() projectUpdate {
	return projectUpdate{}
}

func (u projectUpdate) UpdatePositions(positions []string) {
	u["$set"] = bson.M{
		"positions": positions,
	}
}

func (u projectUpdate) UpdateWorkflows(workflows []bson.M) {
	u["$set"] = bson.M{
		"workflows": workflows,
	}
}

func (u projectUpdate) IncrementSprintRunningNumber() {
	u["$inc"] = bson.M{
		"sprint_running_number": 1,
	}
}

func (u projectUpdate) IncrementTaskRunningNumber() {
	u["$inc"] = bson.M{
		"task_running_number": 1,
	}
}

func (u projectUpdate) UpdateAttributeTemplates(attributeTemplates []bson.M) {
	u["$set"] = bson.M{
		"attributes_templates": attributeTemplates,
	}
}

func (u projectUpdate) UpdateSetupStatus(setupStatus models.ProjectSetupStatus) {
	u["$set"] = bson.M{
		"setup_status": setupStatus,
	}
}
