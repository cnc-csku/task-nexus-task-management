package services

import (
	"context"

	"github.com/cnc-csku/task-nexus-go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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
