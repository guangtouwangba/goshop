package model

import (
	"time"

	"gorm.io/gorm"
)

// TokenType 表示令牌类型
type TokenType string

const (
	// TokenTypeAccess 访问令牌
	TokenTypeAccess TokenType = "access"
	// TokenTypeRefresh 刷新令牌
	TokenTypeRefresh TokenType = "refresh"
	// TokenTypeEmailVerification 邮箱验证令牌
	TokenTypeEmailVerification TokenType = "email_verification"
	// TokenTypePasswordReset 密码重置令牌
	TokenTypePasswordReset TokenType = "password_reset"
	// TokenTypePhoneVerification 手机验证令牌
	TokenTypePhoneVerification TokenType = "phone_verification"
)

// Token 表示认证令牌
type Token struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	Token     string    `json:"token" gorm:"size:255;uniqueIndex;not null"`
	Type      TokenType `json:"type" gorm:"size:30;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsRevoked bool      `json:"is_revoked" gorm:"default:false"`
	IP        *string   `json:"ip" gorm:"size:50"`
	UserAgent *string   `json:"user_agent" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Role 表示用户角色
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;uniqueIndex;not null"`
	Description string         `json:"description" gorm:"size:255"`
	Permissions []Permission   `json:"permissions" gorm:"many2many:role_permissions"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Permission 表示权限
type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;uniqueIndex;not null"`
	Code        string         `json:"code" gorm:"size:50;uniqueIndex;not null"`
	Description string         `json:"description" gorm:"size:255"`
	Module      string         `json:"module" gorm:"size:50;not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserRole 表示用户角色关联
type UserRole struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"uniqueIndex:idx_user_role;not null"`
	RoleID    uint      `json:"role_id" gorm:"uniqueIndex:idx_user_role;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// APIKey 表示API密钥
type APIKey struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"index;not null"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Key         string         `json:"key" gorm:"size:100;uniqueIndex;not null"`
	Secret      string         `json:"-" gorm:"size:255;not null"` // 不暴露给客户端
	Permissions []Permission   `json:"permissions" gorm:"many2many:api_key_permissions"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// LoginLog 表示登录日志
type LoginLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	IP        string    `json:"ip" gorm:"size:50;not null"`
	UserAgent string    `json:"user_agent" gorm:"size:255"`
	Status    string    `json:"status" gorm:"size:20;not null"` // success, failed
	Message   *string   `json:"message" gorm:"size:255"`        // 登录失败原因
	Location  *string   `json:"location" gorm:"size:100"`       // 登录地点
	CreatedAt time.Time `json:"created_at"`
}

// TwoFactorAuth 表示二次验证配置
type TwoFactorAuth struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"uniqueIndex;not null"`
	Secret        string    `json:"-" gorm:"size:50;not null"` // 不暴露给客户端
	IsEnabled     bool      `json:"is_enabled" gorm:"default:false"`
	RecoveryCodes string    `json:"-" gorm:"type:text"` // JSON格式的恢复码数组
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
