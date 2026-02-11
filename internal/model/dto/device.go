package dto

import (
	"robot_scheduler/internal/model/entity"
	"time"
)

// DeviceCreateRequest 创建设备请求
type DeviceCreateRequest struct {
	Type      entity.DeviceType  `json:"type" binding:"required"`    // 设备类型
	Company   entity.CompanyType `json:"company" binding:"required"` // 设备厂商
	IP        *string            `json:"ip,omitempty"`               // 设备IP
	Port      int                `json:"port" binding:"required"`    // 设备端口
	UserName  *string            `json:"userName,omitempty"`         // 登录用户名
	Password  *string            `json:"password,omitempty"`         // 登录密码
	ExtraInfo *string            `json:"extraInfo,omitempty"`        // 扩展信息
}

// DeviceUpdateRequest 更新设备请求
type DeviceUpdateRequest struct {
	Type      *entity.DeviceType   `json:"type,omitempty"`      // 设备类型
	Company   *entity.CompanyType  `json:"company,omitempty"`   // 设备厂商
	IP        *string              `json:"ip,omitempty"`        // 设备IP
	Port      *int                 `json:"port,omitempty"`      // 设备端口
	UserName  *string              `json:"userName,omitempty"`  // 登录用户名
	Password  *string              `json:"password,omitempty"`  // 登录密码
	Status    *entity.DeviceStatus `json:"status,omitempty"`    // 设备状态
	ExtraInfo *string              `json:"extraInfo,omitempty"` // 扩展信息
}

// DeviceResponse 设备响应
type DeviceResponse struct {
	ID         uint                 `json:"id"`                  // 设备ID
	Type       entity.DeviceType    `json:"type"`                // 设备类型
	Company    entity.CompanyType   `json:"company"`             // 设备厂商
	IP         *string              `json:"ip,omitempty"`        // 设备IP
	Port       int                  `json:"port"`                // 设备端口
	UserName   *string              `json:"userName,omitempty"`  // 登录用户名
	Status     *entity.DeviceStatus `json:"status"`              // 设备状态
	CreateTime *time.Time           `json:"createTime"`          // 创建时间
	UpdateTime *time.Time           `json:"updateTime"`          // 更新时间
	ExtraInfo  *string              `json:"extraInfo,omitempty"` // 扩展信息
}

// DeviceListResponse 设备列表响应
type DeviceListResponse struct {
	PageResponse
	List []*DeviceResponse `json:"list"` // 设备列表
}

// NewDeviceResponseFromEntity 从实体对象构建设备响应
func NewDeviceResponseFromEntity(d *entity.Device) *DeviceResponse {
	if d == nil {
		return nil
	}
	return &DeviceResponse{
		ID:         d.ID,
		Type:       d.Type,
		Company:    d.Company,
		IP:         d.IP,
		Port:       d.Port,
		UserName:   d.UserName,
		Status:     d.Status,
		CreateTime: &d.CreatedAt,
		UpdateTime: &d.UpdatedAt,
		ExtraInfo:  d.ExtraInfo,
	}
}

// NewDeviceListResponseFromEntities 从实体列表构建设备列表响应
func NewDeviceListResponseFromEntities(list []*entity.Device, page PageResponse) *DeviceListResponse {
	resp := &DeviceListResponse{
		PageResponse: page,
		List:         make([]*DeviceResponse, 0, len(list)),
	}
	for _, d := range list {
		resp.List = append(resp.List, NewDeviceResponseFromEntity(d))
	}
	return resp
}