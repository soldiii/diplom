package service

import "github.com/soldiii/diplom/internal/repository"

type AgentService struct {
	repo     repository.Agent
	infoRepo repository.Information
}

func NewAgentService(repo repository.Agent, infoRepo repository.Information) *AgentService {
	return &AgentService{repo: repo, infoRepo: infoRepo}
}
func (s *AgentService) DeleteAgent(agentID string) (int, error) {
	return s.repo.DeleteAgent(agentID)
}
