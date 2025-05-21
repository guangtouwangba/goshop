package model

import (
	"time"

	"gorm.io/gorm"
)

// OrderStatus 表示订单状态
type OrderStatus string

const (
	// OrderStatusPending 待付款
	OrderStatusPending OrderStatus = "pending"
	// OrderStatusPaid 已付款
	OrderStatusPaid OrderStatus = "paid"
	// OrderStatusProcessing 处理中
	OrderStatusProcessing OrderStatus = "processing"
	// OrderStatusShipped 已发货
	OrderStatusShipped OrderStatus = "shipped"
	// OrderStatusDelivered 已送达
	OrderStatusDelivered OrderStatus = "delivered"
	// OrderStatusCompleted 已完成
	OrderStatusCompleted OrderStatus = "completed"
	// OrderStatusCancelled 已取消
	OrderStatusCancelled OrderStatus = "cancelled"
	// OrderStatusRefunded 已退款
	OrderStatusRefunded OrderStatus = "refunded"
	// OrderStatusPartiallyRefunded 部分退款
	OrderStatusPartiallyRefunded OrderStatus = "partially_refunded"
	// OrderStatusFailed 失败
	OrderStatusFailed OrderStatus = "failed"
)

// PaymentStatus 表示支付状态
type PaymentStatus string

const (
	// PaymentStatusPending 待支付
	PaymentStatusPending PaymentStatus = "pending"
	// PaymentStatusPaid 已支付
	PaymentStatusPaid PaymentStatus = "paid"
	// PaymentStatusFailed 支付失败
	PaymentStatusFailed PaymentStatus = "failed"
	// PaymentStatusRefunded 已退款
	PaymentStatusRefunded PaymentStatus = "refunded"
	// PaymentStatusPartiallyRefunded 部分退款
	PaymentStatusPartiallyRefunded PaymentStatus = "partially_refunded"
)

// Order 表示订单
type Order struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	OrderNumber     string         `json:"order_number" gorm:"uniqueIndex;size:50;not null"` // 订单号
	UserID          uint           `json:"user_id" gorm:"index"`                             // 用户ID
	Status          OrderStatus    `json:"status" gorm:"size:30;not null;default:'pending'"`
	PaymentStatus   PaymentStatus  `json:"payment_status" gorm:"size:30;not null;default:'pending'"`
	PaymentMethod   string         `json:"payment_method" gorm:"size:50"`                             // 支付方式
	TransactionID   *string        `json:"transaction_id" gorm:"size:100"`                            // 支付交易号
	ShippingMethod  string         `json:"shipping_method" gorm:"size:50"`                            // 配送方式
	ShippingCarrier *string        `json:"shipping_carrier" gorm:"size:50"`                           // 配送公司
	TrackingNumber  *string        `json:"tracking_number" gorm:"size:100"`                           // 物流单号
	Items           []OrderItem    `json:"items" gorm:"foreignKey:OrderID"`                           // 订单项
	CouponCode      *string        `json:"coupon_code" gorm:"size:50"`                                // 优惠券码
	ShippingAddress Address        `json:"shipping_address" gorm:"embedded;embeddedPrefix:shipping_"` // 收货地址
	BillingAddress  Address        `json:"billing_address" gorm:"embedded;embeddedPrefix:billing_"`   // 账单地址
	Subtotal        float64        `json:"subtotal" gorm:"type:decimal(10,2);not null"`               // 小计（未含税、运费）
	ShippingFee     float64        `json:"shipping_fee" gorm:"type:decimal(10,2);not null"`           // 运费
	Tax             float64        `json:"tax" gorm:"type:decimal(10,2);not null"`                    // 税费
	Discount        float64        `json:"discount" gorm:"type:decimal(10,2);not null"`               // 优惠金额
	GrandTotal      float64        `json:"grand_total" gorm:"type:decimal(10,2);not null"`            // 总计
	Note            *string        `json:"note" gorm:"type:text"`                                     // 订单备注
	CustomerNote    *string        `json:"customer_note" gorm:"type:text"`                            // 客户备注
	InternalNote    *string        `json:"internal_note" gorm:"type:text"`                            // 内部备注
	PaidAt          *time.Time     `json:"paid_at"`                                                   // 支付时间
	ShippedAt       *time.Time     `json:"shipped_at"`                                                // 发货时间
	DeliveredAt     *time.Time     `json:"delivered_at"`                                              // 送达时间
	CompletedAt     *time.Time     `json:"completed_at"`                                              // 完成时间
	CancelledAt     *time.Time     `json:"cancelled_at"`                                              // 取消时间
	RefundedAt      *time.Time     `json:"refunded_at"`                                               // 退款时间
	ExpiredAt       *time.Time     `json:"expired_at"`                                                // 过期时间（未支付自动取消）
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// OrderItem 表示订单项
type OrderItem struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	OrderID       uint      `json:"order_id" gorm:"index;not null"`
	ProductID     uint      `json:"product_id" gorm:"index;not null"`
	SKUID         uint      `json:"sku_id" gorm:"index;not null"`
	ProductName   string    `json:"product_name" gorm:"size:255;not null"`
	SKUCode       string    `json:"sku_code" gorm:"size:50;not null"`
	VariantName   string    `json:"variant_name" gorm:"size:255"`
	Price         float64   `json:"price" gorm:"type:decimal(10,2);not null"`    // 单价
	OriginalPrice float64   `json:"original_price" gorm:"type:decimal(10,2)"`    // 原价
	Quantity      int       `json:"quantity" gorm:"not null"`                    // 数量
	Subtotal      float64   `json:"subtotal" gorm:"type:decimal(10,2);not null"` // 小计
	Tax           float64   `json:"tax" gorm:"type:decimal(10,2);not null"`      // 税费
	Discount      float64   `json:"discount" gorm:"type:decimal(10,2);not null"` // 折扣
	Total         float64   `json:"total" gorm:"type:decimal(10,2);not null"`    // 总计
	Weight        *float64  `json:"weight" gorm:"type:decimal(10,2)"`            // 重量
	Image         *string   `json:"image" gorm:"size:255"`                       // 图片
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Address 表示地址
type Address struct {
	Name         string `json:"name" gorm:"size:50"`           // 收货人姓名
	Phone        string `json:"phone" gorm:"size:20"`          // 联系电话
	Province     string `json:"province" gorm:"size:50"`       // 省
	City         string `json:"city" gorm:"size:50"`           // 市
	District     string `json:"district" gorm:"size:50"`       // 区
	DetailedInfo string `json:"detailed_info" gorm:"size:255"` // 详细地址
	PostalCode   string `json:"postal_code" gorm:"size:20"`    // 邮编
}

// OrderLog 表示订单操作日志
type OrderLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	OrderID     uint      `json:"order_id" gorm:"index;not null"`
	UserID      *uint     `json:"user_id" gorm:"index"`           // 操作人ID，可能是系统操作
	Action      string    `json:"action" gorm:"size:30;not null"` // 操作类型：如 status_change, payment, note_update
	StatusFrom  *string   `json:"status_from" gorm:"size:30"`     // 状态变更前
	StatusTo    *string   `json:"status_to" gorm:"size:30"`       // 状态变更后
	Description string    `json:"description" gorm:"size:255"`    // 描述
	IP          *string   `json:"ip" gorm:"size:50"`              // 操作IP
	UserAgent   *string   `json:"user_agent" gorm:"size:255"`     // 用户代理
	CreatedAt   time.Time `json:"created_at"`                     // 操作时间
}

// Cart 表示购物车
type Cart struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    *uint      `json:"user_id" gorm:"index"`             // 用户ID，游客可以为空
	SessionID string     `json:"session_id" gorm:"size:100;index"` // 会话ID，用于游客
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CartItem 表示购物车项
type CartItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CartID    uint      `json:"cart_id" gorm:"index;not null"`
	ProductID uint      `json:"product_id" gorm:"index;not null"`
	SKUID     uint      `json:"sku_id" gorm:"index;not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
