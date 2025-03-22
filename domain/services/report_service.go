package services

import (
	"context"
	"fmt"

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
	GetEpicTaskOverview(ctx context.Context, req *requests.GetEpicTaskOverviewRequest, userID string) (*responses.GetEpicTaskOverviewResponse, *errutils.Error)
	GetAssigneeOverviewBySprint(ctx context.Context, req *requests.GetTaskAssigneeOverviewBySprintRequest, userID string) (*responses.GetAssigneeOverviewBySprintResponse, *errutils.Error)
}

type reportServiceImpl struct {
	userRepo      repositories.UserRepository
	projectRepo   repositories.ProjectRepository
	projectMember repositories.ProjectMemberRepository
	sprintRepo    repositories.SprintRepository
	taskRepo      repositories.TaskRepository
}

func NewReportService(
	userRepo repositories.UserRepository,
	projectRepo repositories.ProjectRepository,
	projectMember repositories.ProjectMemberRepository,
	sprintRepo repositories.SprintRepository,
	taskRepo repositories.TaskRepository,
) ReportService {
	return &reportServiceImpl{
		userRepo:      userRepo,
		projectRepo:   projectRepo,
		projectMember: projectMember,
		sprintRepo:    sprintRepo,
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

	epicTasks := make(map[string]responses.GetEpicTaskOverviewResponseEpics) // map[epicID]epicResponse
	statusCounts := make(map[string]map[string]int)                          // map[epicID]map[status]count

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

func (s *reportServiceImpl) GetAssigneeOverviewBySprint(ctx context.Context, req *requests.GetTaskAssigneeOverviewBySprintRequest, userID string) (*responses.GetAssigneeOverviewBySprintResponse, *errutils.Error) {
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

	var (
		tasks     []*models.Task
		sprintMap = make(map[string]models.Sprint)
	)
	if req.GetAllSprint != nil && *req.GetAllSprint {
		sprints, err := s.sprintRepo.FindByProjectID(ctx, bsonProjectID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if len(sprints) == 0 {
			return &responses.GetAssigneeOverviewBySprintResponse{
				Sprints:    []responses.GetAssigneeOverviewBySprintResponseSprint{},
				TotalCount: 0,
			}, nil
		}

		for _, sprint := range sprints {
			sprintMap[sprint.ID.Hex()] = sprint
		}

		tasks, err = s.taskRepo.FindByProjectID(ctx, bsonProjectID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if len(tasks) == 0 {
			return &responses.GetAssigneeOverviewBySprintResponse{
				Sprints:    []responses.GetAssigneeOverviewBySprintResponseSprint{},
				TotalCount: 0,
			}, nil
		}
	} else {
		activeSprints, err := s.sprintRepo.FindByProjectIDAndStatus(ctx, bsonProjectID, models.SprintStatusInProgress)
		fmt.Println("len(activeSprints)", len(activeSprints))
		for _, sprint := range activeSprints {
			fmt.Println("sprint", sprint)
		}
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if len(activeSprints) == 0 {
			return &responses.GetAssigneeOverviewBySprintResponse{
				Sprints:    []responses.GetAssigneeOverviewBySprintResponseSprint{},
				TotalCount: 0,
			}, nil
		}

		for _, sprint := range activeSprints {
			sprintMap[sprint.ID.Hex()] = sprint
		}

		activeSprintIDs := make([]bson.ObjectID, 0, len(activeSprints))
		for _, sprint := range activeSprints {
			activeSprintIDs = append(activeSprintIDs, sprint.ID)
		}

		tasks, err := s.taskRepo.FindByCurrentSprintIDs(ctx, activeSprintIDs)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if len(tasks) == 0 {
			return &responses.GetAssigneeOverviewBySprintResponse{
				Sprints:    []responses.GetAssigneeOverviewBySprintResponseSprint{},
				TotalCount: 0,
			}, nil
		}
	}

	type Total struct {
		TotalTask  int
		TotalPoint int
	}

	var (
		sprints           = make(map[string]responses.GetAssigneeOverviewBySprintResponseSprint) // map[sprintID]sprintResponse
		assigneeCounts    = make(map[string]map[string]Total)                                    // map[sprintID]map[assigneeID]Total
		assigneeUserIDMap = make(map[bson.ObjectID]struct{})                                     // map[assigneeID]struct{}
	)

	// Step 1: Filter tasks by sprint and initialize counts
	for _, task := range tasks {
		if task.Sprint == nil || task.Sprint.CurrentSprintID == nil {
			continue
		}

		sprintID := *task.Sprint.CurrentSprintID
		sprints[sprintID.Hex()] = responses.GetAssigneeOverviewBySprintResponseSprint{
			SprintID:    sprintID.Hex(),
			SprintTitle: sprintMap[sprintID.Hex()].Title,
		}
		assigneeCounts[sprintID.Hex()] = make(map[string]Total)
	}

	// Step 2: Count tasks by assignee
	for _, task := range tasks {
		if task.Sprint == nil || task.Sprint.CurrentSprintID == nil {
			continue
		}

		sprintID := *task.Sprint.CurrentSprintID
		if counts, exists := assigneeCounts[sprintID.Hex()]; exists {
			for _, assignee := range task.Assignees {
				var point int
				if assignee.Point != nil {
					point = *assignee.Point
				}
				if assignee.UserID != nil {
					counts[assignee.UserID.Hex()] = Total{
						TotalTask:  counts[assignee.UserID.Hex()].TotalTask + 1,
						TotalPoint: counts[assignee.UserID.Hex()].TotalPoint + point,
					}

					assigneeUserIDMap[*assignee.UserID] = struct{}{}
				}
			}
		}
	}

	assigneeUserIDs := make([]bson.ObjectID, 0, len(assigneeUserIDMap))
	for assigneeUserID := range assigneeUserIDMap {
		assigneeUserIDs = append(assigneeUserIDs, assigneeUserID)
	}

	users, err := s.userRepo.FindByIDs(ctx, assigneeUserIDs)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	userMap := make(map[bson.ObjectID]models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// Step 3: Build response with percentage calculation
	responseSprints := make([]responses.GetAssigneeOverviewBySprintResponseSprint, 0, len(sprints))
	for sprintID, sprint := range sprints {
		var (
			assigneeList = make([]responses.GetAssigneeOverviewBySprintResponseSprintAssignee, 0, len(assigneeCounts[sprintID]))
			totalTasks   = 0
			totalPoints  = 0
		)

		for assigneeID, total := range assigneeCounts[sprintID] {
			totalTasks += total.TotalTask
			totalPoints += total.TotalPoint

			userID, err := bson.ObjectIDFromHex(assigneeID)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}

			var profileURL = userMap[userID].DefaultProfileUrl
			if userMap[userID].UploadedProfileUrl != nil {
				profileURL = *userMap[userID].UploadedProfileUrl
			}

			assigneeList = append(assigneeList, responses.GetAssigneeOverviewBySprintResponseSprintAssignee{
				UserID:      assigneeID,
				FullName:    userMap[userID].FullName,
				DisplayName: userMap[userID].DisplayName,
				ProfileUrl:  profileURL,
				TaskCount:   total.TotalTask,
				PointCount:  total.TotalPoint,
			})
		}

		// Calculate percentage
		for i := range assigneeList {
			if totalTasks > 0 {
				assigneeList[i].TaskPercent = float64(assigneeList[i].TaskCount) / float64(totalTasks) * 100
			}
			if totalPoints > 0 {
				assigneeList[i].PointPercent = float64(assigneeList[i].PointCount) / float64(totalPoints) * 100
			}
		}

		sprint.Assignees = assigneeList
		sprint.TotalTask = totalTasks
		sprint.TotalPoint = totalPoints
		responseSprints = append(responseSprints, sprint)
	}

	return &responses.GetAssigneeOverviewBySprintResponse{
		Sprints:    responseSprints,
		TotalCount: len(responseSprints),
	}, nil
}
