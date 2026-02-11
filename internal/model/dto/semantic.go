package dto

import (
	"robot_scheduler/internal/model/entity"
	"time"
)

// SemanticMapCreateRequest 创建语义地图请求
type SemanticMapCreateRequest struct {
	PCDFileID    uint    `json:"pcdFileId" binding:"required"`    // 对应的pcd地图文件id
	UserName     string  `json:"userName" binding:"required"`     // 编辑人员
	SemanticInfo string  `json:"semanticInfo" binding:"required"` // 语义信息
	ExtraInfo    *string `json:"extraInfo,omitempty"`             // 扩展信息
}

// SemanticMapUpdateRequest 更新语义地图请求
type SemanticMapUpdateRequest struct {
	PCDFileID    *uint   `json:"pcdFileId,omitempty"`    // 对应的pcd地图文件id
	UserName     *string `json:"userName,omitempty"`     // 编辑人员
	SemanticInfo *string `json:"semanticInfo,omitempty"` // 语义信息
	ExtraInfo    *string `json:"extraInfo,omitempty"`    // 扩展信息
}

// SemanticMapResponse 语义地图响应
type SemanticMapResponse struct {
	ID           uint            `json:"id"`                  // 语义地图ID
	PCDFileID    uint            `json:"pcdFileId"`           // 对应的pcd地图文件id
	PCDFile      *entity.PCDFile `json:"pcdFile,omitempty"`   // 关联的点云地图
	UserName     string          `json:"userName"`            // 编辑人员
	SemanticInfo string          `json:"semanticInfo"`        // 语义信息
	CreateTime   *time.Time      `json:"createTime"`          // 创建时间
	UpdateTime   *time.Time      `json:"updateTime"`          // 更新时间
	ExtraInfo    *string         `json:"extraInfo,omitempty"` // 扩展信息
}

// SemanticMapListResponse 语义地图列表响应
type SemanticMapListResponse struct {
	PageResponse
	List []*SemanticMapResponse `json:"list"` // 语义地图列表
}

// NewSemanticMapResponseFromEntity 从实体对象构建语义地图响应
func NewSemanticMapResponseFromEntity(m *entity.SemanticMap) *SemanticMapResponse {
	if m == nil {
		return nil
	}
	return &SemanticMapResponse{
		ID:           m.ID,
		PCDFileID:    m.PCDFileID,
		PCDFile:      &m.PCDFile,
		UserName:     m.UserName,
		SemanticInfo: m.SemanticInfo,
		CreateTime:   &m.CreatedAt,
		UpdateTime:   &m.UpdatedAt,
		ExtraInfo:    m.ExtraInfo,
	}
}

// NewSemanticMapListResponseFromEntities 从实体列表构建语义地图列表响应
func NewSemanticMapListResponseFromEntities(list []*entity.SemanticMap, page PageResponse) *SemanticMapListResponse {
	resp := &SemanticMapListResponse{
		PageResponse: page,
		List:         make([]*SemanticMapResponse, 0, len(list)),
	}
	for _, m := range list {
		resp.List = append(resp.List, NewSemanticMapResponseFromEntity(m))
	}
	return resp
}