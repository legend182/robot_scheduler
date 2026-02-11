package impl

import (
	"context"
	"errors"
	dao "robot_scheduler/internal/dao/interfaces"

	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TaskDAOImpl struct {
	db *gorm.DB
}

func NewTaskDAO(db *gorm.DB) dao.TaskDAO {
	return &TaskDAOImpl{db: db}
}

func (d *TaskDAOImpl) Create(ctx context.Context, task *entity.Task) error {
	logger.Info("creating task")

	if err := d.db.WithContext(ctx).Create(task).Error; err != nil {
		logger.Error("failed to create task", zap.Error(err))
		return err
	}

	logger.Info("task created successfully", zap.Uint("id", task.ID))
	return nil
}

func (d *TaskDAOImpl) Update(ctx context.Context, task *entity.Task) error {
	logger.Info("updating task", zap.Uint("id", task.ID))

	result := d.db.WithContext(ctx).Save(task)
	if err := result.Error; err != nil {
		logger.Error("failed to update task", zap.Error(err), zap.Uint("id", task.ID))
		return err
	}

	if result.RowsAffected == 0 {
		logger.Warn("task not found for update", zap.Uint("id", task.ID))
		return errors.New("task not found")
	}

	logger.Info("task updated successfully", zap.Uint("id", task.ID))
	return nil
}

func (d *TaskDAOImpl) Delete(ctx context.Context, id uint) error {
	logger.Info("deleting task", zap.Uint("id", id))

	result := d.db.WithContext(ctx).Delete(&entity.Task{}, id)
	if err := result.Error; err != nil {
		logger.Error("failed to delete task", zap.Error(err), zap.Uint("id", id))
		return err
	}

	if result.RowsAffected == 0 {
		logger.Warn("task not found for deletion", zap.Uint("id", id))
		return errors.New("task not found")
	}

	logger.Info("task deleted successfully", zap.Uint("id", id))
	return nil
}

func (d *TaskDAOImpl) FindByID(ctx context.Context, id uint) (*entity.Task, error) {
	logger.Debug("finding task by id", zap.Uint("id", id))

	var task entity.Task
	err := d.db.WithContext(ctx).Preload("SemanticMap.PCDFile").First(&task, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Debug("task not found", zap.Uint("id", id))
			return nil, nil
		}
		logger.Error("failed to find task by id", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	logger.Debug("task found", zap.Uint("id", id))
	return &task, nil
}

func (d *TaskDAOImpl) FindAll(ctx context.Context) ([]*entity.Task, error) {
	logger.Debug("finding all tasks")

	var tasks []*entity.Task
	err := d.db.WithContext(ctx).Preload("SemanticMap.PCDFile").Find(&tasks).Error
	if err != nil {
		logger.Error("failed to find all tasks", zap.Error(err))
		return nil, err
	}

	logger.Debug("found tasks", zap.Int("count", len(tasks)))
	return tasks, nil
}

// FindPage 分页查询任务
func (d *TaskDAOImpl) FindPage(ctx context.Context, offset, limit int) ([]*entity.Task, int64, error) {
	logger.Debug("finding tasks with pagination", zap.Int("offset", offset), zap.Int("limit", limit))

	var (
		tasks []*entity.Task
		total int64
	)

	db := d.db.WithContext(ctx).Model(&entity.Task{}).Preload("SemanticMap.PCDFile")

	if err := db.Count(&total).Error; err != nil {
		logger.Error("failed to count tasks for pagination", zap.Error(err))
		return nil, 0, err
	}

	if total == 0 {
		return []*entity.Task{}, 0, nil
	}

	if err := db.Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		logger.Error("failed to find tasks with pagination", zap.Error(err))
		return nil, 0, err
	}

	logger.Debug("found tasks with pagination", zap.Int("count", len(tasks)), zap.Int64("total", total))
	return tasks, total, nil
}