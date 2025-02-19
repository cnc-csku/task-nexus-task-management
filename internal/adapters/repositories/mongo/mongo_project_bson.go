package mongo

import "go.mongodb.org/mongo-driver/v2/bson"

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

func (u projectUpdate) AddPositions(position []string) {
	u["$push"] = bson.M{
		"positions": bson.M{
			"$each": position,
		},
	}
}

func (u projectUpdate) AddMembers(member []bson.M) {
	u["$push"] = bson.M{
		"members": bson.M{
			"$each": member,
		},
	}
}

func (u projectUpdate) AddWorkflows(workflow []bson.M) {
	u["$push"] = bson.M{
		"workflows": bson.M{
			"$each": workflow,
		},
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
