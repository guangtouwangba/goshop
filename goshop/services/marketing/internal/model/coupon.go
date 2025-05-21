package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CouponType 表示优惠券类型
type CouponType string

const (
	// CouponTypeFixedAmount 满减券
	CouponTypeFixedAmount CouponType = "fixed_amount"
	// CouponTypePercentage 折扣券
	CouponTypePercentage CouponType = "percentage"
	// CouponTypeFreeShipping 包邮券
	CouponTypeFreeShipping CouponType = "free_shipping"
	// CouponTypeProductSpecific 指定商品券
	CouponTypeProductSpecific CouponType = "product_specific"
	// CouponTypeCategorySpecific 指定分类券
	CouponTypeCategorySpecific CouponType = "category_specific"
	// CouponTypeFirstOrder 首单券
	CouponTypeFirstOrder CouponType = "first_order"
)

// StringSlice 是一个自定义类型，用于存储字符串数组
type StringSlice []string

// Value 实现 driver.Valuer 接口
func (a StringSlice) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口
func (a *StringSlice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为 []byte 失败")
	}
	return json.Unmarshal(b, &a)
}

// UintSlice 是一个自定义类型，用于存储uint数组
type UintSlice []uint

// Value 实现 driver.Valuer 接口
func (a UintSlice) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口
func (a *UintSlice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为 []byte 失败")
	}
	return json.Unmarshal(b, &a)
}

// Coupon 表示优惠券
type Coupon struct {
	ID                   uint           `json:"id" gorm:"primaryKey"`
	Code                 string         `json:"code" gorm:"size:50;uniqueIndex;not null"`             // 优惠码
	Name                 string         `json:"name" gorm:"size:100;not null"`                        // 优惠券名称
	Description          string         `json:"description" gorm:"size:255"`                          // 优惠券描述
	Type                 CouponType     `json:"type" gorm:"size:20;not null"`                         // 优惠券类型
	Value                float64        `json:"value" gorm:"type:decimal(10,2);not null"`             // 优惠金额或折扣百分比
	MinOrderAmount       float64        `json:"min_order_amount" gorm:"type:decimal(10,2);default:0"` // 最低订单金额
	MaxDiscountAmount    *float64       `json:"max_discount_amount" gorm:"type:decimal(10,2)"`        // 最大折扣金额（对于百分比折扣）
	StartAt              time.Time      `json:"start_at" gorm:"not null"`                             // 生效时间
	EndAt                time.Time      `json:"end_at" gorm:"not null"`                               // 失效时间
	TotalQuantity        int            `json:"total_quantity" gorm:"default:0"`                      // 发行量，0表示不限量
	UsedQuantity         int            `json:"used_quantity" gorm:"default:0"`                       // 已使用数量
	UserLimit            int            `json:"user_limit" gorm:"default:1"`                          // 每个用户可使用次数，0表示不限制
	IsActive             bool           `json:"is_active" gorm:"default:true"`                        // 是否激活
	ApplicableProducts   UintSlice      `json:"applicable_products" gorm:"type:jsonb"`                // 适用商品ID
	ApplicableCategories UintSlice      `json:"applicable_categories" gorm:"type:jsonb"`              // 适用分类ID
	ExcludedProducts     UintSlice      `json:"excluded_products" gorm:"type:jsonb"`                  // 排除商品ID
	ExcludedCategories   UintSlice      `json:"excluded_categories" gorm:"type:jsonb"`                // 排除分类ID
	IsForNewUser         bool           `json:"is_for_new_user" gorm:"default:false"`                 // 是否仅限新用户使用
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`
}

// CouponUsage 表示优惠券使用记录
type CouponUsage struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	CouponID       uint      `json:"coupon_id" gorm:"index;not null"`
	UserID         uint      `json:"user_id" gorm:"index;not null"`
	OrderID        uint      `json:"order_id" gorm:"index;not null"`
	OrderNumber    string    `json:"order_number" gorm:"size:50;not null"`
	UsedAt         time.Time `json:"used_at"`
	DiscountAmount float64   `json:"discount_amount" gorm:"type:decimal(10,2);not null"` // 优惠金额
	CreatedAt      time.Time `json:"created_at"`
}
