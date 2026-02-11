package dto

import (
	"robot_scheduler/internal/model/entity"
	"time"
)

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	UserName  string          `json:"userName" binding:"required,min=3,max=50"`                          // 用户名
	Password  string          `json:"password" binding:"required,min=6"`                                 // 密码
	Role      entity.RoleType `json:"role" binding:"required,oneof=administrator manager operator user"` // 角色
	ExtraInfo *string         `json:"extraInfo,omitempty"`                                               // 扩展信息
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Password  *string          `json:"password,omitempty"`                     // 密码
	Role      *entity.RoleType `json:"role,omitempty"`                         // 角色
	IsLocked  *int             `json:"isLocked,omitempty" binding:"oneof=0 1"` // 是否锁定
	ExtraInfo *string          `json:"extraInfo,omitempty"`                    // 扩展信息
}

// UserResponse 用户响应
type UserResponse struct {
	ID         uint            `json:"id"`                  // 用户ID
	UserName   string          `json:"userName"`            // 用户名
	Role       entity.RoleType `json:"role"`                // 角色
	IsLocked   int             `json:"isLocked"`            // 是否锁定
	CreateTime *time.Time      `json:"createTime"`          // 创建时间
	UpdateTime *time.Time      `json:"updateTime"`          // 更新时间
	ExtraInfo  *string         `json:"extraInfo,omitempty"` // 扩展信息
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	PageResponse
	List []*UserResponse `json:"list"` // 用户列表
}

// NewUserResponseFromEntity 从实体对象构建用户响应
func NewUserResponseFromEntity(u *entity.UserInfo) *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
		ID:         u.ID,
		UserName:   u.UserName,
		Role:       u.Role,
		IsLocked:   u.IsLocked,
		CreateTime: &u.CreatedAt,
		UpdateTime: &u.UpdatedAt,
		ExtraInfo:  u.ExtraInfo,
	}
}

// NewUserListResponseFromEntities 从实体列表构建用户列表响应
// 分页信息由调用方填充 PageResponse 字段
func NewUserListResponseFromEntities(list []*entity.UserInfo, page PageResponse) *UserListResponse {
	resp := &UserListResponse{
		PageResponse: page,
		List:         make([]*UserResponse, 0, len(list)),
	}
	for _, u := range list {
		resp.List = append(resp.List, NewUserResponseFromEntity(u))
	}
	return resp
}

// LoginRequest 登录请求
type LoginRequest struct {
	UserName string `json:"userName" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string       `json:"token"` // JWT token
	User  *UserResponse `json:"user"`  // 用户信息
}