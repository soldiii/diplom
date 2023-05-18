package service

import (
	"time"

	"github.com/soldiii/diplom/internal/model"
	"github.com/soldiii/diplom/internal/repository"
)

type PlanService struct {
	repo     repository.Plan
	infoRepo repository.Information
}

func NewPlanService(repo repository.Plan, infoRepo repository.Information) *PlanService {
	return &PlanService{repo: repo, infoRepo: infoRepo}
}

func (s *PlanService) GetPlanBySupervisorID(supID int) ([]*repository.PlanStructure, error) {
	return s.repo.GetPlanBySupervisorID(supID)
}

func (s *PlanService) CreatePlan(plan *model.Plan) (int, error) {
	plan.DateTime = time.Now()
	flag, err := s.repo.IsPlanWasCreatedByThisMonth(plan)
	if err != nil {
		return 0, err
	}
	if flag {
		return s.repo.UpdatePlan(plan)
	}
	return s.repo.CreatePlan(plan)
}
