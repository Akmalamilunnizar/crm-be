package dashboard

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"skripsi-be/internal/models/dto"
	"skripsi-be/internal/models/entities"
)

type CustomerDashboardRepositoryInterface interface {
	MyUserCustomerDashboard(request string) (dto.CustomerDTO, error)
	MyProductCustomerDashboard(request string) (entities.Products, error)
	MyInvoiceCustomerDashboard(request string) ([]entities.Invoice, error)
	MyInvoiceIdCustomerDashboard(request SearchInvoice) (entities.Invoice, error)
	MyInvoiceUpdatePaymentCustomerDashboard(request SearchInvoice, link string) (entities.Invoice, error)
	GetProductIDFromNetworkDevice(customerID string) (string, error)
}

type CustomerDashboardRepositoryStruct struct {
	db *gorm.DB
}

func NewCustomerDashboardRepository(db *gorm.DB) CustomerDashboardRepositoryStruct {
	return CustomerDashboardRepositoryStruct{db}
}

func (r CustomerDashboardRepositoryStruct) MyUserCustomerDashboard(request string) (dto.CustomerDTO, error) {
	myAccount := entities.Customer{}
	customerDto := dto.CustomerDTO{}
	err := r.db.Where("id = ?", &request).First(&myAccount).Error
	if err != nil {
		return customerDto, err
	}

	copier.Copy(&customerDto, &myAccount)

	return customerDto, nil
}

func (r CustomerDashboardRepositoryStruct) MyProductCustomerDashboard(request string) (entities.Products, error) {
	product := entities.Products{}
	err := r.db.Where("id = ?", &request).First(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r CustomerDashboardRepositoryStruct) MyInvoiceCustomerDashboard(request string) ([]entities.Invoice, error) {
	invoice := []entities.Invoice{}
	err := r.db.Where("customer_id = ?", &request).Order("createdAt desc").Find(&invoice).Error
	if err != nil {
		return invoice, err
	}

	return invoice, nil
}

func (r CustomerDashboardRepositoryStruct) MyInvoiceIdCustomerDashboard(request SearchInvoice) (entities.Invoice, error) {
	invoice := entities.Invoice{}
	err := r.db.Where("customer_id = ?", request.UserId).Where("id = ?", request.InvoiceId).First(&invoice).Error
	if err != nil {
		return invoice, err
	}

	return invoice, nil
}

func (r CustomerDashboardRepositoryStruct) MyInvoiceUpdatePaymentCustomerDashboard(request SearchInvoice, link string) (entities.Invoice, error) {
	invoice := entities.Invoice{}
	err := r.db.Where("customer_id = ?", request.UserId).Where("id = ?", request.InvoiceId).First(&invoice).Error
	if err != nil {
		return invoice, err
	}
	invoice.Link = link
	invoice.Status = entities.InvoiceStatusPending
	err = r.db.Save(&invoice).Error
	if err != nil {
		return invoice, err
	}

	return invoice, nil
}

func (r CustomerDashboardRepositoryStruct) GetProductIDFromNetworkDevice(customerID string) (string, error) {
	networkDevice := entities.NetworkDevice{}
	err := r.db.Where("customer_id = ?", customerID).First(&networkDevice).Error
	if err != nil {
		// If no network device found, return empty string instead of error
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}

	if networkDevice.ProductID == nil {
		return "", nil
	}

	return *networkDevice.ProductID, nil
}
