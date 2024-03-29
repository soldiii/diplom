package service

import (
	"errors"

	"github.com/soldiii/diplom/internal/model"
	"github.com/soldiii/diplom/internal/repository"
)

type AdService struct {
	repo     repository.Advertisement
	infoRepo repository.Information
}

func NewAdService(repo repository.Advertisement, infoRepo repository.Information) *AdService {
	return &AdService{repo: repo, infoRepo: infoRepo}
}

func (s *AdService) CreateAd(ad *model.Advertisement) (int, error) {
	return s.repo.CreateAd(ad)
}

func (s *AdService) GetAdsByUserID(userID int, userRole string) ([]*model.Advertisement, error) {
	var supervisorID int
	switch userRole {
	case "agent", "Agent":
		sup_id, err := s.infoRepo.GetSupervisorIDByAgentID(userID)
		if err != nil {
			return nil, err
		}
		supervisorID = sup_id
	case "supervisor", "Supervisor":
		supervisorID = userID
	}
	flag, err := s.repo.IsSupervisorHaveAds(supervisorID)
	if err != nil {
		return nil, err
	}
	if !flag {
		err = errors.New("объявлений нет")
		return nil, err
	}
	return s.repo.GetAdsBySupervisorID(supervisorID)
}

func (s *AdService) UpdateAd(title string, text string, adID string) (int, error) {
	return s.repo.UpdateAd(title, text, adID)
}

func (s *AdService) DeleteAd(adID string) (int, error) {
	return s.repo.DeleteAd(adID)
}
