package domain

import (
	"context"
)

type TaskService struct {
	taskRepository TaskRepository
}

func NewTaskService(taskRepository TaskRepository) *TaskService {
	return &TaskService{
		taskRepository: taskRepository,
	}
}

func (ts *TaskService) CreateTask(ctx context.Context, r *CreateTaskRequest) (*Task, error) {
	t, err := CreateTask(r.TaskText, r.TaskOrder, r.UserID)
	if err != nil {
		return nil, err
	}

	if err := t.Validate(); err != nil {
		return nil, err
	}

	if err := ts.taskRepository.SaveTask(ctx, t); err != nil {
		return nil, err
	}

	return t, err
}

func (ts *TaskService) UpdateTask(ctx context.Context, r *UpdateTaskRequest) (*Task, error) {
	t, err := ts.taskRepository.FindByID(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	r.EnrichTask(t)

	if err := t.Validate(); err != nil {
		return nil, err
	}

	if err := ts.taskRepository.SaveTask(ctx, t); err != nil {
		return nil, err
	}

	return t, err
}

func (ts *TaskService) DeleteTask(ctx context.Context, r *DeleteTaskRequest) error {
	t, err := ts.taskRepository.FindByID(ctx, r.ID)
	if err != nil {
		return err
	}

	if err := ts.taskRepository.DeleteTask(ctx, t); err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) MarkTaskDone(ctx context.Context, r *MarkTaskDoneRequest) (*Task, error) {
	t, err := ts.taskRepository.FindByID(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	t.MarkDone()
	t.BeforeUpdate()

	if err := t.Validate(); err != nil {
		return nil, err
	}

	if err := ts.taskRepository.SaveTask(ctx, t); err != nil {
		return nil, err
	}

	return t, err
}

func (ts *TaskService) MarkTaskNotDone(ctx context.Context, r *MarkTaskNotDoneRequest) (*Task, error) {
	t, err := ts.taskRepository.FindByID(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	t.MarkNotDone()
	t.BeforeUpdate()

	if err := t.Validate(); err != nil {
		return nil, err
	}

	if err := ts.taskRepository.SaveTask(ctx, t); err != nil {
		return nil, err
	}

	return t, err
}

func (ts *TaskService) GetAllByUser(ctx context.Context, u *User) ([]*Task, error) {
	tasks, err := ts.taskRepository.GetAllByUserID(ctx, u.GetID())
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
