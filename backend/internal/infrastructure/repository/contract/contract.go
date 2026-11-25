// Package contract has test suites every repository implementation must pass.
package contract

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
)

// TaskRepository runs the contract; newRepo must return an empty repository.
func TaskRepository(t *testing.T, newRepo func(t *testing.T) task.Repository) {
	t.Helper()
	ctx := context.Background()

	t.Run("GetByID_NotFound", func(t *testing.T) {
		_, err := newRepo(t).GetByID(ctx, uuid.New())
		if !errors.Is(err, task.ErrTaskNotFound) {
			t.Fatalf("err = %v, want ErrTaskNotFound", err)
		}
	})

	t.Run("Save_ThenGetByID_RoundTripsAllFields", func(t *testing.T) {
		r := newRepo(t)
		want := mustTask(t, "write tests", uuid.New())
		want.SetCompleted(true)

		if err := r.Save(ctx, want); err != nil {
			t.Fatalf("Save: %v", err)
		}
		got, err := r.GetByID(ctx, want.ID)
		if err != nil {
			t.Fatalf("GetByID: %v", err)
		}
		assertTaskEqual(t, got, want)
	})

	t.Run("Save_Duplicate_ReturnsErrTaskExists", func(t *testing.T) {
		r := newRepo(t)
		tk := mustTask(t, "x", uuid.New())
		if err := r.Save(ctx, tk); err != nil {
			t.Fatalf("Save: %v", err)
		}
		if err := r.Save(ctx, tk); !errors.Is(err, task.ErrTaskExists) {
			t.Fatalf("second Save err = %v, want ErrTaskExists", err)
		}
	})

	t.Run("Update_PersistsTextCompletedAndUpdatedAt", func(t *testing.T) {
		r := newRepo(t)
		tk := mustTask(t, "before", uuid.New())
		if err := r.Save(ctx, tk); err != nil {
			t.Fatalf("Save: %v", err)
		}

		if err := tk.SetText("after"); err != nil {
			t.Fatal(err)
		}
		tk.SetCompleted(true)
		tk.UpdatedAt = tk.CreatedAt.Add(time.Minute)
		if err := r.Update(ctx, tk); err != nil {
			t.Fatalf("Update: %v", err)
		}

		got, err := r.GetByID(ctx, tk.ID)
		if err != nil {
			t.Fatalf("GetByID: %v", err)
		}
		assertTaskEqual(t, got, tk)
		if !got.UpdatedAt.After(got.CreatedAt) {
			t.Errorf("UpdatedAt %v not after CreatedAt %v", got.UpdatedAt, got.CreatedAt)
		}
	})

	t.Run("Update_NotFound_ReturnsErrTaskNotFound", func(t *testing.T) {
		err := newRepo(t).Update(ctx, mustTask(t, "x", uuid.New()))
		if !errors.Is(err, task.ErrTaskNotFound) {
			t.Fatalf("err = %v, want ErrTaskNotFound", err)
		}
	})

	t.Run("Delete_RemovesTask", func(t *testing.T) {
		r := newRepo(t)
		tk := mustTask(t, "x", uuid.New())
		if err := r.Save(ctx, tk); err != nil {
			t.Fatalf("Save: %v", err)
		}
		if err := r.Delete(ctx, tk.ID); err != nil {
			t.Fatalf("Delete: %v", err)
		}
		if _, err := r.GetByID(ctx, tk.ID); !errors.Is(err, task.ErrTaskNotFound) {
			t.Fatalf("GetByID after Delete err = %v, want ErrTaskNotFound", err)
		}
	})

	t.Run("Delete_NotFound_ReturnsErrTaskNotFound", func(t *testing.T) {
		if err := newRepo(t).Delete(ctx, uuid.New()); !errors.Is(err, task.ErrTaskNotFound) {
			t.Fatalf("err = %v, want ErrTaskNotFound", err)
		}
	})

	t.Run("GetAllByUserID_Empty_ReturnsEmptyNonNil", func(t *testing.T) {
		got, err := newRepo(t).GetAllByUserID(ctx, uuid.New())
		if err != nil {
			t.Fatalf("GetAllByUserID: %v", err)
		}
		if got == nil || len(got) != 0 {
			t.Fatalf("got %#v, want empty non-nil slice", got)
		}
	})

	t.Run("GetAllByUserID_OnlyOwnTasks_OrderedByCreatedAt", func(t *testing.T) {
		r := newRepo(t)
		owner, other := uuid.New(), uuid.New()

		base := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
		var want []task.Task
		for i, text := range []string{"first", "second", "third"} {
			tk := mustTask(t, text, owner)
			tk.CreatedAt = base.Add(time.Duration(i) * time.Hour)
			tk.UpdatedAt = tk.CreatedAt
			want = append(want, tk)
		}
		// Save out of order to make sure ordering comes from CreatedAt.
		for _, i := range []int{2, 0, 1} {
			if err := r.Save(ctx, want[i]); err != nil {
				t.Fatalf("Save: %v", err)
			}
		}
		if err := r.Save(ctx, mustTask(t, "not mine", other)); err != nil {
			t.Fatalf("Save: %v", err)
		}

		got, err := r.GetAllByUserID(ctx, owner)
		if err != nil {
			t.Fatalf("GetAllByUserID: %v", err)
		}
		if len(got) != len(want) {
			t.Fatalf("got %d tasks, want %d", len(got), len(want))
		}
		for i := range want {
			assertTaskEqual(t, got[i], want[i])
		}
	})

	t.Run("Concurrent_SaveAndList", func(t *testing.T) {
		r := newRepo(t)
		owner := uuid.New()
		const n = 20

		var wg sync.WaitGroup
		for range n {
			wg.Go(func() {
				if err := r.Save(ctx, mustTask(t, "concurrent", owner)); err != nil {
					t.Errorf("Save: %v", err)
				}
				if _, err := r.GetAllByUserID(ctx, owner); err != nil {
					t.Errorf("GetAllByUserID: %v", err)
				}
			})
		}
		wg.Wait()

		got, err := r.GetAllByUserID(ctx, owner)
		if err != nil {
			t.Fatalf("GetAllByUserID: %v", err)
		}
		if len(got) != n {
			t.Fatalf("got %d tasks, want %d", len(got), n)
		}
	})
}

// UserRepository runs the contract; newRepo must return an empty repository.
func UserRepository(t *testing.T, newRepo func(t *testing.T) user.Repository) {
	t.Helper()
	ctx := context.Background()

	t.Run("GetByID_NotFound", func(t *testing.T) {
		if _, err := newRepo(t).GetByID(ctx, uuid.New()); !errors.Is(err, user.ErrUserNotFound) {
			t.Fatalf("err = %v, want ErrUserNotFound", err)
		}
	})

	t.Run("GetByEmail_NotFound", func(t *testing.T) {
		if _, err := newRepo(t).GetByEmail(ctx, "nobody@example.com"); !errors.Is(err, user.ErrUserNotFound) {
			t.Fatalf("err = %v, want ErrUserNotFound", err)
		}
	})

	t.Run("Save_ThenGetByIDAndEmail_RoundTrip", func(t *testing.T) {
		r := newRepo(t)
		want := mustUser(t, "alice@example.com")
		if err := r.Save(ctx, want); err != nil {
			t.Fatalf("Save: %v", err)
		}

		for name, get := range map[string]func() (user.User, error){
			"GetByID":    func() (user.User, error) { return r.GetByID(ctx, want.ID) },
			"GetByEmail": func() (user.User, error) { return r.GetByEmail(ctx, want.Email) },
		} {
			got, err := get()
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			if got.ID != want.ID || got.Email != want.Email || got.PasswordHash != want.PasswordHash ||
				!got.CreatedAt.Equal(want.CreatedAt) || !got.UpdatedAt.Equal(want.UpdatedAt) {
				t.Fatalf("%s = %+v, want %+v", name, got, want)
			}
		}
	})

	t.Run("Save_DuplicateEmail_ReturnsErrEmailTaken", func(t *testing.T) {
		r := newRepo(t)
		if err := r.Save(ctx, mustUser(t, "bob@example.com")); err != nil {
			t.Fatalf("Save: %v", err)
		}
		if err := r.Save(ctx, mustUser(t, "bob@example.com")); !errors.Is(err, user.ErrEmailTaken) {
			t.Fatalf("second Save err = %v, want ErrEmailTaken", err)
		}
	})

	t.Run("Concurrent_SaveSameEmail_ExactlyOneWins", func(t *testing.T) {
		r := newRepo(t)
		const n = 10

		var (
			wg      sync.WaitGroup
			mu      sync.Mutex
			success int
		)
		for range n {
			wg.Go(func() {
				err := r.Save(ctx, mustUser(t, "race@example.com"))
				switch {
				case err == nil:
					mu.Lock()
					success++
					mu.Unlock()
				case !errors.Is(err, user.ErrEmailTaken):
					t.Errorf("Save err = %v, want nil or ErrEmailTaken", err)
				}
			})
		}
		wg.Wait()

		if success != 1 {
			t.Fatalf("%d saves succeeded, want exactly 1", success)
		}
	})
}

func mustTask(t *testing.T, text string, userID uuid.UUID) task.Task {
	t.Helper()
	tk, err := task.NewTask(text, userID)
	if err != nil {
		t.Fatalf("NewTask: %v", err)
	}
	return tk
}

func mustUser(t *testing.T, email string) user.User {
	t.Helper()
	u, err := user.NewUser(email, "$2a$10$abcdefghijklmnopqrstuuSomeFakeBcryptHashValueForTests")
	if err != nil {
		t.Fatalf("NewUser: %v", err)
	}
	return u
}

func assertTaskEqual(t *testing.T, got, want task.Task) {
	t.Helper()
	if got.ID != want.ID || got.Text != want.Text || got.Completed != want.Completed || got.UserID != want.UserID ||
		!got.CreatedAt.Equal(want.CreatedAt) || !got.UpdatedAt.Equal(want.UpdatedAt) {
		t.Fatalf("task = %+v, want %+v", got, want)
	}
}
