package responses

type GetStatusOverviewResponse struct {
	Statuses   []GetStatusOverviewResponseStatuses `json:"statuses"`
	TotalCount int                                 `json:"totalCount"`
}

type GetStatusOverviewResponseStatuses struct {
	Status  string  `json:"status"`
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

type GetPriorityOverviewResponse struct {
	Priorities []GetPriorityOverviewResponsePriorities `json:"priorities"`
	TotalCount int                                     `json:"totalCount"`
}

type GetPriorityOverviewResponsePriorities struct {
	Priority string `json:"priority"`
	Count    int    `json:"count"`
	Percent  int    `json:"percent"`
}
