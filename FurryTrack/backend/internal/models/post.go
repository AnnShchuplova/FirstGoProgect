package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Post - модель данных для публикации
type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	AuthorID  uuid.UUID `gorm:"type:uuid;not null"`
	Author    User      `gorm:"foreignKey:AuthorID"`
	PetID     uuid.UUID `gorm:"type:uuid" json:"pet_id"`
	Pet       Pet       `gorm:"foreignKey:PetID"`
	Content   string    `gorm:"type:text;not null"`
	PhotoURL   string    `gorm:"type:varchar(255)"`
	PostType  string    `gorm:"type:varchar(20);default:'regular'" json:"post_type"`
	Price     float64   `gorm:"type:decimal(10,2)"`
	CreatedAt time.Time
	UpdatedAt time.Time

	LikesCount int       `gorm:"default:0"` 
	
	// Связи
	Comments []Comment   `gorm:"foreignKey:PostID"`
	Likes    []PostLike  `gorm:"foreignKey:PostID"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Content string `json:"content" gorm:"type:text;not null"`
	UserID  uuid.UUID  `gorm:"type:uuid;not null"`
	PostID  uuid.UUID   `gorm:"type:uuid;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp without time zone;index"`
}
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type PostLike struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	PostID uuid.UUID `gorm:"type:uuid;not null"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp without time zone;index"`
}