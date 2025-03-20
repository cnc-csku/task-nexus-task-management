package services

import (
	"context"

	"github.com/cnc-csku/task-nexus-go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ReportService interface {
	GetStatusOverview(ctx context.Context, req *requests.GetTaskStatusOverviewRequest, userID string) (*responses.GetTaskStatusOverviewResponse, *errutils.Error)
	GetPriorityOverview(ctx context.Context, req *requests.GetTaskPriorityOverviewRequest, userID string) (*responses.GetTaskPriorityOverviewResponse, *errutils.Error)
	GetTypeOverview(ctx context.Context, req *requests.GetTaskTypeOverviewRequest, userID string) (*responses.GetTaskTypeOverviewResponse, *errutils.Error)
	GetAssigneeOverview(ctx context.Context, req *requests.GetTaskAssigneeOverviewRequest, userID string) (*responses.GetTaskAssigneeOverviewResponse, *errutils.Error)
	GetEpicTaskOverview(ctx context.Context, req *requests.GetEpicTaskOverviewRequest, userID string) (*responses.GetEpicTaskOverviewResponse, *errutils.Error)
}

type reportServiceImpl struct {
	userRepo      repositories.UserRepository
	projectRepo   repositories.ProjectRepository
	projectMember repositories.ProjectMemberRepository
	taskRepo      repositories.TaskRepository
}

func NewReportService(
	userRepo repositories.UserRepository,
	projectRepo repositories.ProjectRepository,
	projectMember repositories.ProjectMemberRepository,
	taskRepo repositories.TaskRepository,
) ReportService {
	return &reportServiceImpl{
		userRepo:      userRepo,
		projectRepo:   projectRepo,
		projectMember: projectMember,
		taskRepo:      taskRepo,
	}
}

func (s *reportServiceImpl) GetStatusOverview(ctx context.Context, req *requests.GetTaskStatusOverviewRequest, userID string) (*responses.GetTaskStatusOverviewResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	project, err := s.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest).WithDebugMessage("project not found")
	}

	member, err := s.projectMember.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.Forbidden).WithDebugMessage("permission denied")
	}

	tasks, err := s.taskRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if len(tasks) == 0 {
		return &responses.GetTaskStatusOverviewResponse{
			Statuses:   []responses.GetTaskStatusOverviewResponseStatuses{},
			TotalCount: 0,
		}, nil
	}

	statuses := make(map[string]int)
	for _, task := range tasks {
		statuses[task.Status]++
	}

	statusOverview := make([]responses.GetTaskStatusOverviewResponseStatuses, 0)
	for status, count := range statuses {
		statusOverview = append(statusOverview, responses.GetTaskStatusOverviewResponseStatuses{
			Status:  status,
			Count:   count,
			Percent: float64(count) / float64(len(tasks)) * 100,
		})
	}

	return &responses.GetTaskStatusOverviewResponse{
		Statuses:   statusOverview,
		TotalCount: len(tasks),
	}, nil
}

func (s *reportServiceImpl) GetPriorityOverview(ctx context.Context, req *requests.GetTaskPriorityOverviewRequest, userID string) (*responses.GetTaskPriorityOverviewResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	project, err := s.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest).WithDebugMessage("project not found")
	}

	member, err := s.projectMember.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.Forbidden).WithDebugMessage("permission denied")
	}

	tasks, err := s.taskRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if len(tasks) == 0 {
		return &responses.GetTaskPriorityOverviewResponse{
			Priorities: []responses.GetTaskPriorityOverviewResponsePriorities{},
		}, nil
	}

	priorities := make(map[string]int)
	for _, task := range tasks {
		priorities[task.Priority.String()]++
	}

	priorityOverview := make([]responses.GetTaskPriorityOverviewResponsePriorities, 0)
	for priority, count := range priorities {
		priorityOverview = append(priorityOverview, responses.GetTaskPriorityOverviewResponsePriorities{
			Priority: priority,
			Count:    count,
			Percent:  float64(count) / float64(len(tasks)) * 100,
		})
	}

	return &responses.GetTaskPriorityOverviewResponse{
		Priorities: priorityOverview,
		TotalCount: len(tasks),
	}, nil
}

func (s *reportServiceImpl) GetTypeOverview(ctx context.Context, req *requests.GetTaskTypeOverviewRequest, userID string) (*responses.GetTaskTypeOverviewResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	project, err := s.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest).WithDebugMessage("project not found")
	}

	member, err := s.projectMember.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.Forbidden).WithDebugMessage("permission denied")
	}

	tasks, err := s.taskRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if len(tasks) == 0 {
		return &responses.GetTaskTypeOverviewResponse{
			Types: []responses.GetTaskTypeOverviewResponseTypes{},
		}, nil
	}

	types := make(map[string]int)
	for _, task := range tasks {
		types[task.Type.String()]++
	}

	typeOverview := make([]responses.GetTaskTypeOverviewResponseTypes, 0)
	for taskType, count := range types {
		typeOverview = append(typeOverview, responses.GetTaskTypeOverviewResponseTypes{
			Type:    taskType,
			Count:   count,
			Percent: float64(count) / float64(len(tasks)) * 100,
		})
	}

	return &responses.GetTaskTypeOverviewResponse{
		Types:      typeOverview,
		TotalCount: len(tasks),
	}, nil
}

// To Be Disccused
func (s *reportServiceImpl) GetAssigneeOverview(ctx context.Context, req *requests.GetTaskAssigneeOverviewRequest, userID string) (*responses.GetTaskAssigneeOverviewResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	project, err := s.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest).WithDebugMessage("project not found")
	}

	member, err := s.projectMember.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.Forbidden).WithDebugMessage("permission denied")
	}

	tasks, err := s.taskRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if len(tasks) == 0 {
		return &responses.GetTaskAssigneeOverviewResponse{
			Assignees: []responses.GetTaskAssigneeOverviewResponseAssignees{},
		}, nil
	}

	assignees := make(map[string]int)
	for _, task := range tasks {
		for _, assignee := range task.Assignees {
			assignees[assignee.UserID.Hex()]++
		}
	}

	assigneeOverview := make([]responses.GetTaskAssigneeOverviewResponseAssignees, 0)
	for assigneeID, count := range assignees {
		bsonUserID, err := bson.ObjectIDFromHex(assigneeID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		user, err := s.userRepo.FindByID(ctx, bsonUserID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if user == nil {
			return nil, errutils.NewError(exceptions.ErrUserNotFound, errutils.BadRequest).WithDebugMessage("user not found")
		}

		var profileUrl = user.DefaultProfileUrl
		if user.UploadedProfileUrl != nil {
			profileUrl = *user.UploadedProfileUrl
		}

		assigneeOverview = append(assigneeOverview, responses.GetTaskAssigneeOverviewResponseAssignees{
			UserID:      assigneeID,
			FullName:    user.FullName,
			DisplayName: user.DisplayName,
			ProfileUrl:  profileUrl,
			Count:       count,
			Percent:     float64(count) / float64(len(tasks)) * 100,
		})
	}

	unassignedTaskCount := 0
	for _, task := range tasks {
		if len(task.Assignees) == 0 {
			unassignedTaskCount++
		}
	}
	assigneeOverview = append(assigneeOverview, responses.GetTaskAssigneeOverviewResponseAssignees{
		UserID:      "",
		FullName:    "Unassigned",
		DisplayName: "Unassigned",
		ProfileUrl:  "",
		Count:       unassignedTaskCount,
		Percent:     float64(unassignedTaskCount) / float64(len(tasks)) * 100,
	})

	return &responses.GetTaskAssigneeOverviewResponse{
		Assignees:  assigneeOverview,
		TotalCount: len(tasks),
	}, nil
}

func (s *reportServiceImpl) GetEpicTaskOverview(ctx context.Context, req *requests.GetEpicTaskOverviewRequest, userID string) (*responses.GetEpicTaskOverviewResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	project, err := s.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest).WithDebugMessage("project not found")
	}

	member, err := s.projectMember.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.Forbidden).WithDebugMessage("permission denied")
	}

	tasks, err := s.taskRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if len(tasks) == 0 {
		return &responses.GetEpicTaskOverviewResponse{
			Epics: []responses.GetEpicTaskOverviewResponseEpics{},
		}, nil
	}

	epicTasks := make(map[string]responses.GetEpicTaskOverviewResponseEpics)
	statusCounts := make(map[string]map[string]int)

	// Step 1: Get all epics and initialize status counts
	for _, task := range tasks {
		if task.Type == models.TaskTypeEpic {
			epicTasks[task.ID.Hex()] = responses.GetEpicTaskOverviewResponseEpics{
				TaskID:  task.ID.Hex(),
				TaskRef: task.TaskRef,
				Title:   task.Title,
			}
			statusCounts[task.ID.Hex()] = make(map[string]int)
		}
	}

	// Step 2: Count status of each epic
	for _, task := range tasks {
		if task.ParentID != nil {
			if counts, exists := statusCounts[task.ParentID.Hex()]; exists {
				counts[task.Status]++
			}
		}
	}

	// Step 3: Build response with percentage calculation
	responseEpics := make([]responses.GetEpicTaskOverviewResponseEpics, 0, len(epicTasks))
	for taskID, epic := range epicTasks {
		statusList := make([]responses.GetEpicTaskOverviewResponseEpicsStatuses, 0, len(statusCounts[taskID]))
		totalTasks := 0

		for status, count := range statusCounts[taskID] {
			totalTasks += count
			statusList = append(statusList, responses.GetEpicTaskOverviewResponseEpicsStatuses{
				Status: status,
				Count:  count,
			})
		}

		// Calculate percentage
		for i := range statusList {
			if totalTasks > 0 {
				statusList[i].Percent = float64(statusList[i].Count) / float64(totalTasks) * 100
			}
		}

		epic.Statuses = statusList
		epic.TotalCount = totalTasks
		responseEpics = append(responseEpics, epic)
	}

	return &responses.GetEpicTaskOverviewResponse{
		Epics:      responseEpics,
		TotalCount: len(responseEpics),
	}, nil
}
