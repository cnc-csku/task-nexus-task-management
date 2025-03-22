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

type GetAssigneeOverviewBySprintResponse struct {
	Sprints    []GetAssigneeOverviewBySprintResponseSprint `json:"sprints"`
	TotalCount int                                         `json:"totalCount"`
}

type GetAssigneeOverviewBySprintResponseSprint struct {
	SprintID    string                                              `json:"sprintID"`
	SprintTitle string                                              `json:"sprintTitle"`
	Assignees   []GetAssigneeOverviewBySprintResponseSprintAssignee `json:"assignees"`
	TotalTask   int                                                 `json:"totalTask"`
	TotalPoint  int                                                 `json:"totalPoint"`
}

type GetAssigneeOverviewBySprintResponseSprintAssignee struct {
	UserID       string  `json:"userID"`
	FullName     string  `json:"fullName"`
	DisplayName  string  `json:"displayName"`
	ProfileUrl   string  `json:"profileUrl"`
	TaskCount    int     `json:"taskCount"`
	PointCount   int     `json:"pointCount"`
	TaskPercent  float64 `json:"taskPercent"`
	PointPercent float64 `json:"pointPercent"`
}
