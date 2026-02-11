package entity

import "gorm.io/gorm"

// DeviceType 设备类型枚举
type DeviceType string

const (
	DeviceTypeWheelRobot DeviceType = "robot_wheel" // 轮式机器人
	DeviceTypeBipedRobot DeviceType = "robot_biped" // 双足机器人
)

// CompanyType 厂商类型枚举
type CompanyType string

const (
	CompanyCyborg CompanyType = "cyborg" // 赛博格
)

// Device 设备表
type Device struct {
	gorm.Model
	Type      DeviceType    `gorm:"type:text;not null;comment:设备类型"`
	Company   CompanyType   `gorm:"type:text;not null;comment:设备厂商"`
	IP        *string       `gorm:"type:text;comment:设备IP"`
	Port      int           `gorm:"comment:设备端口"`
	UserName  *string       `gorm:"type:text;comment:登录用户名"`
	Password  *string       `gorm:"type:text;comment:登录密码(RSA加密)"`
	Status    *DeviceStatus `gorm:"type:text;default:'offline';comment:设备状态"`
	ExtraInfo *string       `gorm:"type:text;comment:扩展信息(JSON)"`
}

// DeviceStatus 设备状态枚举
type DeviceStatus string

const (
	DeviceStatusOffline DeviceStatus = "offline" // 离线
	DeviceStatusOnline  DeviceStatus = "online"  // 在线
	DeviceStatusBusy    DeviceStatus = "busy"    // 忙碌
	DeviceStatusError   DeviceStatus = "error"   // 错误
)

func (Device) TableName() string {
	return "device"
}
