package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ApplicationStatus defines the status of a job application
type ApplicationStatus string

const (
	StatusApplied   ApplicationStatus = "APPLIED"
	StatusReviewed  ApplicationStatus = "REVIEWED"
	StatusRejected  ApplicationStatus = "REJECTED"
	StatusAccepted  ApplicationStatus = "ACCEPTED"
)

// Application represents the applications table
type Application struct {
	ID        uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID         `gorm:"type:uuid;not null" json:"user_id"`
	JobID     uuid.UUID         `gorm:"type:uuid;not null" json:"job_id"`
	Status    ApplicationStatus `gorm:"type:varchar(20);not null;default:'APPLIED'" json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt gorm.DeletedAt    `gorm:"index" json:"-"`

	// Associations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Job  Job  `gorm:"foreignKey:JobID" json:"job,omitempty"`
}

// BeforeCreate hook sets UUID before inserting
func (a *Application) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
