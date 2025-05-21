package model

import (
	"time"

	"gorm.io/gorm"
)

// User 表示系统中的用户实体
type User struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Email         string         `json:"email" gorm:"uniqueIndex;size:255"`
	Phone         *string        `json:"phone" gorm:"uniqueIndex;size:20;default:null"`
	Username      string         `json:"username" gorm:"uniqueIndex;size:50"`
	Password      string         `json:"-" gorm:"size:255"` // 加密的密码，不暴露给客户端
	FirstName     string         `json:"first_name" gorm:"size:50"`
	LastName      string         `json:"last_name" gorm:"size:50"`
	FullName      string         `json:"full_name" gorm:"-"` // 虚拟字段，不存储在数据库中
	Avatar        *string        `json:"avatar" gorm:"size:255;default:null"`
	Role          string         `json:"role" gorm:"size:20;default:'shopper'"`  // 角色: shopper, admin, staff
	Status        string         `json:"status" gorm:"size:20;default:'active'"` // 状态: active, inactive, suspended
	EmailVerified bool           `json:"email_verified" gorm:"default:false"`
	PhoneVerified bool           `json:"phone_verified" gorm:"default:false"`
	LastLoginAt   *time.Time     `json:"last_login_at" gorm:"default:null"`
	MemberLevel   int            `json:"member_level" gorm:"default:0"` // 会员等级: 0=普通会员，1，2，3 等为更高等级
	Points        int            `json:"points" gorm:"default:0"`       // 积分
	TwoFactorAuth bool           `json:"two_factor_auth" gorm:"default:false"`
	Addresses     []Address      `json:"addresses" gorm:"foreignKey:UserID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// Address 表示用户的收货地址
type Address struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"index"`
	Name         string         `json:"name" gorm:"size:50"` // 收货人姓名
	Phone        string         `json:"phone" gorm:"size:20"`
	Province     string         `json:"province" gorm:"size:50"`
	City         string         `json:"city" gorm:"size:50"`
	District     string         `json:"district" gorm:"size:50"`
	DetailedInfo string         `json:"detailed_info" gorm:"size:255"` // 详细地址
	PostalCode   string         `json:"postal_code" gorm:"size:20"`
	IsDefault    bool           `json:"is_default" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// LoginHistory 表示用户的登录历史
type LoginHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`
	IP        string    `json:"ip" gorm:"size:50"`
	UserAgent string    `json:"user_agent" gorm:"size:255"`
	Location  string    `json:"location" gorm:"size:100"`
	CreatedAt time.Time `json:"created_at"` // 登录时间
}

// BeforeSave 在保存前处理 User
func (u *User) BeforeSave(tx *gorm.DB) error {
	u.FullName = u.FirstName + " " + u.LastName
	return nil
}

// AfterFind 在查询后处理 User
func (u *User) AfterFind(tx *gorm.DB) error {
	u.FullName = u.FirstName + " " + u.LastName
	return nil
}
