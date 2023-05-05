package service

import (
	"github.com/soldiii/diplom/internal/model"
	"github.com/soldiii/diplom/internal/repository"
)

type InfoService struct {
	repo repository.Information
}

func NewInfoService(repo repository.Information) *InfoService {
	return &InfoService{repo: repo}
}

func (s *InfoService) GetAllSupervisors() ([]*model.Supervisor, error) {
	return s.repo.GetAllSupervisors()
}

func (s *InfoService) GetUserRoleByID(uID string) (string, error) {
	return s.repo.GetUserRoleByID(uID)
}

type InfoAboutAgent struct {
	FullName           string
	SupervisorFullName string
	Report             *repository.Rates
	Plan               *repository.Rates
}

func (s *InfoService) GetInfoAboutAgentByID(agentID string) (*InfoAboutAgent, error) {

	fullName, err := s.repo.GetFullNameByID(agentID)
	if err != nil {
		return nil, err
	}
	supervisorFullName, err := s.repo.GetSupervisorFullNameByID(agentID)
	if err != nil {
		return nil, err
	}
	report, err := s.repo.GetReportByID(agentID)
	if err != nil {
		return nil, err
	}
	plan, err := s.repo.GetPlanByID(agentID)
	if err != nil {
		return nil, err
	}

	info := &InfoAboutAgent{FullName: fullName, SupervisorFullName: supervisorFullName, Report: report, Plan: plan}
	return info, nil
}
