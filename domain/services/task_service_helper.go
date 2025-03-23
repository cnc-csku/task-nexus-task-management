package services

import (
	"context"
	"fmt"

	"github.com/cnc-csku/task-nexus-go-lib/utils/array"
	"github.com/cnc-csku/task-nexus-go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func validateParentTaskType(childTaskType string, parentTaskType models.TaskType) *errutils.Error {
	switch childTaskType {
	case models.TaskTypeEpic.String():
		return errutils.NewError(exceptions.ErrInvalidParentTaskType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task type is not valid: %s", parentTaskType))
	case models.TaskTypeStory.String(), models.TaskTypeTask.String(), models.TaskTypeBug.String():
		if parentTaskType != models.TaskTypeEpic {
			return errutils.NewError(exceptions.ErrInvalidParentTaskType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task type is not valid: %s", parentTaskType))
		}
	case models.TaskTypeSubTask.String():
		if !array.ContainAny(
			[]string{parentTaskType.String()},
			[]string{
				models.TaskTypeStory.String(),
				models.TaskTypeTask.String(),
				models.TaskTypeBug.String(),
			},
		) {
			return errutils.NewError(exceptions.ErrInvalidParentTaskType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task type is not valid: %s", parentTaskType))
		}
	}
	return nil
}

func getDoneStatuses(project *models.Project) []string {
	isDoneStatuses := make([]string, 0)
	for _, workflow := range project.Workflows {
		if workflow.IsDone {
			isDoneStatuses = append(isDoneStatuses, workflow.Status)
		}
	}
	return isDoneStatuses
}

func getParentTasksMap(ctx context.Context, taskRepo repositories.TaskRepository, tasks []*models.Task) (map[string]*string, *errutils.Error) {
	parentTaskIDs := make([]bson.ObjectID, 0, len(tasks))
	for _, task := range tasks {
		if task.ParentID != nil {
			parentTaskIDs = append(parentTaskIDs, *task.ParentID)
		}
	}

	parentTasks, err := taskRepo.FindByIDs(ctx, parentTaskIDs)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	parentTasksMap := make(map[string]*string, len(parentTasks))
	for _, parentTask := range parentTasks {
		parentTasksMap[parentTask.ID.Hex()] = &parentTask.Title
	}

	return parentTasksMap, nil
}

func updatePreviousParentTask(
	ctx context.Context,
	taskRepo repositories.TaskRepository,
	updatedTask *models.Task,
) *errutils.Error {
	if updatedTask.ParentID == nil {
		return nil
	}

	previousParentTask, err := taskRepo.FindByID(ctx, *updatedTask.ParentID)
	if err != nil {
		return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if previousParentTask == nil {
		return errutils.NewError(exceptions.ErrParentTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task not found: %s", *updatedTask.ParentID))
	}

	childrenOfPreviousParentTasks, err := taskRepo.FindByParentID(ctx, previousParentTask.ID)
	if err != nil {
		return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	if len(childrenOfPreviousParentTasks) == 0 {
		_, err := taskRepo.UpdateHasChildren(ctx, &repositories.UpdateTaskHasChildrenRequest{
			ID:          previousParentTask.ID,
			HasChildren: false,
		})
		if err != nil {
			return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}
	}

	serviceErr := updateChildrenPointOfPreviousParentTask(ctx, taskRepo, updatedTask, previousParentTask)
	if serviceErr != nil {
		return serviceErr
	}

	return nil
}

func updateChildrenPointOfPreviousParentTask(
	ctx context.Context,
	taskRepo repositories.TaskRepository,
	updatedTask, previousParentTask *models.Task,
) *errutils.Error {
	if previousParentTask.Type == models.TaskTypeEpic {
		_, err := taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
			ID:            previousParentTask.ID,
			ChildrenPoint: previousParentTask.ChildrenPoint - updatedTask.ChildrenPoint,
		})
		if err != nil {
			return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}
	} else if array.ContainAny(
		[]string{previousParentTask.Type.String()},
		[]string{
			models.TaskTypeStory.String(),
			models.TaskTypeTask.String(),
			models.TaskTypeBug.String(),
		},
	) {
		// Current previous parent task is a story, task, or bug
		var totalPoint int
		for _, assignee := range updatedTask.Assignees {
			if assignee.Point != nil {
				totalPoint += *assignee.Point
			}
		}

		_, err := taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
			ID:            previousParentTask.ID,
			ChildrenPoint: previousParentTask.ChildrenPoint - totalPoint,
		})
		if err != nil {
			return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		// Update children point of the previous parent task's parent task (EPIC)
		if previousParentTask.ParentID != nil {
			epicParentTask, err := taskRepo.FindByID(ctx, *previousParentTask.ParentID)
			if err != nil {
				return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}

			_, err = taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
				ID:            *previousParentTask.ParentID,
				ChildrenPoint: epicParentTask.ChildrenPoint - totalPoint,
			})
			if err != nil {
				return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}
		}
	}

	return nil
}

func sortWorkflows(workflows []models.ProjectWorkflow) []models.ProjectWorkflow {
	// Step 1: Build graph adjacency list and in-degree map
	graph := make(map[string][]string)
	inDegree := make(map[string]int)
	statusMap := make(map[string]models.ProjectWorkflow)

	// Initialize graph nodes
	for _, workflow := range workflows {
		statusMap[workflow.Status] = workflow
		if _, exists := inDegree[workflow.Status]; !exists {
			inDegree[workflow.Status] = 0
		}
	}

	// Populate the adjacency list and compute in-degree
	for _, workflow := range workflows {
		for _, prevStatus := range workflow.PreviousStatuses {
			graph[prevStatus] = append(graph[prevStatus], workflow.Status)
			inDegree[workflow.Status]++ // Increase in-degree for the dependent status
		}
	}

	// Step 2: Find all statuses with in-degree 0 (starting points)
	var queue []string
	for status, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, status)
		}
	}

	// Step 3: Perform Topological Sort
	var sortedWorkflows []models.ProjectWorkflow
	for len(queue) > 0 {
		// Dequeue a status
		current := queue[0]
		queue = queue[1:]

		// Append to the sorted order
		sortedWorkflows = append(sortedWorkflows, statusMap[current])

		// Reduce in-degree for its children and enqueue if they become 0
		for _, nextStatus := range graph[current] {
			inDegree[nextStatus]--
			if inDegree[nextStatus] == 0 {
				queue = append(queue, nextStatus)
			}
		}
	}

	return sortedWorkflows
}

type UpdateParentTaskStatusToLowestWorkflowStatus struct {
	taskRepo      repositories.TaskRepository
	Workflows     []models.ProjectWorkflow
	ChildrenTasks []*models.Task
	ParentTask    *models.Task
	UpdaterUserID bson.ObjectID
}

func updateParentTaskStatusToLowestWorkflowStatus(ctx context.Context, in *UpdateParentTaskStatusToLowestWorkflowStatus) *errutils.Error {
	// Step 1: Sort workflows
	sortedWorkflows := sortWorkflows(in.Workflows)

	// Step 2: Find the lowest workflow status in children tasks
	statusIndexMap := make(map[string]int)
	for i, workflow := range sortedWorkflows {
		statusIndexMap[workflow.Status] = i
	}

	minIndex := len(sortedWorkflows) // Set to max initially
	for _, child := range in.ChildrenTasks {
		if idx, exists := statusIndexMap[child.Status]; exists && idx < minIndex {
			minIndex = idx
		}
	}

	// Step 3: Update the parent task to the lowest status found
	if minIndex < len(sortedWorkflows) {
		newParentStatus := sortedWorkflows[minIndex].Status
		if in.ParentTask.Status != newParentStatus { // Only update if different
			_, err := in.taskRepo.UpdateStatus(ctx, &repositories.UpdateTaskStatusRequest{
				ID:        in.ParentTask.ID,
				Status:    newParentStatus,
				UpdatedBy: in.UpdaterUserID,
			})
			if err != nil {
				return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}
		}
	}

	return nil
}

func assigneeesContainPoint(assignees []models.TaskAssignee) bool {
	for _, assignee := range assignees {
		if assignee.Point != nil {
			return true
		}
	}
	return false
}

func updateParentTaskPoints(
	ctx context.Context,
	taskRepo repositories.TaskRepository,
	task *models.Task,
	assignees []repositories.UpdateTaskAssigneesRequestAssignee,
	currentTotalPoint int,
) *errutils.Error {
	if task.ParentID == nil {
		return nil
	}

	parentTask, err := taskRepo.FindByID(ctx, *task.ParentID)
	if err != nil {
		return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if parentTask == nil {
		return errutils.NewError(exceptions.ErrParentTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task not found: %s", task.ParentID.Hex()))
	}

	var newTotalPoint int
	for _, assignee := range assignees {
		if assignee.Point != nil {
			newTotalPoint += *assignee.Point
		}
	}

	if parentTask.Type == models.TaskTypeEpic {
		_, err := taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
			ID:            *task.ParentID,
			ChildrenPoint: parentTask.ChildrenPoint - currentTotalPoint + newTotalPoint,
		})
		if err != nil {
			return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}
	} else if array.ContainAny(
		[]string{parentTask.Type.String()},
		[]string{
			models.TaskTypeStory.String(),
			models.TaskTypeTask.String(),
			models.TaskTypeBug.String(),
		},
	) {
		// Modify point for Story, Task, Bug
		_, err := taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
			ID:            *task.ParentID,
			ChildrenPoint: parentTask.ChildrenPoint - currentTotalPoint + newTotalPoint,
		})
		if err != nil {
			return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		// Update children point of the parent task's parent task (EPIC)
		if parentTask.ParentID != nil {
			epicParentTask, err := taskRepo.FindByID(ctx, *parentTask.ParentID)
			if err != nil {
				return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}

			_, err = taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
				ID:            *parentTask.ParentID,
				ChildrenPoint: epicParentTask.ChildrenPoint - currentTotalPoint + newTotalPoint,
			})
			if err != nil {
				return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}
		}
	}

	return nil
}
