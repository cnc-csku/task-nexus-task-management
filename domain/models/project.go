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
	SetupStatus         ProjectStatus              `bson:"setup_status" json:"setupStatus"`
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
	IsDone           bool     `bson:"is_done" json:"isDone"`
}

func GetDefaultWorkflows() []ProjectWorkflow {
	return []ProjectWorkflow{
		{
			Status:           "Todo",
			PreviousStatuses: nil,
			IsDefault:        true,
			IsDone:           false,
		},
		{
			Status:           "In Progress",
			PreviousStatuses: []string{"Todo"},
			IsDefault:        false,
			IsDone:           false,
		},
		{
			Status:           "Done",
			PreviousStatuses: []string{"In Progress"},
			IsDefault:        false,
			IsDone:           true,
		},
	}
}

func GetDefaultPositions() []string {
	return []string{"Backend Developer", "Frontend Developer", "UX/UI Designer", "Quality Assurance"}
}

type ProjectAttributeTemplate struct {
	Name string           `bson:"name" json:"name"`
	Type KeyValuePairType `bson:"type" json:"type"`
}

type ProjectSetupStatus string

const (
	ProjectSetupStatusProjectCreated  ProjectSetupStatus = "PROJECT_CREATED"
	ProjectSetupStatusPositionConfig  ProjectSetupStatus = "POSITION_CONFIGURATION"
	ProjectSetupStatusOwnerConfig     ProjectSetupStatus = "OWNER_POSITION_CONFIGURATION"
	ProjectSetupStatusWorkflowConfig  ProjectSetupStatus = "WORKFLOW_CONFIGURATION"
	ProjectSetupStatusAttributeConfig ProjectSetupStatus = "ATTRIBUTE_TEMPLATE_CONFIGURATION"
	ProjectSetupStatusCompleted       ProjectSetupStatus = "COMPLETED"
)

func (p ProjectSetupStatus) String() string {
	return string(p)
}

func (p ProjectSetupStatus) IsValid() bool {
	switch p {
	case ProjectSetupStatusProjectCreated, ProjectSetupStatusPositionConfig, ProjectSetupStatusOwnerConfig, ProjectSetupStatusWorkflowConfig, ProjectSetupStatusAttributeConfig, ProjectSetupStatusCompleted:
		return true
	}
	return false
}
