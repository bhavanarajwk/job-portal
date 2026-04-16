package services

import (
	"errors"
	"log"

	"job-portal/models"
	"job-portal/repository"
	"job-portal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateJobInput holds data for creating a job
type CreateJobInput struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Description string `json:"description" binding:"required,min=10"`
	Company     string `json:"company" binding:"required,min=2,max=255"`
	Location    string `json:"location" binding:"required,min=2,max=255"`
}

// UpdateJobInput holds data for updating a job
type UpdateJobInput struct {
	Title       string `json:"title" binding:"omitempty,min=3,max=255"`
	Description string `json:"description" binding:"omitempty,min=10"`
	Company     string `json:"company" binding:"omitempty,min=2,max=255"`
	Location    string `json:"location" binding:"omitempty,min=2,max=255"`
}

// JobService defines job business logic
type JobService interface {
	CreateJob(input CreateJobInput, recruiterID uuid.UUID) (*models.Job, error)
	GetJobByID(id uuid.UUID) (*models.Job, error)
	GetAllJobs(pagination utils.PaginationParams, filter repository.JobFilter) (utils.PaginatedResponse, error)
	UpdateJob(id uuid.UUID, input UpdateJobInput, recruiterID uuid.UUID) (*models.Job, error)
	DeleteJob(id uuid.UUID, recruiterID uuid.UUID, role models.Role) error
}

type jobService struct {
	jobRepo  repository.JobRepository
	userRepo repository.UserRepository
	appRepo  repository.ApplicationRepository
}

// NewJobService creates a new JobService instance
func NewJobService(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	appRepo repository.ApplicationRepository,
) JobService {
	return &jobService{jobRepo: jobRepo, userRepo: userRepo, appRepo: appRepo}
}

// ── Create ────────────────────────────────────────────────────────────────────

func (s *jobService) CreateJob(input CreateJobInput, recruiterID uuid.UUID) (*models.Job, error) {
	job := &models.Job{
		Title:       input.Title,
		Description: input.Description,
		Company:     input.Company,
		Location:    input.Location,
		RecruiterID: recruiterID,
	}

	if err := s.jobRepo.Create(job); err != nil {
		return nil, errors.New("failed to create job")
	}

	// Notify all candidates asynchronously
	go s.notifyAllCandidatesNewJob(job)

	return job, nil
}

// notifyAllCandidatesNewJob emails every CANDIDATE about the new job posting
func (s *jobService) notifyAllCandidatesNewJob(job *models.Job) {
	candidates, err := s.userRepo.FindAllByRole(models.RoleCandidate)
	if err != nil {
		log.Printf("[EMAIL] Failed to fetch candidates for new job alert: %v", err)
		return
	}
	if len(candidates) == 0 {
		log.Println("[EMAIL] No candidates registered to notify")
		return
	}
	log.Printf("[EMAIL] Notifying %d candidate(s) about new job: %s", len(candidates), job.Title)
	for _, c := range candidates {
		utils.SendEmailAsync(utils.NewJobAlertEmail(
			c.Email, c.Name,
			job.Title, job.Company, job.Location, job.Description,
		))
	}
}

// ── Read ──────────────────────────────────────────────────────────────────────

func (s *jobService) GetJobByID(id uuid.UUID) (*models.Job, error) {
	job, err := s.jobRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job not found")
		}
		return nil, errors.New("database error")
	}
	return job, nil
}

func (s *jobService) GetAllJobs(pagination utils.PaginationParams, filter repository.JobFilter) (utils.PaginatedResponse, error) {
	jobs, total, err := s.jobRepo.FindAll(pagination.Offset, pagination.PageSize, filter)
	if err != nil {
		return utils.PaginatedResponse{}, errors.New("failed to fetch jobs")
	}
	return utils.BuildPaginatedResponse(jobs, total, pagination), nil
}

// ── Update ────────────────────────────────────────────────────────────────────

func (s *jobService) UpdateJob(id uuid.UUID, input UpdateJobInput, recruiterID uuid.UUID) (*models.Job, error) {
	job, err := s.jobRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job not found")
		}
		return nil, errors.New("database error")
	}

	if job.RecruiterID != recruiterID {
		return nil, errors.New("forbidden: you do not own this job")
	}

	if input.Title != "" {
		job.Title = input.Title
	}
	if input.Description != "" {
		job.Description = input.Description
	}
	if input.Company != "" {
		job.Company = input.Company
	}
	if input.Location != "" {
		job.Location = input.Location
	}

	if err := s.jobRepo.Update(job); err != nil {
		return nil, errors.New("failed to update job")
	}

	// Notify candidates who applied for this job
	go s.notifyApplicantsJobUpdated(job)

	return job, nil
}

// notifyApplicantsJobUpdated emails candidates who have an active application for this job
func (s *jobService) notifyApplicantsJobUpdated(job *models.Job) {
	// Fetch all applications for this job (no pagination — get all)
	apps, _, err := s.appRepo.FindByJobID(job.ID, 0, 10000)
	if err != nil {
		log.Printf("[EMAIL] Failed to fetch applicants for job update: %v", err)
		return
	}
	if len(apps) == 0 {
		log.Printf("[EMAIL] No applicants to notify for job update: %s", job.Title)
		return
	}

	log.Printf("[EMAIL] Notifying %d applicant(s) about job update: %s", len(apps), job.Title)
	for _, app := range apps {
		candidate, err := s.userRepo.FindByID(app.UserID)
		if err != nil {
			log.Printf("[EMAIL] Could not find candidate %s: %v", app.UserID, err)
			continue
		}
		utils.SendEmailAsync(utils.JobUpdatedEmail(
			candidate.Email, candidate.Name,
			job.Title, job.Company, job.Location, job.Description,
		))
	}
}

// ── Delete ────────────────────────────────────────────────────────────────────

func (s *jobService) DeleteJob(id uuid.UUID, recruiterID uuid.UUID, role models.Role) error {
	job, err := s.jobRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("job not found")
		}
		return errors.New("database error")
	}

	if role != models.RoleAdmin && job.RecruiterID != recruiterID {
		return errors.New("forbidden: you do not own this job")
	}

	return s.jobRepo.Delete(id)
}
