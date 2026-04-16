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

// UpdateStatusInput holds data for updating application status
type UpdateStatusInput struct {
	Status models.ApplicationStatus `json:"status" binding:"required,oneof=APPLIED REVIEWED REJECTED ACCEPTED"`
}

// ApplyInput holds data submitted when a candidate applies for a job
type ApplyInput struct {
	CoverLetter string // from multipart form field
	ResumeURL   string // set after file is saved
}

// ApplicationService defines application business logic
type ApplicationService interface {
	Apply(candidateID, jobID uuid.UUID, input ApplyInput) (*models.Application, error)
	GetMyApplications(candidateID uuid.UUID, pagination utils.PaginationParams) (utils.PaginatedResponse, error)
	GetApplicationsByJob(jobID uuid.UUID, recruiterID uuid.UUID, pagination utils.PaginationParams) (utils.PaginatedResponse, error)
	GetAllApplications(pagination utils.PaginationParams) (utils.PaginatedResponse, error)
	UpdateStatus(appID uuid.UUID, recruiterID uuid.UUID, input UpdateStatusInput) (*models.Application, error)
}

type applicationService struct {
	appRepo  repository.ApplicationRepository
	jobRepo  repository.JobRepository
	userRepo repository.UserRepository
}

// NewApplicationService creates a new ApplicationService instance
func NewApplicationService(
	appRepo repository.ApplicationRepository,
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
) ApplicationService {
	return &applicationService{appRepo: appRepo, jobRepo: jobRepo, userRepo: userRepo}
}

func (s *applicationService) Apply(candidateID, jobID uuid.UUID, input ApplyInput) (*models.Application, error) {
	// Verify job exists
	job, err := s.jobRepo.FindByID(jobID)
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
		UserID:      candidateID,
		JobID:       jobID,
		CoverLetter: input.CoverLetter,
		ResumeURL:   input.ResumeURL,
		Status:      models.StatusApplied,
	}

	if err := s.appRepo.Create(app); err != nil {
		return nil, errors.New("failed to submit application")
	}

	// Snapshot values for goroutine
	jobTitle    := job.Title
	jobCompany  := job.Company
	jobLocation := job.Location
	recruiterID := job.RecruiterID
	candID      := candidateID

	go func() {
		candidate, err := s.userRepo.FindByID(candID)
		if err != nil {
			log.Printf("[EMAIL] Could not find candidate %s: %v", candID, err)
			return
		}
		recruiter, err := s.userRepo.FindByID(recruiterID)
		if err != nil {
			log.Printf("[EMAIL] Could not find recruiter %s: %v", recruiterID, err)
			return
		}

		// Email 1 → Candidate: confirmation
		utils.SendEmailAsync(utils.ApplicationConfirmationEmail(
			candidate.Email, candidate.Name,
			jobTitle, jobCompany, jobLocation,
		))

		// Email 2 → Recruiter: new application alert
		utils.SendEmailAsync(utils.ApplicationReceivedEmail(
			recruiter.Email, recruiter.Name,
			candidate.Name, candidate.Email,
			jobTitle,
		))
	}()

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

	// Notify candidate on meaningful status changes
	if input.Status == models.StatusAccepted ||
		input.Status == models.StatusRejected ||
		input.Status == models.StatusReviewed {

		jobTitle   := job.Title
		jobCompany := job.Company
		status     := string(input.Status)
		candID     := app.UserID

		go func() {
			candidate, err := s.userRepo.FindByID(candID)
			if err != nil {
				log.Printf("[EMAIL] Could not find candidate %s: %v", candID, err)
				return
			}
			utils.SendEmailAsync(utils.ApplicationStatusEmail(
				candidate.Email, candidate.Name,
				jobTitle, jobCompany, status,
			))
		}()
	}

	return app, nil
}
