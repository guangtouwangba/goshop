package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ProductType 定义商品类型
type ProductType string

const (
	// ProductTypePhysical 实体商品
	ProductTypePhysical ProductType = "physical"
	// ProductTypeVirtual 虚拟商品（如电子书、音乐）
	ProductTypeVirtual ProductType = "virtual"
	// ProductTypeSubscription 订阅商品
	ProductTypeSubscription ProductType = "subscription"
	// ProductTypeBundle 商品套装
	ProductTypeBundle ProductType = "bundle"
	// ProductTypeGiftBox 主题礼盒/福袋
	ProductTypeGiftBox ProductType = "gift_box"
)

// ProductStatus 定义商品状态
type ProductStatus string

const (
	// ProductStatusDraft 草稿状态
	ProductStatusDraft ProductStatus = "draft"
	// ProductStatusActive 已上架
	ProductStatusActive ProductStatus = "active"
	// ProductStatusInactive 已下架
	ProductStatusInactive ProductStatus = "inactive"
	// ProductStatusScheduled 定时上下架
	ProductStatusScheduled ProductStatus = "scheduled"
)

// StringArray 是一个自定义类型，用于存储字符串数组
type StringArray []string

// Value 实现 driver.Valuer 接口
func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口
func (a *StringArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为 []byte 失败")
	}
	return json.Unmarshal(b, &a)
}

// Product 定义商品模型
type Product struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	Name              string         `json:"name" gorm:"size:255;not null"`
	Description       string         `json:"description" gorm:"type:text"`
	ShortDescription  string         `json:"short_description" gorm:"size:500"`
	Type              ProductType    `json:"type" gorm:"size:20;not null;default:'physical'"`
	Status            ProductStatus  `json:"status" gorm:"size:20;not null;default:'draft'"`
	SKUs              []SKU          `json:"skus" gorm:"foreignKey:ProductID"`
	RegularPrice      float64        `json:"regular_price" gorm:"type:decimal(10,2);not null"`
	SalePrice         *float64       `json:"sale_price" gorm:"type:decimal(10,2)"`
	SaleStartDate     *time.Time     `json:"sale_start_date"`
	SaleEndDate       *time.Time     `json:"sale_end_date"`
	InventoryTracking bool           `json:"inventory_tracking" gorm:"default:true"`
	PublishDate       *time.Time     `json:"publish_date"`
	UnpublishDate     *time.Time     `json:"unpublish_date"`
	Weight            *float64       `json:"weight" gorm:"type:decimal(10,2)"` // 重量（公斤）
	Length            *float64       `json:"length" gorm:"type:decimal(10,2)"` // 长度（厘米）
	Width             *float64       `json:"width" gorm:"type:decimal(10,2)"`  // 宽度（厘米）
	Height            *float64       `json:"height" gorm:"type:decimal(10,2)"` // 高度（厘米）
	Images            StringArray    `json:"images" gorm:"type:jsonb"`
	Videos            StringArray    `json:"videos" gorm:"type:jsonb"`
	Categories        []Category     `json:"categories" gorm:"many2many:product_categories"`
	Brand             *Brand         `json:"brand" gorm:"foreignKey:BrandID"`
	BrandID           *uint          `json:"brand_id"`
	Tags              StringArray    `json:"tags" gorm:"type:jsonb"`
	SEOTitle          string         `json:"seo_title" gorm:"size:255"`
	SEOKeywords       string         `json:"seo_keywords" gorm:"size:255"`
	SEODescription    string         `json:"seo_description" gorm:"size:500"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// SKU 定义商品规格单元
type SKU struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ProductID   uint           `json:"product_id" gorm:"index;not null"`
	SKUCode     string         `json:"sku_code" gorm:"size:50;uniqueIndex;not null"`
	VariantName string         `json:"variant_name" gorm:"size:255;not null"` // 如 "红色，XL"
	Attributes  Attributes     `json:"attributes" gorm:"type:jsonb"`          // 如 {color: "red", size: "XL"}
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	SalePrice   *float64       `json:"sale_price" gorm:"type:decimal(10,2)"`
	StockQty    int            `json:"stock_qty" gorm:"default:0"`
	Image       *string        `json:"image" gorm:"size:255"`
	Weight      *float64       `json:"weight" gorm:"type:decimal(10,2)"`
	IsDefault   bool           `json:"is_default" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Attributes 定义商品属性
type Attributes map[string]string

// Value 实现 driver.Valuer 接口
func (a Attributes) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口
func (a *Attributes) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为 []byte 失败")
	}
	return json.Unmarshal(b, &a)
}

// Category 定义商品分类
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Slug        string         `json:"slug" gorm:"size:50;uniqueIndex;not null"`
	Description string         `json:"description" gorm:"size:500"`
	Image       *string        `json:"image" gorm:"size:255"`
	ParentID    *uint          `json:"parent_id" gorm:"index"`
	Parent      *Category      `json:"parent" gorm:"foreignKey:ParentID"`
	Children    []Category     `json:"children" gorm:"foreignKey:ParentID"`
	Level       int            `json:"level" gorm:"default:0"` // 0 = 根分类，1 = 二级分类，以此类推
	Sort        int            `json:"sort" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Brand 定义商品品牌
type Brand struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;uniqueIndex;not null"`
	Description string         `json:"description" gorm:"size:500"`
	Logo        *string        `json:"logo" gorm:"size:255"`
	Website     *string        `json:"website" gorm:"size:255"`
	Country     *string        `json:"country" gorm:"size:50"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
