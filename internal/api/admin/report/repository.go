package report

import (
	"skripsi-be/internal/models/dto"
	"skripsi-be/internal/models/entities"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type AdminReportRepositoryInterface interface {
	InternetReportAdminReportRepository() ([]dto.ReportInternet, error)
}
type AdminReportRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminReportRepository(db *gorm.DB) AdminReportRepositoryStruct {
	return AdminReportRepositoryStruct{db}
}

func (r AdminReportRepositoryStruct) InternetReportAdminReportRepository() ([]dto.ReportInternet, error) {
	customers := []entities.Customer{}
	reports := []dto.ReportInternet{}
	tx := r.db.Preload("Product").Find(&customers)
	if tx.Error != nil {
		return nil, tx.Error
	}

	copier.Copy(&reports, &customers)

	return reports, nil
}
