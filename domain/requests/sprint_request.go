package requests

import "time"

type CreateSprintRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
}

type GetSprintByIDRequest struct {
	SprintID string `param:"sprintId" validate:"required"`
}

type EditSprintRequest struct {
	ProjectID  string     `param:"projectId" validate:"required"`
	SprintID   string     `param:"sprintId" validate:"required"`
	Title      string     `json:"title" validate:"required"`
	SprintGoal string     `json:"sprintGoal"`
	Duration   *int       `json:"duration"`
	StartDate  *time.Time `json:"startDate"`
	EndDate    *time.Time `json:"endDate"`
}

type ListSprintPathParam struct {
	ProjectID string `param:"projectId" validate:"required"`
	IsActive  *bool  `query:"isActive"`
}
