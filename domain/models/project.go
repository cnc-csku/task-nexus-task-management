package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Project struct {
	ID                  bson.ObjectID              `bson:"_id" json:"id"`
	WorkspaceID         bson.ObjectID              `bson:"workspace_id" json:"workspaceId"`
	Name                string                     `bson:"name" json:"name"`
	ProjectPrefix       string                     `bson:"project_prefix" json:"projectPrefix"`
	Description         string                     `bson:"description" json:"description"`
	Status              ProjectStatus              `bson:"status" json:"status"`
	SprintRunningNumber int                        `bson:"sprint_running_number" json:"sprintRunningNumber"`
	TaskRunningNumber   int                        `bson:"task_running_number" json:"taskRunningNumber"`
	Workflows           []ProjectWorkflow          `bson:"workflows" json:"workflows"`
	AttributeTemplates  []ProjectAttributeTemplate `bson:"attributes_templates" json:"attributesTemplates"`
	Positions           []string                   `bson:"positions" json:"positions"`
	CreatedAt           time.Time                  `bson:"created_at" json:"createdAt"`
	CreatedBy           bson.ObjectID              `bson:"created_by" json:"createdBy"`
	UpdatedAt           time.Time                  `bson:"updated_at" json:"updatedAt"`
	UpdatedBy           bson.ObjectID              `bson:"updated_by" json:"updatedBy"`
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

type ProjectWorkflow struct {
	PreviousStatuses []string `bson:"previous_statuses" json:"previousStatuses"`
	Status           string   `bson:"status" json:"status"`
	IsDefault        bool     `bson:"is_default" json:"isDefault"`
}

func GetDefaultWorkflows() []ProjectWorkflow {
	return []ProjectWorkflow{
		{Status: "Todo", IsDefault: true},
		{Status: "In Progress", PreviousStatuses: []string{"Todo"}},
		{Status: "Done", PreviousStatuses: []string{"In Progress"}},
	}
}

func GetDefaultPositions() []string {
	return []string{"Backend Developer", "Frontend Developer", "UX/UI Designer", "Quality Assurance"}
}

type ProjectAttributeTemplate struct {
	Name string           `bson:"name" json:"name"`
	Type KeyValuePairType `bson:"type" json:"type"`
}
