package models

import (
	"github.com/google/uuid"
	"time"
)


// UserRelation - модель для создания связи "подписчик - подписка"
type UserRelation struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	FollowerID  uuid.UUID `json:"follower_id" gorm:"type:uuid;not null"`
	FollowingID uuid.UUID `json:"following_id" gorm:"type:uuid;not null"`
	CreatedAt   time.Time
}
