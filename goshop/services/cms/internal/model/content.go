package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ContentType 表示内容类型
type ContentType string

const (
	// ContentTypePage 页面
	ContentTypePage ContentType = "page"
	// ContentTypePost 博文
	ContentTypePost ContentType = "post"
	// ContentTypeBanner 横幅
	ContentTypeBanner ContentType = "banner"
)

// ContentStatus 表示内容状态
type ContentStatus string

const (
	// ContentStatusDraft 草稿
	ContentStatusDraft ContentStatus = "draft"
	// ContentStatusPublished 已发布
	ContentStatusPublished ContentStatus = "published"
	// ContentStatusArchived 已归档
	ContentStatusArchived ContentStatus = "archived"
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

// Content 表示CMS内容
type Content struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Type            ContentType    `json:"type" gorm:"size:20;not null;index"`
	Title           string         `json:"title" gorm:"size:255;not null"`
	Slug            string         `json:"slug" gorm:"size:255;uniqueIndex;not null"`
	Content         string         `json:"content" gorm:"type:text"`
	Excerpt         string         `json:"excerpt" gorm:"size:500"`
	CoverImage      *string        `json:"cover_image" gorm:"size:255"`
	Author          string         `json:"author" gorm:"size:50"`
	AuthorID        uint           `json:"author_id" gorm:"index"`
	Status          ContentStatus  `json:"status" gorm:"size:20;not null;default:'draft'"`
	Tags            StringArray    `json:"tags" gorm:"type:jsonb"`
	Categories      []Category     `json:"categories" gorm:"many2many:content_categories"`
	PublishedAt     *time.Time     `json:"published_at"`
	ViewCount       int            `json:"view_count" gorm:"default:0"`
	IsSticky        bool           `json:"is_sticky" gorm:"default:false"`   // 是否置顶
	SortOrder       int            `json:"sort_order" gorm:"default:0"`      // 排序顺序
	MetaTitle       string         `json:"meta_title" gorm:"size:255"`       // SEO标题
	MetaKeywords    string         `json:"meta_keywords" gorm:"size:255"`    // SEO关键词
	MetaDescription string         `json:"meta_description" gorm:"size:500"` // SEO描述
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// Category 表示内容分类
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Slug        string         `json:"slug" gorm:"size:50;uniqueIndex;not null"`
	Description string         `json:"description" gorm:"size:255"`
	ParentID    *uint          `json:"parent_id" gorm:"index"`
	Parent      *Category      `json:"parent" gorm:"foreignKey:ParentID"`
	Children    []Category     `json:"children" gorm:"foreignKey:ParentID"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Menu 表示导航菜单
type Menu struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:50;not null"`
	Location  string         `json:"location" gorm:"size:50;not null"` // 位置：如header, footer, sidebar
	Items     []MenuItem     `json:"items" gorm:"foreignKey:MenuID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// MenuItem 表示导航菜单项
type MenuItem struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	MenuID    uint           `json:"menu_id" gorm:"index;not null"`
	ParentID  *uint          `json:"parent_id" gorm:"index"`
	Parent    *MenuItem      `json:"parent" gorm:"foreignKey:ParentID"`
	Children  []MenuItem     `json:"children" gorm:"foreignKey:ParentID"`
	Title     string         `json:"title" gorm:"size:50;not null"`
	URL       string         `json:"url" gorm:"size:255;not null"`
	Target    string         `json:"target" gorm:"size:20;default:'_self'"` // _self, _blank
	Icon      *string        `json:"icon" gorm:"size:50"`
	SortOrder int            `json:"sort_order" gorm:"default:0"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	ContentID *uint          `json:"content_id" gorm:"index"` // 关联的内容ID
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Banner 表示广告横幅
type Banner struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:100;not null"`
	Image       string         `json:"image" gorm:"size:255;not null"`
	URL         string         `json:"url" gorm:"size:255"`
	Position    string         `json:"position" gorm:"size:50;not null"` // 位置：如home_top, sidebar
	Description string         `json:"description" gorm:"size:255"`
	StartAt     time.Time      `json:"start_at" gorm:"not null"`
	EndAt       time.Time      `json:"end_at" gorm:"not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
