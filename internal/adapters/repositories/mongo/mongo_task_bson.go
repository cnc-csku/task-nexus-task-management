package mongo

import (
	"fmt"
	"time"

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

func (f taskFilter) WithTaskRef(taskRef string) {
	f["task_ref"] = taskRef
}

func (f taskFilter) WithUserApproval(userID bson.ObjectID) {
	f["approval.user_id"] = userID
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
		"type":        in.Type,
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
	fmt.Println("assignees", assignees)

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
