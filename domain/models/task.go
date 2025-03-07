package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Task struct {
	ID          bson.ObjectID  `bson:"_id" json:"id"`
	TaskRef     string         `bson:"task_ref" json:"taskRef"`
	ProjectID   bson.ObjectID  `bson:"project_id" json:"projectId"`
	Title       string         `bson:"title" json:"title"`
	Description string         `bson:"description" json:"description"`
	ParentID    *bson.ObjectID `bson:"parent_id" json:"parentId"`
	Type        TaskType       `bson:"type" json:"type"`
	Status      string         `bson:"status" json:"status"`
	Priority    *TaskPriority  `bson:"priority" json:"priority"`
	Approvals   []TaskApproval `bson:"approvals" json:"approvals"`
	Assignees   []TaskAssignee `bson:"assignees" json:"assignees"`
	Sprint      *TaskSprint    `bson:"sprint" json:"sprint"`
	CreatedAt   time.Time      `bson:"created_at" json:"createdAt"`
	CreatedBy   bson.ObjectID  `bson:"created_by" json:"createdBy"`
	UpdatedAt   time.Time      `bson:"updated_at" json:"updatedAt"`
	UpdatedBy   bson.ObjectID  `bson:"updated_by" json:"updatedBy"`
}

type TaskType string

const (
	TaskTypeEpic    TaskType = "EPIC"
	TaskTypeStory   TaskType = "STORY"
	TaskTypeTask    TaskType = "TASK"
	TaskTypeBug     TaskType = "BUG"
	TaskTypeSubTask TaskType = "SUB_TASK"
)

func (t TaskType) String() string {
	return string(t)
}

func (t TaskType) IsValid() bool {
	switch t {
	case TaskTypeEpic, TaskTypeStory, TaskTypeTask, TaskTypeBug, TaskTypeSubTask:
		return true
	}
	return false
}

type TaskPriority string

const (
	TaskPriorityLow      TaskPriority = "LOW"
	TaskPriorityMedium   TaskPriority = "MEDIUM"
	TaskPriorityHigh     TaskPriority = "HIGH"
	TaskPriotityCritical TaskPriority = "CRITICAL"
)

func (t TaskPriority) String() string {
	return string(t)
}

func (t TaskPriority) IsValid() bool {
	switch t {
	case TaskPriorityLow, TaskPriorityMedium, TaskPriorityHigh, TaskPriotityCritical:
		return true
	}
	return false
}

type TaskApproval struct {
	Reason string        `bson:"reason" json:"reason"`
	UserID bson.ObjectID `bson:"user_id" json:"userId"`
}

type TaskAssignee struct {
	Postion string        `bson:"position" json:"position"`
	UserID  bson.ObjectID `bson:"user_id" json:"userId"`
	Point   *int          `bson:"point" json:"point"`
}

type TaskSprint struct {
	PreviousSprintIDs []bson.ObjectID `bson:"previous_sprint_ids" json:"previousSprintIds"`
	CurrentSprintID   bson.ObjectID   `bson:"current_sprint_id" json:"currentSprintId"`
}
