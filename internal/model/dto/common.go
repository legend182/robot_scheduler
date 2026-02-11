package dto

import (
	"time"
)

// PageRequest 分页请求
// 使用 omitempty 允许不传分页参数，由业务层设置默认值
type PageRequest struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`         // 页码
	PageSize int `form:"pageSize" json:"pageSize" binding:"omitempty,min=1"` // 每页大小
}

// PageResponse 分页响应
type PageResponse struct {
	Total    int64 `json:"total"`    // 总数
	Pages    int   `json:"pages"`    // 总页数
	Page     int   `json:"page"`     // 当前页
	PageSize int   `json:"pageSize"` // 每页大小
}

// BaseResponse 基础响应
type BaseResponse struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

// TimeRange 时间范围
type TimeRange struct {
	StartTime *time.Time `form:"startTime" json:"startTime"` // 开始时间
	EndTime   *time.Time `form:"endTime" json:"endTime"`     // 结束时间
}
