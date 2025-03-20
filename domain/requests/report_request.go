package requests

type GetStatusOverviewRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
}

type GetPriorityOverviewRequest struct {
	ProjectID string `param:"projectId" validate:"required"`
}
