package model

import (
	"time"

	"gorm.io/gorm"
)

// PromotionType 表示促销活动类型
type PromotionType string

const (
	// PromotionTypeFlashSale 限时特价
	PromotionTypeFlashSale PromotionType = "flash_sale"
	// PromotionTypeBundleSale 捆绑销售
	PromotionTypeBundleSale PromotionType = "bundle_sale"
	// PromotionTypeBuyXGetY 买X送Y
	PromotionTypeBuyXGetY PromotionType = "buy_x_get_y"
	// PromotionTypeSecondHalfPrice 第二件半价
	PromotionTypeSecondHalfPrice PromotionType = "second_half_price"
	// PromotionTypeSpendGetFree 满赠活动
	PromotionTypeSpendGetFree PromotionType = "spend_get_free"
	// PromotionTypeQuantityDiscount 阶梯式优惠
	PromotionTypeQuantityDiscount PromotionType = "quantity_discount"
)

// Promotion 表示促销活动
type Promotion struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"size:100;not null"`
	Description    string         `json:"description" gorm:"size:500"`
	Type           PromotionType  `json:"type" gorm:"size:30;not null"`
	StartAt        time.Time      `json:"start_at" gorm:"not null"`
	EndAt          time.Time      `json:"end_at" gorm:"not null"`
	IsActive       bool           `json:"is_active" gorm:"default:true"`
	Priority       int            `json:"priority" gorm:"default:0"`                  // 优先级，越高越优先
	ProductIDs     UintSlice      `json:"product_ids" gorm:"type:jsonb"`              // 适用商品ID
	CategoryIDs    UintSlice      `json:"category_ids" gorm:"type:jsonb"`             // 适用分类ID
	DiscountValue  float64        `json:"discount_value" gorm:"type:decimal(10,2)"`   // 折扣值（金额或百分比）
	DiscountType   string         `json:"discount_type" gorm:"size:20"`               // amount或percentage
	MinOrderAmount *float64       `json:"min_order_amount" gorm:"type:decimal(10,2)"` // 最低订单金额
	MinQuantity    *int           `json:"min_quantity"`                               // 最低购买数量
	MaxUsesPerUser *int           `json:"max_uses_per_user"`                          // 每个用户最大使用次数
	TotalUses      int            `json:"total_uses" gorm:"default:0"`                // 总使用次数
	MaxUses        *int           `json:"max_uses"`                                   // 最大使用次数，null表示不限
	FreeProductID  *uint          `json:"free_product_id"`                            // 赠品ID
	FreeProductQty *int           `json:"free_product_qty"`                           // 赠品数量
	Rules          StringSlice    `json:"rules" gorm:"type:jsonb"`                    // 促销规则，例如阶梯式优惠规则
	Image          *string        `json:"image" gorm:"size:255"`                      // 活动图片
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// PromotionUsage 表示促销活动使用记录
type PromotionUsage struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	PromotionID    uint      `json:"promotion_id" gorm:"index;not null"`
	UserID         uint      `json:"user_id" gorm:"index;not null"`
	OrderID        uint      `json:"order_id" gorm:"index;not null"`
	OrderNumber    string    `json:"order_number" gorm:"size:50;not null"`
	DiscountAmount float64   `json:"discount_amount" gorm:"type:decimal(10,2);not null"` // 优惠金额
	UsedAt         time.Time `json:"used_at"`
	CreatedAt      time.Time `json:"created_at"`
}

// LoyaltyPointRule 表示积分规则
type LoyaltyPointRule struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"size:100;not null"`
	Description    string         `json:"description" gorm:"size:500"`
	PointsPerSpend int            `json:"points_per_spend" gorm:"default:1"`                    // 每消费1元获得的积分
	MinOrderAmount float64        `json:"min_order_amount" gorm:"type:decimal(10,2);default:0"` // 最低订单金额
	IsActive       bool           `json:"is_active" gorm:"default:true"`
	StartAt        *time.Time     `json:"start_at"` // 生效时间，null表示永久有效
	EndAt          *time.Time     `json:"end_at"`   // 失效时间，null表示永久有效
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// LoyaltyPointTransaction 表示积分交易
type LoyaltyPointTransaction struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	UserID        uint       `json:"user_id" gorm:"index;not null"`
	Points        int        `json:"points" gorm:"not null"`        // 正值为获得，负值为使用
	Balance       int        `json:"balance" gorm:"not null"`       // 交易后的积分余额
	Type          string     `json:"type" gorm:"size:20;not null"`  // earn, redeem, expire, adjust
	ReferenceID   *string    `json:"reference_id" gorm:"size:50"`   // 关联ID（如订单ID）
	ReferenceType *string    `json:"reference_type" gorm:"size:20"` // 关联类型（如order）
	Description   string     `json:"description" gorm:"size:255"`
	ExpiresAt     *time.Time `json:"expires_at"` // 过期时间，null表示永不过期
	CreatedAt     time.Time  `json:"created_at"`
}

// MemberLevel 表示会员等级
type MemberLevel struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"size:50;not null"`
	Level          int            `json:"level" gorm:"not null;uniqueIndex"`                   // 等级值，数字越大等级越高
	RequiredPoints int            `json:"required_points" gorm:"not null"`                     // 所需积分
	DiscountRate   float64        `json:"discount_rate" gorm:"type:decimal(5,2);default:1.00"` // 折扣率，例如0.95表示95折
	Description    string         `json:"description" gorm:"size:500"`
	Icon           *string        `json:"icon" gorm:"size:255"`
	IsActive       bool           `json:"is_active" gorm:"default:true"`
	Benefits       StringSlice    `json:"benefits" gorm:"type:jsonb"` // 会员权益列表
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}
