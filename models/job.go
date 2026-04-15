package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Job represents the jobs table
type Job struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text;not null" json:"description"`
	Company     string         `gorm:"type:varchar(255);not null" json:"company"`
	Location    string         `gorm:"type:varchar(255);not null" json:"location"`
	RecruiterID uuid.UUID      `gorm:"type:uuid;not null" json:"recruiter_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	Recruiter    User          `gorm:"foreignKey:RecruiterID" json:"recruiter,omitempty"`
	Applications []Application `gorm:"foreignKey:JobID" json:"-"`
}

// BeforeCreate hook sets UUID before inserting
func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}
