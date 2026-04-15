package repository

import (
	"job-portal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ApplicationRepository defines the interface for application data operations
type ApplicationRepository interface {
	Create(application *models.Application) error
	FindByID(id uuid.UUID) (*models.Application, error)
	FindByUserID(userID uuid.UUID, offset, limit int) ([]models.Application, int64, error)
	FindByJobID(jobID uuid.UUID, offset, limit int) ([]models.Application, int64, error)
	FindAll(offset, limit int) ([]models.Application, int64, error)
	ExistsByUserAndJob(userID, jobID uuid.UUID) (bool, error)
	UpdateStatus(id uuid.UUID, status models.ApplicationStatus) error
}

type applicationRepository struct {
	db *gorm.DB
}

// NewApplicationRepository creates a new ApplicationRepository instance
func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(application *models.Application) error {
	return r.db.Create(application).Error
}

func (r *applicationRepository) FindByID(id uuid.UUID) (*models.Application, error) {
	var app models.Application
	err := r.db.Preload("User").Preload("Job").Where("id = ?", id).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepository) FindByUserID(userID uuid.UUID, offset, limit int) ([]models.Application, int64, error) {
	var apps []models.Application
	var total int64

	r.db.Model(&models.Application{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).
		Preload("Job").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&apps).Error
	return apps, total, err
}

func (r *applicationRepository) FindByJobID(jobID uuid.UUID, offset, limit int) ([]models.Application, int64, error) {
	var apps []models.Application
	var total int64

	r.db.Model(&models.Application{}).Where("job_id = ?", jobID).Count(&total)
	err := r.db.Where("job_id = ?", jobID).
		Preload("User").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&apps).Error
	return apps, total, err
}

func (r *applicationRepository) FindAll(offset, limit int) ([]models.Application, int64, error) {
	var apps []models.Application
	var total int64

	r.db.Model(&models.Application{}).Count(&total)
	err := r.db.Preload("User").Preload("Job").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&apps).Error
	return apps, total, err
}

func (r *applicationRepository) ExistsByUserAndJob(userID, jobID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Application{}).
		Where("user_id = ? AND job_id = ?", userID, jobID).
		Count(&count).Error
	return count > 0, err
}

func (r *applicationRepository) UpdateStatus(id uuid.UUID, status models.ApplicationStatus) error {
	return r.db.Model(&models.Application{}).
		Where("id = ?", id).
		Update("status", status).Error
}
