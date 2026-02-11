package entity

import "gorm.io/gorm"

// PCDFile 点云地图表
type PCDFile struct {
	gorm.Model
	Name      string  `gorm:"type:text;not null;comment:地图名称"`
	Area      string  `gorm:"type:text;not null;comment:区域描述"`
	Path      string  `gorm:"type:text;not null;comment:文件存储路径"`
	UserName  string  `gorm:"type:text;not null;comment:上传人员"`
	Size      int     `gorm:"comment:文件大小(字节)"`
	MinioPath *string `gorm:"type:text;comment:MinIO存储路径"`
	ExtraInfo *string `gorm:"type:text;comment:扩展信息(JSON)"`
}

func (PCDFile) TableName() string {
	return "pcd_file"
}
