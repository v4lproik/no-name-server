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

func (interactor *AppInteractor) FindReport(id string) (domain.Report, error) {
	report, err := interactor.ReportRepository.Find(id)
	return report, err
}

func (interactor *AppInteractor) CreateReport(report *domain.Report) error {
	return interactor.ReportRepository.Store(report)

}

func (interactor *AppInteractor) UpdateReport(report *domain.Report) error {
	return interactor.ReportRepository.Update(report)
}

func (interactor *AppInteractor) DeleteReport(id string) error {
	return interactor.ReportRepository.Delete(id)
}