package task_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
)

func TestNewTask(t *testing.T) {
	t.Parallel()
	owner := uuid.New()

	tests := []struct {
		name     string
		text     string
		userID   uuid.UUID
		wantText string
		wantErr  error
	}{
		{name: "valid", text: "buy milk", userID: owner, wantText: "buy milk"},
		{name: "trims whitespace", text: "  buy milk \n", userID: owner, wantText: "buy milk"},
		{name: "max length", text: strings.Repeat("ж", task.MaxTextLength), userID: owner, wantText: strings.Repeat("ж", task.MaxTextLength)},
		{name: "empty", text: "", userID: owner, wantErr: task.ErrInvalidText},
		{name: "whitespace only", text: " \t\n", userID: owner, wantErr: task.ErrInvalidText},
		{name: "too long", text: strings.Repeat("a", task.MaxTextLength+1), userID: owner, wantErr: task.ErrInvalidText},
		{name: "nil owner", text: "x", userID: uuid.Nil, wantErr: task.ErrInvalidUserID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := task.NewTask(tt.text, tt.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if got.ID == uuid.Nil || got.Text != tt.wantText || got.Completed || got.UserID != tt.userID {
				t.Fatalf("unexpected task %+v", got)
			}
			if got.CreatedAt.Location() != time.UTC || !got.CreatedAt.Equal(got.UpdatedAt) {
				t.Fatalf("timestamps not UTC/equal: %v %v", got.CreatedAt, got.UpdatedAt)
			}
			if got.CreatedAt.Nanosecond()%int(time.Millisecond) != 0 {
				t.Fatalf("CreatedAt not truncated to milliseconds: %v", got.CreatedAt)
			}
		})
	}
}

func TestTask_SetText(t *testing.T) {
	t.Parallel()

	t.Run("updates text and UpdatedAt", func(t *testing.T) {
		t.Parallel()
		tk := task.Task{Text: "old"}
		if err := tk.SetText(" new "); err != nil {
			t.Fatal(err)
		}
		if tk.Text != "new" || tk.UpdatedAt.IsZero() {
			t.Fatalf("got %+v", tk)
		}
	})

	t.Run("invalid text leaves task unchanged", func(t *testing.T) {
		t.Parallel()
		before := task.Task{Text: "old"}
		tk := before
		if err := tk.SetText("   "); !errors.Is(err, task.ErrInvalidText) {
			t.Fatalf("err = %v, want ErrInvalidText", err)
		}
		if tk != before {
			t.Fatalf("task changed on error: %+v", tk)
		}
	})
}

func TestTask_SetCompleted(t *testing.T) {
	t.Parallel()
	for _, completed := range []bool{true, false} {
		tk := task.Task{Completed: !completed}
		tk.SetCompleted(completed)
		if tk.Completed != completed || tk.UpdatedAt.IsZero() {
			t.Fatalf("SetCompleted(%v) = %+v", completed, tk)
		}
	}
}

func TestTask_IsOwnedBy(t *testing.T) {
	t.Parallel()
	owner := uuid.New()
	tk := task.Task{UserID: owner}
	if !tk.IsOwnedBy(owner) || tk.IsOwnedBy(uuid.New()) {
		t.Fatal("IsOwnedBy returned the wrong answer")
	}
}
