package service

import (
	"context"
	"errors"
	dao "robot_scheduler/internal/dao/interfaces"
	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/dto"
	"robot_scheduler/internal/model/entity"

	"go.uber.org/zap"
)

// TaskService 任务服务
type TaskService struct {
	taskDAO dao.TaskDAO
}

func NewTaskService(taskDAO dao.TaskDAO) *TaskService {
	return &TaskService{
		taskDAO: taskDAO,
	}
}

// CreateTask 创建任务
func (s *TaskService) CreateTask(ctx context.Context, req *dto.TaskCreateRequest) (*dto.TaskResponse, error) {
	logger.Info("creating task in service", zap.Uint("semanticMapID", req.SemanticMapID))

	// 创建任务实体
	task := &entity.Task{
		SemanticMapID: req.SemanticMapID,
		UserName:      req.UserName,
		TaskInfo:      req.TaskInfo,
		Status:        &[]entity.TaskStatus{entity.TaskStatusPending}[0],
		ExtraInfo:     req.ExtraInfo,
	}

	// 保存到数据库
	if err := s.taskDAO.Create(ctx, task); err != nil {
		logger.Error("failed to create task in service", zap.Error(err))
		return nil, err
	}

	logger.Info("task created successfully in service", zap.Uint("id", task.ID))
	return dto.NewTaskResponseFromEntity(task), nil
}

// UpdateTask 更新任务
func (s *TaskService) UpdateTask(ctx context.Context, id uint, req *dto.TaskUpdateRequest) error {
	logger.Info("updating task in service", zap.Uint("id", id))

	// 获取任务
	task, err := s.taskDAO.FindByID(ctx, id)
	if err != nil {
		logger.Error("failed to find task for update", zap.Error(err), zap.Uint("id", id))
		return err
	}

	if task == nil {
		logger.Warn("task not found for update", zap.Uint("id", id))
		return errors.New("task not found")
	}

	// 更新字段
	if req.SemanticMapID != nil {
		task.SemanticMapID = *req.SemanticMapID
	}
	if req.UserName != nil {
		task.UserName = *req.UserName
	}
	if req.TaskInfo != nil {
		task.TaskInfo = *req.TaskInfo
	}
	if req.Status != nil {
		task.Status = req.Status
	}
	if req.ExtraInfo != nil {
		task.ExtraInfo = req.ExtraInfo
	}

	// 保存更新
	if err := s.taskDAO.Update(ctx, task); err != nil {
		logger.Error("failed to update task in service", zap.Error(err), zap.Uint("id", id))
		return err
	}

	logger.Info("task updated successfully in service", zap.Uint("id", id))
	return nil
}

// DeleteTask 删除任务
func (s *TaskService) DeleteTask(ctx context.Context, id uint) error {
	logger.Info("deleting task in service", zap.Uint("id", id))
	return s.taskDAO.Delete(ctx, id)
}

// GetTaskByID 根据ID获取任务
func (s *TaskService) GetTaskByID(ctx context.Context, id uint) (*dto.TaskResponse, error) {
	logger.Debug("getting task by id in service", zap.Uint("id", id))
	task, err := s.taskDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, nil
	}
	return dto.NewTaskResponseFromEntity(task), nil
}

// ListTasks 分页获取任务列表
func (s *TaskService) ListTasks(ctx context.Context, req dto.PageRequest) (*dto.TaskListResponse, error) {
	logger.Debug("listing tasks in service with pagination", zap.Int("page", req.Page), zap.Int("pageSize", req.PageSize))

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize

	tasks, total, err := s.taskDAO.FindPage(ctx, offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	pages := 0
	if req.PageSize > 0 {
		pages = int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	}

	page := dto.PageResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Pages:    pages,
	}

	return dto.NewTaskListResponseFromEntities(tasks, page), nil
}
