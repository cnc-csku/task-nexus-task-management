package mongo

import (
	"time"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type taskFilter bson.M

func NewTaskFilter() taskFilter {
	return taskFilter{}
}

func (f taskFilter) WithID(id bson.ObjectID) {
	f["_id"] = id
}

func (f taskFilter) WithIDs(ids []bson.ObjectID) {
	f["_id"] = bson.M{
		"$in": ids,
	}
}

func (f taskFilter) WithTaskRef(taskRef string) {
	f["task_ref"] = taskRef
}

func (f taskFilter) WithUserApproval(userID bson.ObjectID) {
	f["approval.user_id"] = userID
}

func (f taskFilter) WithProjectID(projectID bson.ObjectID) {
	f["project_id"] = projectID
}

func (f taskFilter) WithStatuses(statuses []string) {
	f["status"] = bson.M{
		"$in": statuses,
	}
}

func (f taskFilter) WithNotInStatuses(statuses []string) {
	f["status"] = bson.M{
		"$nin": statuses,
	}
}

func (f taskFilter) WithParentID(parentID bson.ObjectID) {
	f["parent_id"] = parentID
}

func (f taskFilter) WithType(taskType models.TaskType) {
	f["type"] = taskType
}

func (f taskFilter) WithTypes(taskTypes []models.TaskType) {
	f["type"] = bson.M{
		"$in": taskTypes,
	}
}

func (f taskFilter) WithSprintID(sprintID bson.ObjectID) {
	f["sprint.current_sprint_id"] = sprintID
}

func (f taskFilter) WithUserIDs(userIDs []bson.ObjectID) {
	f["assignees.user_id"] = bson.M{
		"$in": userIDs,
	}
}

func (f taskFilter) WithSearchKeyword(keyword string) {
	f["$or"] = []bson.M{
		{"task_ref": bson.M{"$regex": keyword, "$options": "i"}},
		{"title": bson.M{"$regex": keyword, "$options": "i"}},
	}
}

type taskUpdate bson.M

func NewTaskUpdate() taskUpdate {
	return taskUpdate{}
}

func (u taskUpdate) UpdateDetail(in *repositories.UpdateTaskDetailRequest) {
	u["$set"] = bson.M{
		"title":       in.Title,
		"description": in.Description,
		"parent_id":   in.ParentID,
		"priority":    in.Priority,
		"updated_at":  time.Now(),
		"updated_by":  in.UpdatedBy,
	}
}

func (u taskUpdate) UpdateStatus(in *repositories.UpdateTaskStatusRequest) {
	u["$set"] = bson.M{
		"status":     in.Status,
		"updated_at": time.Now(),
		"updated_by": in.UpdatedBy,
	}
}

func (u taskUpdate) UpdateApprovals(in *repositories.UpdateTaskApprovalsRequest) {
	approval := make([]bson.M, len(in.Approval))
	for i, a := range in.Approval {
		approval[i] = bson.M{
			"user_id": a.UserID,
		}
	}

	u["$set"] = bson.M{
		"approval":   approval,
		"updated_at": time.Now(),
		"updated_by": in.UpdatedBy,
	}
}

func (u taskUpdate) ApproveTask(reason string) {
	u["$set"] = bson.M{
		"approval.$.reason": reason,
	}
}

func (u taskUpdate) UpdateAssignees(in *repositories.UpdateTaskAssigneesRequest) {
	assignees := make([]bson.M, len(in.Assignees))
	for i, a := range in.Assignees {
		assignees[i] = bson.M{
			"position": a.Position,
			"user_id":  a.UserID,
			"point":    a.Point,
		}
	}

	u["$set"] = bson.M{
		"assignees":  assignees,
		"updated_at": time.Now(),
		"updated_by": in.UpdatedBy,
	}
}

func (u taskUpdate) UpdateSprint(in *repositories.UpdateTaskSprintRequest) {
	u["$set"] = bson.M{
		"sprint": bson.M{
			"current_sprint_id":   in.CurrentSprintID,
			"previous_sprint_ids": in.PreviousSprintIDs,
		},
		"updated_at": time.Now(),
		"updated_by": in.UpdatedBy,
	}
}

func (u taskUpdate) UpdateHasChildren(hasChildren bool) {
	u["$set"] = bson.M{
		"has_children": hasChildren,
		"updated_at":   time.Now(),
	}
}

func (u taskUpdate) UpdateChildrenPoint(point int) {
	u["$set"] = bson.M{
		"children_point": point,
		"updated_at":     time.Now(),
	}
}
