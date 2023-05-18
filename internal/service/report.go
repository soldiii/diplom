package service

import (
	"time"

	"github.com/soldiii/diplom/internal/model"
	"github.com/soldiii/diplom/internal/repository"
)

type ReportService struct {
	repo     repository.Report
	infoRepo repository.Information
}

func NewReportService(repo repository.Report, infoRepo repository.Information) *ReportService {
	return &ReportService{repo: repo, infoRepo: infoRepo}
}

func (s *ReportService) CreateReport(report *model.Report) (int, error) {
	report.DateTime = time.Now()
	flag, err := s.repo.IsReportWasCreatedByThisDay(report)
	if err != nil {
		return 0, err
	}
	if flag {
		return s.repo.UpdateReport(report)
	}
	return s.repo.CreateReport(report)
}

func (s *ReportService) GetRatesByAgentID(agentID int) (*repository.Rates, error) {
	return s.repo.GetRatesByAgentID(agentID)
}

func (s *ReportService) GetRatesBySupervisorIDAndPeriod(supID int, period string) (*repository.Rates, error) {
	return s.repo.GetRatesBySupervisorIDAndPeriod(supID, period)
}

func IncreaseLastDateByOneDay(lastDate string) (string, error) {
	lDate, err := time.Parse("2006-01-02", lastDate)
	if err != nil {
		return "", err
	}
	lDate = lDate.AddDate(0, 0, 1)
	newLastDate := lDate.Format("2006-01-02")
	return newLastDate, nil
}

func (s *ReportService) GetRatesBySupervisorFirstAndLastDates(supID int, firstDate, lastDate string) (*repository.Rates, error) {
	newLastDate, err := IncreaseLastDateByOneDay(lastDate)
	if err != nil {
		return nil, err
	}
	return s.repo.GetRatesBySupervisorFirstAndLastDates(supID, firstDate, newLastDate)
}

func (s *ReportService) GetReportsByAgents(supID int, firstDate, lastDate string) ([]*repository.ReportStructure, error) {
	newLastDate, err := IncreaseLastDateByOneDay(lastDate)
	if err != nil {
		return nil, err
	}
	return s.repo.GetReportsByAgents(supID, firstDate, newLastDate)
}
