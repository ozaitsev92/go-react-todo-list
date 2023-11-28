package sqlstore

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
)

type TaskRepository struct {
	store *Store
}

type DBTaskRecord struct {
	ID        uuid.UUID
	TaskText  string
	TaskOrder int
	IsDone    bool
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *DBTaskRecord) ToTask() (*domain.Task, error) {
	t := &domain.Task{}

	if err := t.SetID(r.ID); err != nil {
		return nil, err
	}

	if err := t.SetTaskText(r.TaskText); err != nil {
		return nil, err
	}

	if err := t.SetTaskOrder(r.TaskOrder); err != nil {
		return nil, err
	}

	if err := t.SetIsDone(r.IsDone); err != nil {
		return nil, err
	}

	if err := t.SetUserID(r.UserID); err != nil {
		return nil, err
	}

	if err := t.SetCreatedAt(r.CreatedAt); err != nil {
		return nil, err
	}

	if err := t.SetUpdatedAt(r.UpdatedAt); err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TaskRepository) SaveTask(ctx context.Context, task *domain.Task) error {
	_, err := r.store.db.Exec(
		`
			INSERT INTO tasks (id, task_text, task_order, is_done, user_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT(id) DO UPDATE
			SET task_text = EXCLUDED.task_text,
				task_order = EXCLUDED.task_order,
				is_done = EXCLUDED.is_done,
				updated_at = EXCLUDED.updated_at;
		`,
		task.GetID(),
		task.GetTaskText(),
		task.GetTaskOrder(),
		task.GetIsDone(),
		task.GetUserID(),
		task.GetCreatedAt(),
		task.GetUpdatedAt(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	row := r.store.db.QueryRow(`
		SELECT id, task_text, task_order, is_done, user_id, created_at, updated_at
		FROM tasks
		WHERE id = $1;
	`, id)

	rec := &DBTaskRecord{}
	err := row.Scan(
		&rec.ID,
		&rec.TaskText,
		&rec.TaskOrder,
		&rec.IsDone,
		&rec.UserID,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return rec.ToTask()
}

func (r *TaskRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Task, error) {
	rows, err := r.store.db.Query(`
		SELECT id, task_text, task_order, is_done, user_id, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY task_order ASC;
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*domain.Task{}
	for rows.Next() {
		rec := &DBTaskRecord{}
		err := rows.Scan(
			&rec.ID,
			&rec.TaskText,
			&rec.TaskOrder,
			&rec.IsDone,
			&rec.UserID,
			&rec.CreatedAt,
			&rec.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if task, err := rec.ToTask(); err == nil {
			tasks = append(tasks, task)
		} else {
			return []*domain.Task{}, nil
		}
	}

	return tasks, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, task *domain.Task) error {
	_, err := r.store.db.Exec("DELETE FROM tasks WHERE id = $1;", task.GetID())
	if err != nil {
		return err
	}

	return nil
}
