package service

import "github.com/soldiii/diplom/internal/repository"

type PlanService struct {
	repo     repository.Plan
	infoRepo repository.Information
}

func NewPlanService(repo repository.Plan, infoRepo repository.Information) *PlanService {
	return &PlanService{repo: repo, infoRepo: infoRepo}
}
