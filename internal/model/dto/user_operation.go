package dto

import (
	"robot_scheduler/internal/model/entity"
	"time"
)

// UserOperationCreateRequest 创建操作记录请求
type UserOperationCreateRequest struct {
	UserName   string               `json:"userName" binding:"required"`  // 操作人员
	Operation  entity.OperationType `json:"operation" binding:"required"` // 操作类型
	Module     string               `json:"module" binding:"required"`    // 操作模块
	TargetID   *uint                `json:"targetId,omitempty"`           // 目标ID
	TargetName *string              `json:"targetName,omitempty"`         // 目标名称
	IP         *string              `json:"ip,omitempty"`                 // 操作IP
	UserAgent  *string              `json:"userAgent,omitempty"`          // 用户代理
	ExtraInfo  string               `json:"extraInfo" binding:"required"` // 操作信息
}

// UserOperationUpdateRequest 更新操作记录请求
type UserOperationUpdateRequest struct {
	UserName   *string               `json:"userName,omitempty"`   // 操作人员
	Operation  *entity.OperationType `json:"operation,omitempty"`  // 操作类型
	Module     *string               `json:"module,omitempty"`     // 操作模块
	TargetID   *uint                 `json:"targetId,omitempty"`   // 目标ID
	TargetName *string               `json:"targetName,omitempty"` // 目标名称
	IP         *string               `json:"ip,omitempty"`         // 操作IP
	UserAgent  *string               `json:"userAgent,omitempty"`  // 用户代理
	ExtraInfo  *string               `json:"extraInfo,omitempty"`  // 操作信息
}

// UserOperationResponse 操作记录响应
type UserOperationResponse struct {
	ID         uint                 `json:"id"`                   // 记录ID
	UserName   string               `json:"userName"`             // 操作人员
	Operation  entity.OperationType `json:"operation"`            // 操作类型
	Module     string               `json:"module"`               // 操作模块
	TargetID   *uint                `json:"targetId,omitempty"`   // 目标ID
	TargetName *string              `json:"targetName,omitempty"` // 目标名称
	IP         *string              `json:"ip,omitempty"`         // 操作IP
	UserAgent  *string              `json:"userAgent,omitempty"`  // 用户代理
	ExtraInfo  string               `json:"extraInfo"`            // 操作信息
	CreateTime *time.Time           `json:"createTime"`           // 操作时间
}

// UserOperationListResponse 操作记录列表响应
type UserOperationListResponse struct {
	PageResponse
	List []*UserOperationResponse `json:"list"` // 操作记录列表
}

// NewUserOperationResponseFromEntity 从实体对象构建操作记录响应
func NewUserOperationResponseFromEntity(o *entity.UserOperation) *UserOperationResponse {
	if o == nil {
		return nil
	}
	return &UserOperationResponse{
		ID:         o.ID,
		UserName:   o.UserName,
		Operation:  o.Operation,
		Module:     o.Module,
		TargetID:   o.TargetID,
		TargetName: o.TargetName,
		IP:         o.IP,
		UserAgent:  o.UserAgent,
		ExtraInfo:  o.ExtraInfo,
		CreateTime: o.CreateTime,
	}
}

// NewUserOperationListResponseFromEntities 从实体列表构建操作记录列表响应
func NewUserOperationListResponseFromEntities(list []*entity.UserOperation, page PageResponse) *UserOperationListResponse {
	resp := &UserOperationListResponse{
		PageResponse: page,
		List:         make([]*UserOperationResponse, 0, len(list)),
	}
	for _, o := range list {
		resp.List = append(resp.List, NewUserOperationResponseFromEntity(o))
	}
	return resp
}