package entity

import "gorm.io/gorm"

// SemanticMap 语义地图表
type SemanticMap struct {
	gorm.Model
	PCDFileID    uint    `gorm:"not null;comment:对应的pcd地图文件id;index"`
	PCDFile      PCDFile `gorm:"foreignKey:PCDFileID"`
	UserName     string  `gorm:"type:text;not null;comment:编辑人员"`
	SemanticInfo string  `gorm:"type:text;comment:语义信息"`
	ExtraInfo    *string `gorm:"type:text;comment:扩展信息(JSON)"`
}

func (SemanticMap) TableName() string {
	return "semantic_map"
}
