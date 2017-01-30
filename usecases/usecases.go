package usecases

import (
	"github.com/yinkozi/no-name-domain"
)

type Logger interface {
	Log(message string) error
}

type ReportRepository interface {
	Store(report *domain.Report) error
	Update(report *domain.Report) error
	Find(id string) (domain.Report, error)
	List() ([]domain.Report, error)
	Delete(id string) error
}

type AppInteractor struct {
	ReportRepository ReportRepository
	Logger          Logger
}

func (u *AppInteractor) FindReport(id string) (domain.Report, error) {
	report, err := u.ReportRepository.Find(id)
	return report, err
}

func (u *AppInteractor) CreateReport(report *domain.Report) error {
	return u.ReportRepository.Store(report)

}

func (u *AppInteractor) UpdateReport(report *domain.Report) error {
	return u.ReportRepository.Update(report)
}

func (r *AppInteractor) DeleteReport(id string) error {
	return r.ReportRepository.Delete(id)
}