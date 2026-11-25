package repository_test

import (
	"testing"

	"github.com/ozaitsev92/tododdd/internal/domain/task"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/repository/contract"
	repository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/memory"
)

func TestTaskRepository_Contract(t *testing.T) {
	contract.TaskRepository(t, func(*testing.T) task.Repository { return repository.NewRepository() })
}
