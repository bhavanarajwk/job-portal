package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role defines the user role type
type Role string

const (
	RoleAdmin     Role = "ADMIN"
	RoleRecruiter Role = "RECRUITER"
	RoleCandidate Role = "CANDIDATE"
)

// User represents the users table
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Role      Role           `gorm:"type:varchar(20);not null;default:'CANDIDATE'" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	Jobs         []Job         `gorm:"foreignKey:RecruiterID" json:"-"`
	Applications []Application `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook sets UUID before inserting
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
