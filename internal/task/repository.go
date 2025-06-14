package task

import (
	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task Task) (Task, error)
	GetAllTasks() ([]Task, error)
	GetTasksByUserID(userID uint32) ([]Task, error)
	UpdateTaskByID(id uint32, task Task) (Task, error)
	DeleteTaskByID(id uint32) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task Task) (Task, error) {
	result := r.db.Create(&task)
	if result.Error != nil {
		return Task{}, result.Error
	}
	return task, nil
}

func (r *taskRepository) GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetTasksByUserID(userID uint32) ([]Task, error) {
	var tasks []Task
	err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) UpdateTaskByID(id uint32, task Task) (Task, error) {
	var existingTask Task
	if err := r.db.First(&existingTask, id).Error; err != nil {
		return Task{}, err
	}

	existingTask.Title = task.Title
	existingTask.IsDone = task.IsDone

	result := r.db.Save(&existingTask)
	if result.Error != nil {
		return Task{}, result.Error
	}

	return existingTask, nil
}

func (r *taskRepository) DeleteTaskByID(id uint32) error {
	result := r.db.Delete(&Task{}, id)
	return result.Error
}
