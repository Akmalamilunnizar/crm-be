package dashboard

import (
	midtrans "skripsi-be/internal/api/common/midtrands"
	"skripsi-be/internal/models/dto"
)

type CustomerDashboardServiceInterface interface {
	MyUserCustomerDashboard(request string) (dto.DashboardDTO, error)
	CreatePaymentCustomerDashboard(request SearchInvoice) (midtrans.MidtransResponsePaymentLink, error)
}

type CustomerDashboardServiceStruct struct {
	repository CustomerDashboardRepositoryInterface
}

func NewCustomerDashboardService(repository CustomerDashboardRepositoryInterface) CustomerDashboardServiceStruct {
	return CustomerDashboardServiceStruct{repository}
}

func (s CustomerDashboardServiceStruct) MyUserCustomerDashboard(request string) (dto.DashboardDTO, error) {
	dashboard := dto.DashboardDTO{}
	myAccount, err := s.repository.MyUserCustomerDashboard(request)
	if err != nil {
		return dashboard, err
	}

	MyProduct, err := s.repository.MyProductCustomerDashboard(myAccount.ProductID)
	if err != nil {
		return dashboard, err
	}

	MyInvoice, err := s.repository.MyInvoiceCustomerDashboard(myAccount.ID)
	if err != nil {
		return dashboard, err
	}

	dashboard.Customer = myAccount
	dashboard.Product = MyProduct
	dashboard.Invoice = MyInvoice

	return dashboard, nil
}

func (s CustomerDashboardServiceStruct) CreatePaymentCustomerDashboard(request SearchInvoice) (midtrans.MidtransResponsePaymentLink, error) {
	midtransResponse := midtrans.MidtransResponsePaymentLink{}
	invoice, err := s.repository.MyInvoiceIdCustomerDashboard(request)
	if err != nil {
		return midtransResponse, err
	}

	mindtransRequest := midtrans.MidtransCreatePaymentLink{
		OrderID:     invoice.ID,
		GrossAmount: invoice.Amount,
	}
	midtransResponse, err = midtrans.CreatePaymentLink(mindtransRequest)

	if err != nil {
		return midtransResponse, err
	}

	// Update invoice with payment link
	_, err = s.repository.MyInvoiceUpdatePaymentCustomerDashboard(request, midtransResponse.RedirectURL)
	if err != nil {
		return midtransResponse, err
	}

	return midtransResponse, nil

}
