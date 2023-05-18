package service

import (
	"errors"

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

func (s *InfoService) GetIsValidByID(uID int) (bool, error) {
	return s.repo.GetIsValidByID(uID)
}

type InfoAboutAgent struct {
	FullName           string
	SupervisorFullName string
	Report             *repository.Rates
	Plan               *repository.Rates
}

func (s *InfoService) GetInfoAboutAgentByID(agentID int) (*InfoAboutAgent, error) {

	fullName, err := s.repo.GetFullNameByAgentID(agentID)
	if err != nil {
		return nil, err
	}
	supervisorFullName, err := s.repo.GetSupervisorFullNameByAgentID(agentID)
	if err != nil {
		return nil, err
	}
	report, err := s.repo.GetReportByAgentID(agentID)
	if err != nil {
		return nil, err
	}
	plan, err := s.repo.GetPlanByAgentID(agentID)
	if err != nil {
		return nil, err
	}

	info := &InfoAboutAgent{FullName: fullName, SupervisorFullName: supervisorFullName, Report: report, Plan: plan}
	return info, nil
}

type InfoAboutSupervisor struct {
	FullName string
	Plan     *repository.Rates
}

func (s *InfoService) GetInfoAboutSupervisorByID(supID int) (*InfoAboutSupervisor, error) {
	fullName, err := s.repo.GetFullNameBySupID(supID)
	if err != nil {
		return nil, err
	}
	plan, err := s.repo.GetPlanBySupID(supID)
	if err != nil {
		return nil, err
	}

	info := &InfoAboutSupervisor{FullName: fullName, Plan: plan}
	return info, nil
}

type AgentIDAndFullName struct {
	ID       int
	FullName string
}

func (s *InfoService) GetAllAgentsBySupID(supID int) ([]*repository.AgentIDAndFullName, error) {
	if err := s.repo.CheckForSupervisor(supID); err != nil {
		err = errors.New("супервайзер с таким id не существует")
		return nil, err
	}
	return s.repo.GetAllAgentsBySupID(supID)
}
