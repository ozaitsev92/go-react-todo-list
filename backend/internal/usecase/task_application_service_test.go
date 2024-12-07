package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/config"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	repo "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/memory"
	"github.com/ozaitsev92/tododdd/internal/usecase"
)

func TestTaskUseCaseCreate(t *testing.T) {
	type args struct {
		text   string
		userId uuid.UUID
	}

	type testCase struct {
		name    string
		args    args
		want    task.Task
		wantErr error
	}

	tests := []testCase{
		{
			name: "Success",
			args: args{
				text:   "test",
				userId: uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
			},
			want: task.Task{
				Text:      "test",
				Completed: false,
				UserID:    uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
			},
			wantErr: nil,
		},
		{
			name: "Empty task",
			args: args{
				text:   "",
				userId: uuid.New(),
			},
			want:    task.Task{},
			wantErr: task.ErrInvalidText,
		},
		{
			name: "Empty userId",
			args: args{
				text:   "test",
				userId: uuid.Nil,
			},
			want:    task.Task{},
			wantErr: task.ErrInvalidUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := repo.NewRepository(config.Config{})

			s := usecase.NewTaskUseCase(repo)
			newTask, err := s.CreateTask(context.Background(), tt.args.text, tt.args.userId)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("s.CreateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && newTask.Text != tt.args.text {
				t.Errorf("s.CreateTask() Text = %v, want %v", newTask.Text, tt.want.Text)
			}
			if err == nil && newTask.Completed != tt.want.Completed {
				t.Errorf("s.CreateTask() Completed = %v, want %v", newTask.Completed, tt.want.Completed)
			}
			if err == nil && newTask.UpdatedAt.IsZero() {
				t.Error("s.CreateTask() UpdatedAt is zero")
			}
		})
	}
}

func TestTaskGetAllTasksForUser(t *testing.T) {
	type testCase struct {
		name    string
		tasks   []task.Task
		userId  uuid.UUID
		want    []task.Task
		wantErr error
	}

	tests := []testCase{
		{
			name: "Success",
			tasks: []task.Task{
				{
					ID:        uuid.New(),
					Text:      "test",
					Completed: false,
					UserID:    uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
				},
				{
					ID:        uuid.New(),
					Text:      "test 123",
					Completed: true,
					UserID:    uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
				},
				{
					ID:        uuid.New(),
					Text:      "lorem ipsum",
					Completed: true,
					UserID:    uuid.MustParse("842efa19-5926-4cac-8ca5-2fcb7d8056b1"),
				},
			},
			userId: uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
			want: []task.Task{
				{
					Text:      "test",
					Completed: false,
					UserID:    uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
				},
				{
					Text:      "test 123",
					Completed: true,
					UserID:    uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := repo.NewRepository(config.Config{})

			for _, aTask := range tt.tasks {
				if err := repo.Save(context.Background(), aTask); err != nil {
					t.Error(err)
				}
			}

			s := usecase.NewTaskUseCase(repo)
			allTasks, err := s.GetAllTasksForUser(context.Background(), tt.userId)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("s.GetAllTasksForUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && len(allTasks) != len(tt.want) {
				t.Errorf("s.GetAllTasksForUser() number if tasks doesn't match actual = %v, expected %v", len(allTasks), len(tt.want))
			}
		})
	}
}

func TestTaskUseCaseUpdateTask(t *testing.T) {
	type testCase struct {
		name    string
		task    task.Task
		want    task.Task
		wantErr error
	}

	tests := []testCase{
		{
			name: "Success",
			task: task.Task{
				ID:        uuid.New(),
				Text:      "test",
				Completed: false,
				UserID:    uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
			},
			want: task.Task{
				Text:      "test 321",
				Completed: false,
				UserID:    uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := repo.NewRepository(config.Config{})

			if err := repo.Save(context.Background(), tt.task); err != nil {
				t.Error(err)
			}

			s := usecase.NewTaskUseCase(repo)
			updatedTask, err := s.UpdateTask(context.Background(), tt.task.ID, tt.want.Text, tt.task.UserID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("s.UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && updatedTask.Text != tt.want.Text {
				t.Errorf("s.UpdateTask() Text = %v, want %v", updatedTask.Text, tt.want.Text)
			}
			if err == nil && updatedTask.Completed != tt.want.Completed {
				t.Errorf("s.UpdateTask() Completed = %v, want %v", updatedTask.Completed, tt.want.Completed)
			}
			if err == nil && updatedTask.UpdatedAt.IsZero() {
				t.Error("s.UpdateTask() UpdatedAt is zero")
			}
		})
	}
}

func TestTaskUseCaseMarkTaskCompleted(t *testing.T) {
	type testCase struct {
		name    string
		task    task.Task
		want    task.Task
		wantErr error
	}

	tests := []testCase{
		{
			name: "Success",
			task: task.Task{
				ID:        uuid.New(),
				Text:      "test",
				Completed: false,
				UserID:    uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
			},
			want: task.Task{
				Text:      "test",
				Completed: true,
				UserID:    uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := repo.NewRepository(config.Config{})

			if err := repo.Save(context.Background(), tt.task); err != nil {
				t.Error(err)
			}

			s := usecase.NewTaskUseCase(repo)
			updatedTask, err := s.MarkTaskCompleted(context.Background(), tt.task.ID, tt.task.UserID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("s.MarkTaskCompleted() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && updatedTask.Text != tt.want.Text {
				t.Errorf("s.MarkTaskCompleted() Text = %v, want %v", updatedTask.Text, tt.want.Text)
			}
			if err == nil && updatedTask.Completed != tt.want.Completed {
				t.Errorf("s.MarkTaskCompleted() Completed = %v, want %v", updatedTask.Text, tt.want.Text)
			}
			if err == nil && updatedTask.UpdatedAt.IsZero() {
				t.Error("s.MarkTaskCompleted() UpdatedAt is zero")
			}
		})
	}
}

func TestTaskUseCaseMarkTaskNotCompleted(t *testing.T) {
	type testCase struct {
		name    string
		task    task.Task
		want    task.Task
		wantErr error
	}

	tests := []testCase{
		{
			name: "Success",
			task: task.Task{
				ID:        uuid.New(),
				Text:      "test",
				Completed: true,
				UserID:    uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
			},
			want: task.Task{
				Text:      "test",
				Completed: false,
				UserID:    uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := repo.NewRepository(config.Config{})

			if err := repo.Save(context.Background(), tt.task); err != nil {
				t.Error(err)
			}

			s := usecase.NewTaskUseCase(repo)
			updatedTask, err := s.MarkTaskNotCompleted(context.Background(), tt.task.ID, tt.task.UserID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("s.MarkTaskNotCompleted() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && updatedTask.Text != tt.want.Text {
				t.Errorf("s.MarkTaskNotCompleted() Text = %v, want %v", updatedTask.Text, tt.want.Text)
			}
			if err == nil && updatedTask.Completed != tt.want.Completed {
				t.Errorf("s.MarkTaskNotCompleted() Completed = %v, want %v", updatedTask.Text, tt.want.Text)
			}
			if err == nil && updatedTask.UpdatedAt.IsZero() {
				t.Error("s.MarkTaskNotCompleted() UpdatedAt is zero")
			}
		})
	}
}

func TestTaskUseCaseDeleteTask(t *testing.T) {
	type testCase struct {
		name    string
		task    task.Task
		want    task.Task
		wantErr error
	}

	tests := []testCase{
		{
			name: "Success",
			task: task.Task{
				ID:        uuid.New(),
				Text:      "test",
				Completed: false,
				UserID:    uuid.MustParse("731efa19-5926-4cac-8ca5-2fcb7d8056b1"),
			},
			want:    task.Task{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := repo.NewRepository(config.Config{})

			if err := repo.Save(context.Background(), tt.task); err != nil {
				t.Error(err)
			}

			s := usecase.NewTaskUseCase(repo)
			err := s.DeleteTask(context.Background(), tt.task.ID, tt.task.UserID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("s.DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
			}

			deletedTask, err := repo.GetByID(context.Background(), tt.task.ID)
			if !errors.Is(err, task.ErrTaskNotFound) {
				t.Errorf("s.DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && deletedTask != tt.want {
				t.Errorf("s.DeleteTask() Text = %v, want %v", deletedTask, tt.want)
			}
		})
	}
}
