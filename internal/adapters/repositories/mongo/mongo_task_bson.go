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
	f["approvals.user_id"] = userID
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

func (f taskFilter) WithNoParentID() {
	f["parent_id"] = bson.M{
		"$eq": nil,
	}
}

func (f taskFilter) WithType(taskType models.TaskType) {
	f["type"] = taskType
}

func (f taskFilter) WithTypes(taskTypes []models.TaskType) {
	f["type"] = bson.M{
		"$in": taskTypes,
	}
}

func (f taskFilter) WithCurrentSprintID(sprintID bson.ObjectID) {
	f["sprint.current_sprint_id"] = sprintID
}

func (f taskFilter) WithCurrentSprintIDs(sprintID []bson.ObjectID) {
	f["sprint.current_sprint_id"] = bson.M{
		"$in": sprintID,
	}
}

func (f taskFilter) WithPreviousSprintID(sprintID bson.ObjectID) {
	f["sprint.previous_sprint_ids"] = sprintID
}

func (f taskFilter) WithNoSprintID() {
	f["sprint.current_sprint_id"] = bson.M{
		"$eq": nil,
	}
}

func (f taskFilter) WithUserIDs(userIDs []bson.ObjectID) {
	f["assignees.user_id"] = bson.M{
		"$in": userIDs,
	}
}

func (f taskFilter) WithPositions(positions []string) {
	f["assignees.position"] = bson.M{
		"$in": positions,
	}
}

func (f taskFilter) WithSearchKeyword(keyword string) {
	f["$or"] = []bson.M{
		{"task_ref": bson.M{"$regex": keyword, "$options": "i"}},
		{"title": bson.M{"$regex": keyword, "$options": "i"}},
	}
}

func (f taskFilter) WithCurrentSprintIDAndPreviousSprintIDs(sprintID bson.ObjectID) {
	f["$or"] = []bson.M{
		{"sprint.current_sprint_id": sprintID},
		{"sprint.previous_sprint_ids": bson.M{
			"$in": []bson.ObjectID{sprintID},
		}},
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
		"priority":    in.Priority,
		"start_date":  in.StartDate,
		"due_date":    in.DueDate,
		"updated_at":  time.Now(),
		"updated_by":  in.UpdatedBy,
	}
}

func (u taskUpdate) UpdateTitle(in *repositories.UpdateTaskTitleRequest) {
	u["$set"] = bson.M{
		"title":      in.Title,
		"updated_at": time.Now(),
		"updated_by": in.UpdatedBy,
	}
}

func (u taskUpdate) UpdateParentID(in *repositories.UpdateTaskParentIDRequest) {
	u["$set"] = bson.M{
		"parent_id":  in.ParentID,
		"updated_at": time.Now(),
		"updated_by": in.UpdatedBy,
	}
}

func (u taskUpdate) UpdateType(in *repositories.UpdateTaskTypeRequest) {
	u["$set"] = bson.M{
		"type":       in.Type,
		"updated_at": time.Now(),
		"updated_by": in.UpdatedBy,
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
	approvals := make([]bson.M, len(in.Approval))
	for i, a := range in.Approval {
		approvals[i] = bson.M{
			"user_id": a.UserID,
		}
	}

	u["$set"] = bson.M{
		"approvals":  approvals,
		"updated_at": time.Now(),
		"updated_by": in.UpdatedBy,
	}
}

func (u taskUpdate) ApproveTask(reason string) {
	u["$set"] = bson.M{
		"approvals.$.is_approved": true,
		"approvals.$.reason":      reason,
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

func (u taskUpdate) UpdateCurrentSprintID(currentSprintID *bson.ObjectID, updatedBy bson.ObjectID) {
	u["$set"] = bson.M{
		"sprint": bson.M{
			"current_sprint_id": currentSprintID,
		},
		"updated_at": time.Now(),
		"updated_by": updatedBy,
	}
}

func (u taskUpdate) UpdatePreviousSprintIDs(in *repositories.UpdateTaskPreviousSprintIDsRequest) {
	u["$set"] = bson.M{
		"sprint.previous_sprint_ids": in.PreviousSprintIDs,
		"updated_at":                 time.Now(),
		"updated_by":                 in.UpdatedBy,
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

func (u taskUpdate) UpdateAttributes(in *repositories.UpdateTaskAttributesRequest) {
	attrs := make([]bson.M, len(in.Attributes))
	for i, a := range in.Attributes {
		attrs[i] = bson.M{
			"key":   a.Key,
			"value": a.Value,
		}
	}

	u["$set"] = bson.M{
		"attributes": attrs,
		"updated_at": time.Now(),
		"updated_by": in.UpdatedBy,
	}
}
