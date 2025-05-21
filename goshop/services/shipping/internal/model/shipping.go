package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ShippingMethod 表示物流配送方式
type ShippingMethod struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"size:50;not null"`
	Code          string         `json:"code" gorm:"size:20;uniqueIndex;not null"`
	Description   string         `json:"description" gorm:"size:255"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	SortOrder     int            `json:"sort_order" gorm:"default:0"`
	EstimatedDays string         `json:"estimated_days" gorm:"size:50"` // 预计送达时间，如"3-5天"
	Icon          *string        `json:"icon" gorm:"size:255"`
	CarrierIDs    UintSlice      `json:"carrier_ids" gorm:"type:jsonb"` // 关联的物流公司ID
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
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

// ShippingCarrier 表示物流承运商/快递公司
type ShippingCarrier struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Code        string         `json:"code" gorm:"size:20;uniqueIndex;not null"`
	TrackingURL string         `json:"tracking_url" gorm:"size:255"` // 物流追踪URL模板，例如"https://example.com/track/{tracking_number}"
	Logo        *string        `json:"logo" gorm:"size:255"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	APICode     *string        `json:"api_code" gorm:"size:50"` // 第三方物流API的代码
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

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

// ShippingRateConditionType 表示运费计算条件类型
type ShippingRateConditionType string

const (
	// ShippingRateConditionTypeWeight 按重量计算
	ShippingRateConditionTypeWeight ShippingRateConditionType = "weight"
	// ShippingRateConditionTypePrice 按价格计算
	ShippingRateConditionTypePrice ShippingRateConditionType = "price"
	// ShippingRateConditionTypeQuantity 按数量计算
	ShippingRateConditionTypeQuantity ShippingRateConditionType = "quantity"
)

// ShippingZone 表示运费区域
type ShippingZone struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Description string         `json:"description" gorm:"size:255"`
	RegionCodes []string       `json:"region_codes" gorm:"type:text[]"` // 地区代码列表，如省份/城市代码
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// ShippingRate 表示运费计算规则
type ShippingRate struct {
	ID               uint                      `json:"id" gorm:"primaryKey"`
	ShippingMethodID uint                      `json:"shipping_method_id" gorm:"index;not null"`
	ShippingZoneID   uint                      `json:"shipping_zone_id" gorm:"index;not null"`
	Name             string                    `json:"name" gorm:"size:50;not null"`
	ConditionType    ShippingRateConditionType `json:"condition_type" gorm:"size:20;not null"`
	ConditionMin     float64                   `json:"condition_min" gorm:"type:decimal(10,2);not null"`    // 条件最小值
	ConditionMax     *float64                  `json:"condition_max" gorm:"type:decimal(10,2)"`             // 条件最大值，null表示无上限
	BaseRate         float64                   `json:"base_rate" gorm:"type:decimal(10,2);not null"`        // 基础运费
	AdditionalRate   float64                   `json:"additional_rate" gorm:"type:decimal(10,2);default:0"` // 附加费率，每超出一个单位的费用
	AdditionalUnit   float64                   `json:"additional_unit" gorm:"type:decimal(10,2);default:1"` // 附加单位，如每超出1公斤
	IsFreeThreshold  bool                      `json:"is_free_threshold" gorm:"default:false"`              // 是否有包邮门槛
	FreeThreshold    *float64                  `json:"free_threshold" gorm:"type:decimal(10,2)"`            // 包邮条件值，如订单金额超过此值免运费
	IsActive         bool                      `json:"is_active" gorm:"default:true"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at"`
	DeletedAt        gorm.DeletedAt            `json:"-" gorm:"index"`
}

// Shipment 表示物流配送信息
type Shipment struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	OrderID             uint           `json:"order_id" gorm:"index;not null"`
	OrderNumber         string         `json:"order_number" gorm:"size:50;not null"`
	UserID              uint           `json:"user_id" gorm:"index"`
	ShippingMethodID    uint           `json:"shipping_method_id" gorm:"index"`
	ShippingMethodName  string         `json:"shipping_method_name" gorm:"size:50"`
	ShippingCarrierID   *uint          `json:"shipping_carrier_id" gorm:"index"`
	ShippingCarrierName *string        `json:"shipping_carrier_name" gorm:"size:50"`
	TrackingNumber      *string        `json:"tracking_number" gorm:"size:100"`
	TrackingURL         *string        `json:"tracking_url" gorm:"size:255"`
	ShippedAt           *time.Time     `json:"shipped_at"`
	DeliveredAt         *time.Time     `json:"delivered_at"`
	Status              string         `json:"status" gorm:"size:20;default:'pending'"`         // pending, shipped, delivered, failed
	Address             JSONMap        `json:"address" gorm:"type:jsonb;not null"`              // 配送地址
	Items               JSONMap        `json:"items" gorm:"type:jsonb;not null"`                // 配送商品信息
	TrackingInfo        JSONMap        `json:"tracking_info" gorm:"type:jsonb"`                 // 物流追踪信息
	ShippingFee         float64        `json:"shipping_fee" gorm:"type:decimal(10,2);not null"` // 运费
	Note                *string        `json:"note" gorm:"size:255"`                            // 配送备注
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
}
