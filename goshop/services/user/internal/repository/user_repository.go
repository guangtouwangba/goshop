package repository

import (
	"context"

	"github.com/yourusername/goshop/services/user/internal/model"
	"gorm.io/gorm"
)

// UserRepository 定义用户仓库接口
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*model.User, int64, error)
	VerifyEmail(ctx context.Context, id uint) error
	VerifyPhone(ctx context.Context, id uint) error
	UpdateLastLogin(ctx context.Context, id uint) error
	AddPoints(ctx context.Context, id uint, points int) error
	UpdateMemberLevel(ctx context.Context, id uint, level int) error
	AddLoginHistory(ctx context.Context, history *model.LoginHistory) error
	GetLoginHistory(ctx context.Context, userID uint, limit int) ([]*model.LoginHistory, error)
}

// GormUserRepository 实现 UserRepository 接口的 GORM 仓库
type GormUserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{
		db: db,
	}
}

// Create 创建新用户
func (r *GormUserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据 ID 获取用户
func (r *GormUserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Addresses").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *GormUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *GormUserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户（软删除）
func (r *GormUserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// List 获取用户列表
func (r *GormUserRepository) List(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// VerifyEmail 验证用户邮箱
func (r *GormUserRepository) VerifyEmail(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("email_verified", true).Error
}

// VerifyPhone 验证用户手机
func (r *GormUserRepository) VerifyPhone(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("phone_verified", true).Error
}

// UpdateLastLogin 更新最后登录时间
func (r *GormUserRepository) UpdateLastLogin(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("last_login_at", gorm.Expr("NOW()")).Error
}

// AddPoints 添加积分
func (r *GormUserRepository) AddPoints(ctx context.Context, id uint, points int) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).
		Update("points", gorm.Expr("points + ?", points)).Error
}

// UpdateMemberLevel 更新会员等级
func (r *GormUserRepository) UpdateMemberLevel(ctx context.Context, id uint, level int) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).
		Update("member_level", level).Error
}

// AddLoginHistory 添加登录历史
func (r *GormUserRepository) AddLoginHistory(ctx context.Context, history *model.LoginHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetLoginHistory 获取用户登录历史
func (r *GormUserRepository) GetLoginHistory(ctx context.Context, userID uint, limit int) ([]*model.LoginHistory, error) {
	var histories []*model.LoginHistory

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories).Error

	if err != nil {
		return nil, err
	}

	return histories, nil
}
