package responses

type GetTaskStatusOverviewResponse struct {
	Statuses   []GetTaskStatusOverviewResponseStatuses `json:"statuses"`
	TotalCount int                                     `json:"totalCount"`
}

type GetTaskStatusOverviewResponseStatuses struct {
	Status  string  `json:"status"`
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

type GetTaskPriorityOverviewResponse struct {
	Priorities []GetTaskPriorityOverviewResponsePriorities `json:"priorities"`
	TotalCount int                                         `json:"totalCount"`
}

type GetTaskPriorityOverviewResponsePriorities struct {
	Priority string  `json:"priority"`
	Count    int     `json:"count"`
	Percent  float64 `json:"percent"`
}

type GetTaskTypeOverviewResponse struct {
	Types      []GetTaskTypeOverviewResponseTypes `json:"types"`
	TotalCount int                                `json:"totalCount"`
}

type GetTaskTypeOverviewResponseTypes struct {
	Type    string  `json:"type"`
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

type GetTaskAssigneeOverviewResponse struct {
	Assignees  []GetTaskAssigneeOverviewResponseAssignees `json:"assignees"`
	TotalCount int                                        `json:"totalCount"`
}

type GetTaskAssigneeOverviewResponseAssignees struct {
	UserID      string  `json:"userID"`
	FullName    string  `json:"fullName"`
	DisplayName string  `json:"displayName"`
	ProfileUrl  string  `json:"profileUrl"`
	Count       int     `json:"count"`
	Percent     float64 `json:"percent"`
}

type GetEpicTaskOverviewResponse struct {
	Epics      []GetEpicTaskOverviewResponseEpics `json:"epics"`
	TotalCount int                                `json:"totalCount"`
}

type GetEpicTaskOverviewResponseEpics struct {
	TaskID     string                                     `json:"taskID"`
	TaskRef    string                                     `json:"taskRef"`
	Title      string                                     `json:"title"`
	Statuses   []GetEpicTaskOverviewResponseEpicsStatuses `json:"statuses"`
	TotalCount int                                        `json:"totalCount"`
}

type GetEpicTaskOverviewResponseEpicsStatuses struct {
	Status  string  `json:"status"`
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}
