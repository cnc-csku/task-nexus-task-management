package requests

type GetTaskStatusOverviewRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
}

type GetTaskPriorityOverviewRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
}

type GetTaskTypeOverviewRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
}

type GetTaskAssigneeOverviewRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
}

type GetEpicTaskOverviewRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
}
