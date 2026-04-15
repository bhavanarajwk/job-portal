package services

import (
	"errors"

	"job-portal/models"
	"job-portal/repository"
	"job-portal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UpdateStatusInput holds data for updating application status
type UpdateStatusInput struct {
	Status models.ApplicationStatus `json:"status" binding:"required,oneof=APPLIED REVIEWED REJECTED ACCEPTED"`
}

// ApplicationService defines application business logic
type ApplicationService interface {
	Apply(candidateID, jobID uuid.UUID) (*models.Application, error)
	GetMyApplications(candidateID uuid.UUID, pagination utils.PaginationParams) (utils.PaginatedResponse, error)
	GetApplicationsByJob(jobID uuid.UUID, recruiterID uuid.UUID, pagination utils.PaginationParams) (utils.PaginatedResponse, error)
	GetAllApplications(pagination utils.PaginationParams) (utils.PaginatedResponse, error)
	UpdateStatus(appID uuid.UUID, recruiterID uuid.UUID, input UpdateStatusInput) (*models.Application, error)
}

type applicationService struct {
	appRepo repository.ApplicationRepository
	jobRepo repository.JobRepository
}

// NewApplicationService creates a new ApplicationService instance
func NewApplicationService(appRepo repository.ApplicationRepository, jobRepo repository.JobRepository) ApplicationService {
	return &applicationService{appRepo: appRepo, jobRepo: jobRepo}
}

func (s *applicationService) Apply(candidateID, jobID uuid.UUID) (*models.Application, error) {
	// Verify job exists
	_, err := s.jobRepo.FindByID(jobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job not found")
		}
		return nil, errors.New("database error")
	}

	// Prevent duplicate applications
	exists, err := s.appRepo.ExistsByUserAndJob(candidateID, jobID)
	if err != nil {
		return nil, errors.New("database error")
	}
	if exists {
		return nil, errors.New("you have already applied for this job")
	}

	app := &models.Application{
		UserID: candidateID,
		JobID:  jobID,
		Status: models.StatusApplied,
	}

	if err := s.appRepo.Create(app); err != nil {
		return nil, errors.New("failed to submit application")
	}

	return app, nil
}

func (s *applicationService) GetMyApplications(candidateID uuid.UUID, pagination utils.PaginationParams) (utils.PaginatedResponse, error) {
	apps, total, err := s.appRepo.FindByUserID(candidateID, pagination.Offset, pagination.PageSize)
	if err != nil {
		return utils.PaginatedResponse{}, errors.New("failed to fetch applications")
	}
	return utils.BuildPaginatedResponse(apps, total, pagination), nil
}

func (s *applicationService) GetApplicationsByJob(jobID uuid.UUID, recruiterID uuid.UUID, pagination utils.PaginationParams) (utils.PaginatedResponse, error) {
	// Verify recruiter owns the job
	job, err := s.jobRepo.FindByID(jobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.PaginatedResponse{}, errors.New("job not found")
		}
		return utils.PaginatedResponse{}, errors.New("database error")
	}

	if job.RecruiterID != recruiterID {
		return utils.PaginatedResponse{}, errors.New("forbidden: you do not own this job")
	}

	apps, total, err := s.appRepo.FindByJobID(jobID, pagination.Offset, pagination.PageSize)
	if err != nil {
		return utils.PaginatedResponse{}, errors.New("failed to fetch applications")
	}
	return utils.BuildPaginatedResponse(apps, total, pagination), nil
}

func (s *applicationService) GetAllApplications(pagination utils.PaginationParams) (utils.PaginatedResponse, error) {
	apps, total, err := s.appRepo.FindAll(pagination.Offset, pagination.PageSize)
	if err != nil {
		return utils.PaginatedResponse{}, errors.New("failed to fetch applications")
	}
	return utils.BuildPaginatedResponse(apps, total, pagination), nil
}

func (s *applicationService) UpdateStatus(appID uuid.UUID, recruiterID uuid.UUID, input UpdateStatusInput) (*models.Application, error) {
	app, err := s.appRepo.FindByID(appID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("application not found")
		}
		return nil, errors.New("database error")
	}

	// Verify recruiter owns the job this application is for
	job, err := s.jobRepo.FindByID(app.JobID)
	if err != nil {
		return nil, errors.New("associated job not found")
	}
	if job.RecruiterID != recruiterID {
		return nil, errors.New("forbidden: you do not own this job")
	}

	if err := s.appRepo.UpdateStatus(appID, input.Status); err != nil {
		return nil, errors.New("failed to update status")
	}

	app.Status = input.Status
	return app, nil
}
