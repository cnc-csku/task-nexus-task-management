package mongo

import "go.mongodb.org/mongo-driver/v2/bson"

type projectFilter bson.M

func NewProjectFilter() projectFilter {
	return projectFilter{}
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
