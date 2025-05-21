package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// PaymentMethod 支付方式
type PaymentMethod string

const (
	// PaymentMethodStripe Stripe 支付
	PaymentMethodStripe PaymentMethod = "stripe"
	// PaymentMethodPayPal PayPal 支付
	PaymentMethodPayPal PaymentMethod = "paypal"
	// PaymentMethodWechat 微信支付
	PaymentMethodWechat PaymentMethod = "wechat"
	// PaymentMethodAlipay 支付宝
	PaymentMethodAlipay PaymentMethod = "alipay"
	// PaymentMethodBankTransfer 银行转账
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	// PaymentMethodCreditCard 信用卡
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	// PaymentMethodCOD 货到付款
	PaymentMethodCOD PaymentMethod = "cod"
)

// PaymentStatus 支付状态
type PaymentStatus string

const (
	// PaymentStatusPending 待支付
	PaymentStatusPending PaymentStatus = "pending"
	// PaymentStatusProcessing 处理中
	PaymentStatusProcessing PaymentStatus = "processing"
	// PaymentStatusSuccess 成功
	PaymentStatusSuccess PaymentStatus = "success"
	// PaymentStatusFailed 失败
	PaymentStatusFailed PaymentStatus = "failed"
	// PaymentStatusRefunding 退款中
	PaymentStatusRefunding PaymentStatus = "refunding"
	// PaymentStatusRefunded 已退款
	PaymentStatusRefunded PaymentStatus = "refunded"
	// PaymentStatusPartialRefunded 部分退款
	PaymentStatusPartialRefunded PaymentStatus = "partially_refunded"
	// PaymentStatusCancelled 已取消
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// JSONMap 是一个自定义类型，用于存储 JSON 对象
type JSONMap map[string]interface{}

// Value 实现 driver.Valuer 接口
func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口
func (j *JSONMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为 []byte 失败")
	}
	return json.Unmarshal(b, &j)
}

// Payment 支付记录
type Payment struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	OrderID           uint           `json:"order_id" gorm:"index;not null"`
	OrderNumber       string         `json:"order_number" gorm:"size:50;index;not null"`
	UserID            uint           `json:"user_id" gorm:"index"`
	PaymentMethod     PaymentMethod  `json:"payment_method" gorm:"size:20;not null"`
	Amount            float64        `json:"amount" gorm:"type:decimal(10,2);not null"`
	Currency          string         `json:"currency" gorm:"size:3;not null;default:'CNY'"`
	Status            PaymentStatus  `json:"status" gorm:"size:20;not null;default:'pending'"`
	TransactionID     *string        `json:"transaction_id" gorm:"size:100;index"` // 支付平台的交易ID
	PaymentGatewayRef *string        `json:"payment_gateway_ref" gorm:"size:100"`  // 支付网关的引用ID
	PaymentData       JSONMap        `json:"payment_data" gorm:"type:jsonb"`       // 支付相关的其他数据
	ErrorMessage      *string        `json:"error_message" gorm:"type:text"`       // 错误信息
	ClientIP          string         `json:"client_ip" gorm:"size:50"`             // 客户端IP
	ReturnURL         string         `json:"return_url" gorm:"size:255"`           // 支付成功后的回调URL
	NotifyURL         string         `json:"notify_url" gorm:"size:255"`           // 支付网关异步通知URL
	ExpiredAt         *time.Time     `json:"expired_at"`                           // 支付过期时间
	PaidAt            *time.Time     `json:"paid_at"`                              // 支付成功时间
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// Refund 退款记录
type Refund struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	PaymentID     uint           `json:"payment_id" gorm:"index;not null"`
	OrderID       uint           `json:"order_id" gorm:"index;not null"`
	UserID        uint           `json:"user_id" gorm:"index"`
	Amount        float64        `json:"amount" gorm:"type:decimal(10,2);not null"`
	Currency      string         `json:"currency" gorm:"size:3;not null;default:'CNY'"`
	Reason        string         `json:"reason" gorm:"size:255"`
	Status        PaymentStatus  `json:"status" gorm:"size:20;not null;default:'processing'"`
	TransactionID *string        `json:"transaction_id" gorm:"size:100;index"` // 退款交易ID
	RefundData    JSONMap        `json:"refund_data" gorm:"type:jsonb"`        // 退款相关的其他数据
	ErrorMessage  *string        `json:"error_message" gorm:"type:text"`       // 错误信息
	OperatorID    *uint          `json:"operator_id" gorm:"index"`             // 操作人ID
	RefundedAt    *time.Time     `json:"refunded_at"`                          // 退款成功时间
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// PaymentGateway 支付网关配置
type PaymentGateway struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	Name                string         `json:"name" gorm:"size:50;not null"`
	Code                PaymentMethod  `json:"code" gorm:"size:20;uniqueIndex;not null"`
	Description         string         `json:"description" gorm:"size:255"`
	Logo                *string        `json:"logo" gorm:"size:255"`
	IsActive            bool           `json:"is_active" gorm:"default:true"`
	IsSandbox           bool           `json:"is_sandbox" gorm:"default:true"`          // 是否沙盒模式
	Config              JSONMap        `json:"config" gorm:"type:jsonb;not null"`       // 配置信息，如 API 密钥等
	SupportedCurrencies []string       `json:"supported_currencies" gorm:"type:text[]"` // 支持的货币
	SortOrder           int            `json:"sort_order" gorm:"default:0"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
}

// PaymentLog 支付操作日志
type PaymentLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PaymentID  uint      `json:"payment_id" gorm:"index"`
	RefundID   *uint     `json:"refund_id" gorm:"index"`
	Action     string    `json:"action" gorm:"size:30;not null"` // 操作类型：如 create, update, notify, refund
	StatusFrom *string   `json:"status_from" gorm:"size:30"`     // 状态变更前
	StatusTo   *string   `json:"status_to" gorm:"size:30"`       // 状态变更后
	Data       JSONMap   `json:"data" gorm:"type:jsonb"`         // 操作相关数据
	IP         *string   `json:"ip" gorm:"size:50"`              // 操作IP
	UserAgent  *string   `json:"user_agent" gorm:"size:255"`     // 用户代理
	OperatorID *uint     `json:"operator_id" gorm:"index"`       // 操作人ID
	CreatedAt  time.Time `json:"created_at"`                     // 操作时间
}
