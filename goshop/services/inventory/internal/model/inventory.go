package model

import (
	"time"

	"gorm.io/gorm"
)

// StockOperation 表示库存操作类型
type StockOperation string

const (
	// StockOperationInitial 初始库存
	StockOperationInitial StockOperation = "initial"
	// StockOperationIncrease 增加库存
	StockOperationIncrease StockOperation = "increase"
	// StockOperationDecrease 减少库存
	StockOperationDecrease StockOperation = "decrease"
	// StockOperationHold 锁定库存（订单创建时）
	StockOperationHold StockOperation = "hold"
	// StockOperationRelease 释放库存（订单取消时）
	StockOperationRelease StockOperation = "release"
	// StockOperationConfirm 确认消耗库存（订单支付时）
	StockOperationConfirm StockOperation = "confirm"
	// StockOperationAdjust 库存调整
	StockOperationAdjust StockOperation = "adjust"
)

// InventoryActionSource 表示库存操作来源
type InventoryActionSource string

const (
	// InventoryActionSourceOrder 来自订单
	InventoryActionSourceOrder InventoryActionSource = "order"
	// InventoryActionSourceManual 手动操作
	InventoryActionSourceManual InventoryActionSource = "manual"
	// InventoryActionSourceSystem 系统操作
	InventoryActionSourceSystem InventoryActionSource = "system"
	// InventoryActionSourceAPI API接口操作
	InventoryActionSourceAPI InventoryActionSource = "api"
)

// StockStrategy 表示库存扣减策略
type StockStrategy string

const (
	// StockStrategyOrder 下单减库存
	StockStrategyOrder StockStrategy = "order"
	// StockStrategyPayment 付款减库存
	StockStrategyPayment StockStrategy = "payment"
)

// SKUStock 表示SKU库存
type SKUStock struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	SKUID           uint           `json:"sku_id" gorm:"uniqueIndex;not null"`                       // SKU ID
	AvailableStock  int            `json:"available_stock" gorm:"not null"`                          // 可用库存
	HoldStock       int            `json:"hold_stock" gorm:"not null"`                               // 锁定库存（未付款）
	StockStrategy   StockStrategy  `json:"stock_strategy" gorm:"size:20;not null;default:'payment'"` // 库存扣减策略
	LowStockAlert   int            `json:"low_stock_alert" gorm:"default:10"`                        // 低库存预警值
	IsInfinite      bool           `json:"is_infinite" gorm:"default:false"`                         // 是否不限库存
	WarehouseID     *uint          `json:"warehouse_id" gorm:"index"`                                // 仓库ID，可选
	LastStockUpdate *time.Time     `json:"last_stock_update"`                                        // 最后库存更新时间
	StockStatus     string         `json:"stock_status" gorm:"size:20;default:'in_stock'"`           // 库存状态：in_stock, out_of_stock, low_stock
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// StockMovement 表示库存流水
type StockMovement struct {
	ID            uint                  `json:"id" gorm:"primaryKey"`
	SKUID         uint                  `json:"sku_id" gorm:"index;not null"`
	Quantity      int                   `json:"quantity" gorm:"not null"` // 正值为增加，负值为减少
	Operation     StockOperation        `json:"operation" gorm:"size:20;not null"`
	BeforeStock   int                   `json:"before_stock"` // 操作前库存
	AfterStock    int                   `json:"after_stock"`  // 操作后库存
	Source        InventoryActionSource `json:"source" gorm:"size:20;not null;default:'manual'"`
	ReferenceID   *string               `json:"reference_id" gorm:"size:50"`   // 关联ID（如订单ID）
	ReferenceType *string               `json:"reference_type" gorm:"size:20"` // 关联类型（如order）
	Note          *string               `json:"note" gorm:"size:255"`
	OperatorID    *uint                 `json:"operator_id" gorm:"index"` // 操作人ID
	WarehouseID   *uint                 `json:"warehouse_id" gorm:"index"`
	CreatedAt     time.Time             `json:"created_at"`
}

// Warehouse 表示仓库
type Warehouse struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name" gorm:"size:50;not null"`
	Code       string         `json:"code" gorm:"size:20;uniqueIndex;not null"`
	Address    string         `json:"address" gorm:"size:255"`
	Province   string         `json:"province" gorm:"size:50"`
	City       string         `json:"city" gorm:"size:50"`
	District   string         `json:"district" gorm:"size:50"`
	PostalCode string         `json:"postal_code" gorm:"size:20"`
	Contact    string         `json:"contact" gorm:"size:50"`
	Phone      string         `json:"phone" gorm:"size:20"`
	Email      *string        `json:"email" gorm:"size:100"`
	IsDefault  bool           `json:"is_default" gorm:"default:false"`
	IsActive   bool           `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// StockAlert 表示库存预警记录
type StockAlert struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SKUID      uint      `json:"sku_id" gorm:"index;not null"`
	StockLevel int       `json:"stock_level" gorm:"not null"`             // 当前库存
	AlertLevel int       `json:"alert_level" gorm:"not null"`             // 预警阈值
	Status     string    `json:"status" gorm:"size:20;default:'pending'"` // pending, processed, dismissed
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
