package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Project struct {
	ID                  bson.ObjectID       `bson:"_id" json:"id"`
	WorkspaceID         bson.ObjectID       `bson:"workspace_id" json:"workspaceId"`
	Name                string              `bson:"name" json:"name"`
	ProjectPrefix       string              `bson:"project_prefix" json:"projectPrefix"`
	Description         *string             `bson:"description" json:"description"`
	Status              ProjectStatus       `bson:"status" json:"status"`
	SprintRunningNumber int                 `bson:"sprint_running_number" json:"sprintRunningNumber"`
	TaskRunningNumber   int                 `bson:"task_running_number" json:"taskRunningNumber"`
	Workflows           []Workflow          `bson:"workflows" json:"workflows"`
	AttributeTemplates  []AttributeTemplate `bson:"attributes_templates" json:"attributesTemplates"`
	Positions           []string            `bson:"positions" json:"positions"`
	CreatedAt           time.Time           `bson:"created_at" json:"createdAt"`
	CreatedBy           bson.ObjectID       `bson:"created_by" json:"createdBy"`
	UpdatedAt           time.Time           `bson:"updated_at" json:"updatedAt"`
	UpdatedBy           bson.ObjectID       `bson:"updated_by" json:"updatedBy"`
}

type Workflow struct {
	PreviousStatuses []string `bson:"previous_statuses" json:"previousStatuses"`
	Status           string   `bson:"status" json:"status"`
}

func GetDefaultWorkflows() []Workflow {
	return []Workflow{
		{Status: "TODO"},
		{Status: "IN_PROGRESS", PreviousStatuses: []string{"TODO"}},
		{Status: "DONE", PreviousStatuses: []string{"IN_PROGRESS"}},
	}
}

type AttributeTemplate struct {
	Name  string      `bson:"name" json:"name"`
	Type  string      `bson:"type" json:"type"`
	Value interface{} `bson:"value" json:"value"`
}

type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "ACTIVE"
	ProjectStatusInactive ProjectStatus = "INACTIVE"
)

func (p ProjectStatus) String() string {
	return string(p)
}

func (p ProjectStatus) IsValid() bool {
	switch p {
	case ProjectStatusActive, ProjectStatusInactive:
		return true
	}
	return false
}
