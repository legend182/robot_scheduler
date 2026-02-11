package dto

import (
	"robot_scheduler/internal/model/entity"
	"time"
)

// TaskCreateRequest 创建任务请求
type TaskCreateRequest struct {
	SemanticMapID uint    `json:"semanticMapId" binding:"required"` // 对应的语义地图id
	UserName      string  `json:"userName" binding:"required"`      // 编辑人员
	TaskInfo      string  `json:"taskInfo" binding:"required"`      // 任务信息
	ExtraInfo     *string `json:"extraInfo,omitempty"`              // 扩展信息
}

// TaskUpdateRequest 更新任务请求
type TaskUpdateRequest struct {
	SemanticMapID *uint              `json:"semanticMapId,omitempty"` // 对应的语义地图id
	UserName      *string            `json:"userName,omitempty"`      // 编辑人员
	TaskInfo      *string            `json:"taskInfo,omitempty"`      // 任务信息
	Status        *entity.TaskStatus `json:"status,omitempty"`        // 任务状态
	ExtraInfo     *string            `json:"extraInfo,omitempty"`     // 扩展信息
}

// TaskResponse 任务响应
type TaskResponse struct {
	ID            uint                `json:"id"`                    // 任务ID
	SemanticMapID uint                `json:"semanticMapId"`         // 对应的语义地图id
	SemanticMap   *entity.SemanticMap `json:"semanticMap,omitempty"` // 关联的语义地图
	UserName      string              `json:"userName"`              // 编辑人员
	TaskInfo      string              `json:"taskInfo"`              // 任务信息
	Status        *entity.TaskStatus  `json:"status"`                // 任务状态
	CreateTime    *time.Time          `json:"createTime"`            // 创建时间
	UpdateTime    *time.Time          `json:"updateTime"`            // 更新时间
	ExtraInfo     *string             `json:"extraInfo,omitempty"`   // 扩展信息
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	PageResponse
	List []*TaskResponse `json:"list"` // 任务列表
}

// NewTaskResponseFromEntity 从实体对象构建任务响应
func NewTaskResponseFromEntity(t *entity.Task) *TaskResponse {
	if t == nil {
		return nil
	}
	return &TaskResponse{
		ID:            t.ID,
		SemanticMapID: t.SemanticMapID,
		SemanticMap:   &t.SemanticMap,
		UserName:      t.UserName,
		TaskInfo:      t.TaskInfo,
		Status:        t.Status,
		CreateTime:    &t.CreatedAt,
		UpdateTime:    &t.UpdatedAt,
		ExtraInfo:     t.ExtraInfo,
	}
}

// NewTaskListResponseFromEntities 从实体列表构建任务列表响应
func NewTaskListResponseFromEntities(list []*entity.Task, page PageResponse) *TaskListResponse {
	resp := &TaskListResponse{
		PageResponse: page,
		List:         make([]*TaskResponse, 0, len(list)),
	}
	for _, t := range list {
		resp.List = append(resp.List, NewTaskResponseFromEntity(t))
	}
	return resp
}