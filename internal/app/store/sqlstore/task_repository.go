package sqlstore

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/model"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/store"
)

type TaskRepository struct {
	store *Store
}

func (r *TaskRepository) Create(t *model.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}

	if err := t.BeforeCreate(); err != nil {
		return err
	}

	row := r.store.db.QueryRow(
		"INSERT INTO tasks (id, task_text, task_order, is_done, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		t.ID,
		t.TaskText,
		t.TaskOrder,
		t.IsDone,
		t.UserID,
	)

	return row.Scan(&t.ID)
}

func (r *TaskRepository) GetAllByUser(UserID uuid.UUID) ([]*model.Task, error) {
	rows, err := r.store.db.Query(`
		SELECT id, task_text, task_order, is_done, user_id
		FROM tasks
		WHERE user_id = $1
		ORDER BY task_order ASC;
	`, UserID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*model.Task{}
	for rows.Next() {
		t := &model.Task{}
		err := rows.Scan(&t.ID, &t.TaskText, &t.TaskOrder, &t.IsDone, &t.UserID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TaskRepository) MarkAsDone(TaskID uuid.UUID) (*model.Task, error) {
	_, err := r.store.db.Exec("UPDATE tasks SET is_done = true WHERE id = $1;", TaskID)
	if err != nil {
		return nil, err
	}

	t := &model.Task{}

	row := r.store.db.QueryRow(
		"SELECT id, task_text, task_order, is_done, user_id FROM tasks WHERE id = $1;",
		TaskID,
	)

	if err := row.Scan(&t.ID, &t.TaskText, &t.TaskOrder, &t.IsDone, &t.UserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return t, nil
}

func (r *TaskRepository) MarkAsNotDone(TaskID uuid.UUID) (*model.Task, error) {
	_, err := r.store.db.Exec("UPDATE tasks SET is_done = false WHERE id = $1;", TaskID)
	if err != nil {
		return nil, err
	}

	t := &model.Task{}

	row := r.store.db.QueryRow(
		"SELECT id, task_text, task_order, is_done, user_id FROM tasks WHERE id = $1;",
		TaskID,
	)

	if err := row.Scan(&t.ID, &t.TaskText, &t.TaskOrder, &t.IsDone, &t.UserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return t, nil
}

func (r *TaskRepository) Delete(TaskID uuid.UUID) error {
	_, err := r.store.db.Exec("DELETE FROM tasks WHERE id = $1;", TaskID)
	if err != nil {
		return err
	}

	return nil
}
