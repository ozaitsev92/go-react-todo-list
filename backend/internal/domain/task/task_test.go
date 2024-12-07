package task_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
)

func TestTaskNewTask(t *testing.T) {
	type args struct {
		task   string
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
				task:   "test",
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
				task:   "",
				userId: uuid.New(),
			},
			want:    task.Task{},
			wantErr: task.ErrInvalidText,
		},
		{
			name: "Empty userId",
			args: args{
				task:   "test",
				userId: uuid.Nil,
			},
			want:    task.Task{},
			wantErr: task.ErrInvalidUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := task.NewTask(tt.args.task, tt.args.userId)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got.Text != tt.want.Text {
				t.Errorf("NewTask() Text = %v, want %v", got.Text, tt.want.Text)
			}
			if got.Completed != tt.want.Completed {
				t.Errorf("NewTask() Completed = %v, want %v", got.Completed, tt.want.Completed)
			}
			if got.UserID != tt.want.UserID {
				t.Errorf("NewTask() UserID = %v, want %v", got.UserID, tt.want.UserID)
			}
			if tt.wantErr == nil && got.CreatedAt.IsZero() {
				t.Error("NewTask() CreatedAt is zero")
			}
			if tt.wantErr == nil && got.UpdatedAt.IsZero() {
				t.Error("NewTask() CreatedAt is UpdatedAt")
			}
			if tt.wantErr == nil && got.ID == uuid.Nil {
				t.Error("NewTask() ID is nil")
			}
		})
	}
}

func TestTaskSetTask(t *testing.T) {
	type fields struct {
		Text string
	}

	type args struct {
		text string
	}

	type testCase struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}

	tests := []testCase{
		{
			name: "Success",
			fields: fields{
				Text: "test",
			},
			args: args{
				text: "new test",
			},
			wantErr: nil,
		},
		{
			name: "Empty task",
			fields: fields{
				Text: "test",
			},
			args: args{
				text: "",
			},
			wantErr: task.ErrInvalidText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ti := task.Task{
				Text: tt.fields.Text,
			}
			err := ti.SetText(tt.args.text)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SetText() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && ti.Text != tt.args.text {
				t.Errorf("SetText() Text = %v, want %v", ti.Text, tt.args.text)
			}
			if err == nil && ti.UpdatedAt.IsZero() {
				t.Error("SetText() UpdatedAt is zero")
			}
		})
	}
}

func TestTaskMarkCompleted(t *testing.T) {
	type testCase struct {
		name string
	}

	tests := []testCase{
		{
			name: "Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ti := task.Task{}
			ti.MarkCompleted()
			if !ti.Completed {
				t.Errorf("MarkCompleted() Completed = %v, want %v", ti.Completed, true)
			}
			if ti.UpdatedAt.IsZero() {
				t.Error("MarkCompleted() UpdatedAt is zero")
			}
		})
	}
}

func TestTaskMarkNotCompleted(t *testing.T) {
	type testCase struct {
		name string
	}

	tests := []testCase{
		{
			name: "Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ti := task.Task{}
			ti.MarkNotCompleted()
			if ti.Completed {
				t.Errorf("MarkNotCompleted() Completed = %v, want %v", ti.Completed, false)
			}
			if ti.UpdatedAt.IsZero() {
				t.Error("MarkNotCompleted() UpdatedAt is zero")
			}
		})
	}
}
