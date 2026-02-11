package dto

import (
	"robot_scheduler/internal/model/entity"
	"time"
)

// PCDFileCreateRequest 创建点云地图请求
type PCDFileCreateRequest struct {
	Name      string  `json:"name" binding:"required,min=1,max=100"` // 地图名称
	Area      string  `json:"area" binding:"required"`               // 区域描述
	Path      string  `json:"path" binding:"required"`               // 文件存储路径
	UserName  string  `json:"userName" binding:"required"`           // 上传人员
	Size      int     `json:"size" binding:"required,min=0"`         // 文件大小
	MinioPath *string `json:"minioPath,omitempty"`                   // MinIO存储路径
	ExtraInfo *string `json:"extraInfo,omitempty"`                   // 扩展信息
}

// PCDFileUpdateRequest 更新点云地图请求
type PCDFileUpdateRequest struct {
	Name      *string `json:"name,omitempty" binding:"omitempty,min=1,max=100"` // 地图名称
	Area      *string `json:"area,omitempty"`                                   // 区域描述
	Path      *string `json:"path,omitempty"`                                   // 文件存储路径
	UserName  *string `json:"userName,omitempty"`                               // 上传人员
	Size      *int    `json:"size,omitempty" binding:"omitempty,min=0"`         // 文件大小
	MinioPath *string `json:"minioPath,omitempty"`                              // MinIO存储路径
	ExtraInfo *string `json:"extraInfo,omitempty"`                              // 扩展信息
}

// PCDFileResponse 点云地图响应
type PCDFileResponse struct {
	ID         uint       `json:"id"`                  // 地图ID
	Name       string     `json:"name"`                // 地图名称
	Area       string     `json:"area"`                // 区域描述
	Path       string     `json:"path"`                // 文件存储路径
	UserName   string     `json:"userName"`            // 上传人员
	Size       int        `json:"size"`                // 文件大小
	MinioPath  *string    `json:"minioPath,omitempty"` // MinIO存储路径
	CreateTime *time.Time `json:"createTime"`          // 创建时间
	UpdateTime *time.Time `json:"updateTime"`          // 更新时间
	ExtraInfo  *string    `json:"extraInfo,omitempty"` // 扩展信息
}

// PCDFileListResponse 点云地图列表响应
type PCDFileListResponse struct {
	PageResponse
	List []*PCDFileResponse `json:"list"` // 点云地图列表
}

// NewPCDFileResponseFromEntity 从实体对象构建点云地图响应
func NewPCDFileResponseFromEntity(f *entity.PCDFile) *PCDFileResponse {
	if f == nil {
		return nil
	}
	return &PCDFileResponse{
		ID:         f.ID,
		Name:       f.Name,
		Area:       f.Area,
		Path:       f.Path,
		UserName:   f.UserName,
		Size:       f.Size,
		MinioPath:  f.MinioPath,
		CreateTime: &f.CreatedAt,
		UpdateTime: &f.UpdatedAt,
		ExtraInfo:  f.ExtraInfo,
	}
}

// NewPCDFileListResponseFromEntities 从实体列表构建点云地图列表响应
func NewPCDFileListResponseFromEntities(list []*entity.PCDFile, page PageResponse) *PCDFileListResponse {
	resp := &PCDFileListResponse{
		PageResponse: page,
		List:         make([]*PCDFileResponse, 0, len(list)),
	}
	for _, f := range list {
		resp.List = append(resp.List, NewPCDFileResponseFromEntity(f))
	}
	return resp
}
