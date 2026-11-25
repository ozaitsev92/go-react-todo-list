package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	taskmemory "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/memory"
	"github.com/ozaitsev92/tododdd/internal/usecase"
)

func setupTask(ctx context.Context, t *testing.T) (*usecase.TaskUseCase, *faultyTasks, task.Task) {
	t.Helper()
	repo := &faultyTasks{Repository: taskmemory.NewRepository()}
	uc := usecase.NewTaskUseCase(repo)
	tk, err := uc.CreateTask(ctx, "original", uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	return uc, repo, tk
}

func TestTaskUseCase_CreateTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("invalid text not saved", func(t *testing.T) {
		t.Parallel()
		repo := taskmemory.NewRepository()
		owner := uuid.New()
		if _, err := usecase.NewTaskUseCase(repo).CreateTask(ctx, " ", owner); !errors.Is(err, task.ErrInvalidText) {
			t.Fatalf("err = %v", err)
		}
		if all, _ := repo.GetAllByUserID(ctx, owner); len(all) != 0 {
			t.Fatal("invalid task was saved")
		}
	})

	t.Run("repository error is wrapped", func(t *testing.T) {
		t.Parallel()
		uc := usecase.NewTaskUseCase(&faultyTasks{Repository: taskmemory.NewRepository(), failSave: true})
		if _, err := uc.CreateTask(ctx, "x", uuid.New()); !errors.Is(err, errDB) {
			t.Fatalf("err = %v, want wrapped errDB", err)
		}
	})
}

// Every operation on an existing task must enforce ownership, report missing
// tasks, and propagate repository failures.
func TestTaskUseCase_OwnedOperations(t *testing.T) {
	t.Parallel()

	ops := map[string]func(uc *usecase.TaskUseCase, id, userID uuid.UUID) error{
		"GetTask": func(uc *usecase.TaskUseCase, id, uid uuid.UUID) error {
			_, err := uc.GetTask(context.Background(), id, uid)
			return err
		},
		"UpdateTask": func(uc *usecase.TaskUseCase, id, uid uuid.UUID) error {
			_, err := uc.UpdateTask(context.Background(), id, "changed", uid)
			return err
		},
		"SetTaskCompleted": func(uc *usecase.TaskUseCase, id, uid uuid.UUID) error {
			_, err := uc.SetTaskCompleted(context.Background(), id, true, uid)
			return err
		},
		"DeleteTask": func(uc *usecase.TaskUseCase, id, uid uuid.UUID) error {
			return uc.DeleteTask(context.Background(), id, uid)
		},
	}

	for name, op := range ops {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			t.Run("owner succeeds", func(t *testing.T) {
				t.Parallel()
				uc, _, tk := setupTask(t.Context(), t)
				if err := op(uc, tk.ID, tk.UserID); err != nil {
					t.Fatalf("err = %v", err)
				}
			})

			t.Run("other user is rejected and task untouched", func(t *testing.T) {
				t.Parallel()
				uc, repo, tk := setupTask(t.Context(), t)
				if err := op(uc, tk.ID, uuid.New()); !errors.Is(err, usecase.ErrUnauthorizedAction) {
					t.Fatalf("err = %v, want ErrUnauthorizedAction", err)
				}
				got, err := repo.GetByID(context.Background(), tk.ID)
				if err != nil || got != tk {
					t.Fatalf("task modified: %+v, %v", got, err)
				}
			})

			t.Run("missing task", func(t *testing.T) {
				t.Parallel()
				uc, _, tk := setupTask(t.Context(), t)
				if err := op(uc, uuid.New(), tk.UserID); !errors.Is(err, task.ErrTaskNotFound) {
					t.Fatalf("err = %v, want ErrTaskNotFound", err)
				}
			})

			t.Run("lookup failure propagates", func(t *testing.T) {
				t.Parallel()
				uc, repo, tk := setupTask(t.Context(), t)
				repo.failGet = true
				if err := op(uc, tk.ID, tk.UserID); !errors.Is(err, errDB) {
					t.Fatalf("err = %v, want errDB", err)
				}
			})
		})
	}
}

func TestTaskUseCase_UpdateTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("changes text", func(t *testing.T) {
		t.Parallel()
		uc, repo, tk := setupTask(ctx, t)
		got, err := uc.UpdateTask(ctx, tk.ID, "changed", tk.UserID)
		if err != nil || got.Text != "changed" {
			t.Fatalf("got %+v, %v", got, err)
		}
		if stored, _ := repo.GetByID(ctx, tk.ID); stored.Text != "changed" {
			t.Fatalf("not persisted: %+v", stored)
		}
	})

	t.Run("invalid text is not written", func(t *testing.T) {
		t.Parallel()
		uc, repo, tk := setupTask(ctx, t)
		repo.failUpdate = true // would surface if Update were called
		if _, err := uc.UpdateTask(ctx, tk.ID, "", tk.UserID); !errors.Is(err, task.ErrInvalidText) {
			t.Fatalf("err = %v, want ErrInvalidText", err)
		}
	})

	t.Run("update failure propagates", func(t *testing.T) {
		t.Parallel()
		uc, repo, tk := setupTask(ctx, t)
		repo.failUpdate = true
		if _, err := uc.UpdateTask(ctx, tk.ID, "changed", tk.UserID); !errors.Is(err, errDB) {
			t.Fatalf("err = %v, want errDB", err)
		}
	})
}

func TestTaskUseCase_SetTaskCompleted_Toggles(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	uc, repo, tk := setupTask(t.Context(), t)

	for _, want := range []bool{true, false, true} {
		got, err := uc.SetTaskCompleted(ctx, tk.ID, want, tk.UserID)
		if err != nil || got.Completed != want {
			t.Fatalf("SetTaskCompleted(%v) = %+v, %v", want, got, err)
		}
		if stored, _ := repo.GetByID(ctx, tk.ID); stored.Completed != want {
			t.Fatalf("not persisted: %+v", stored)
		}
	}
}

func TestTaskUseCase_DeleteTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("removes task", func(t *testing.T) {
		t.Parallel()
		uc, repo, tk := setupTask(ctx, t)
		if err := uc.DeleteTask(ctx, tk.ID, tk.UserID); err != nil {
			t.Fatal(err)
		}
		if _, err := repo.GetByID(ctx, tk.ID); !errors.Is(err, task.ErrTaskNotFound) {
			t.Fatalf("err = %v, want ErrTaskNotFound", err)
		}
	})

	t.Run("delete failure propagates", func(t *testing.T) {
		t.Parallel()
		uc, repo, tk := setupTask(ctx, t)
		repo.failDelete = true
		if err := uc.DeleteTask(ctx, tk.ID, tk.UserID); !errors.Is(err, errDB) {
			t.Fatalf("err = %v, want errDB", err)
		}
	})
}
