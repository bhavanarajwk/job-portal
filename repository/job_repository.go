package repository

import (
	"job-portal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JobFilter holds optional filter parameters for job queries
type JobFilter struct {
	Location string
	Company  string
	Query    string
}

// JobRepository defines the interface for job data operations
type JobRepository interface {
	Create(job *models.Job) error
	FindByID(id uuid.UUID) (*models.Job, error)
	FindAll(offset, limit int, filter JobFilter) ([]models.Job, int64, error)
	FindByRecruiterID(recruiterID uuid.UUID, offset, limit int) ([]models.Job, int64, error)
	Update(job *models.Job) error
	Delete(id uuid.UUID) error
}

type jobRepository struct {
	db *gorm.DB
}

// NewJobRepository creates a new JobRepository instance
func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(job *models.Job) error {
	return r.db.Create(job).Error
}

func (r *jobRepository) FindByID(id uuid.UUID) (*models.Job, error) {
	var job models.Job
	err := r.db.Preload("Recruiter").Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) FindAll(offset, limit int, filter JobFilter) ([]models.Job, int64, error) {
	var jobs []models.Job
	var total int64

	query := r.db.Model(&models.Job{})

	if filter.Query != "" {
		like := "%" + filter.Query + "%"
		query = query.Where(
			"(title ILIKE ? OR description ILIKE ? OR company ILIKE ? OR location ILIKE ?)",
			like, like, like, like,
		)
	}

	if filter.Location != "" {
		query = query.Where("location ILIKE ?", "%"+filter.Location+"%")
	}
	if filter.Company != "" {
		query = query.Where("company ILIKE ?", "%"+filter.Company+"%")
	}

	query.Count(&total)
	err := query.Preload("Recruiter").Offset(offset).Limit(limit).Order("created_at DESC").Find(&jobs).Error
	return jobs, total, err
}

func (r *jobRepository) FindByRecruiterID(recruiterID uuid.UUID, offset, limit int) ([]models.Job, int64, error) {
	var jobs []models.Job
	var total int64

	r.db.Model(&models.Job{}).Where("recruiter_id = ?", recruiterID).Count(&total)
	err := r.db.Where("recruiter_id = ?", recruiterID).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&jobs).Error
	return jobs, total, err
}

func (r *jobRepository) Update(job *models.Job) error {
	return r.db.Save(job).Error
}

func (r *jobRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.Job{}).Error
}
