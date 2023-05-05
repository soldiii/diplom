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
	return s.repo.CreateReport(report)
}
