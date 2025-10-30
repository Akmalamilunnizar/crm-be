package dashboard

import (
	"fmt"
	midtrans "skripsi-be/internal/api/common/midtrands"
	"skripsi-be/internal/models/dto"
	"skripsi-be/internal/models/entities"
)

type CustomerDashboardServiceInterface interface {
	MyUserCustomerDashboard(request string) (dto.DashboardDTO, error)
	CreatePaymentCustomerDashboard(request SearchInvoice) (midtrans.MidtransResponsePaymentLink, error)
	CheckDeviceStatus(customerID string) (string, error)
	GetAvailableProducts() ([]entities.Products, error)
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
		fmt.Printf("CustomerDashboard - Error getting customer: %v\n", err)
		return dashboard, err
	}

	fmt.Printf("CustomerDashboard - Customer ID: %s, Name: %s\n", myAccount.ID, myAccount.Name)

	// Get network device data with mac_address and product info
	networkDevice, err := s.repository.GetNetworkDeviceData(myAccount.ID)
	if err != nil {
		fmt.Printf("CustomerDashboard - Error getting network device: %v\n", err)
		return dashboard, err
	}

	fmt.Printf("CustomerDashboard - Network Device ID: %s, MAC: %v, ProductID: %v\n",
		networkDevice.ID, networkDevice.MacAddress, networkDevice.ProductID)

	// If no network device found, create empty product
	var MyProduct entities.Products
	if networkDevice.ID != "" && networkDevice.Product != nil {
		MyProduct = *networkDevice.Product
		fmt.Printf("CustomerDashboard - Product found: ID=%s, Name=%s, Price=%v\n",
			MyProduct.ID, MyProduct.Name, MyProduct.Price)
	} else {
		fmt.Printf("CustomerDashboard - No product found. NetworkDevice.ID=%s, Product=%v\n",
			networkDevice.ID, networkDevice.Product)

		// Fallback: Try to get product using the old method
		productID, err := s.repository.GetProductIDFromNetworkDevice(myAccount.ID)
		if err == nil && productID != "" {
			fmt.Printf("CustomerDashboard - Fallback: Found product_id=%s\n", productID)
			MyProduct, err = s.repository.MyProductCustomerDashboard(productID)
			if err != nil {
				fmt.Printf("CustomerDashboard - Fallback: Error getting product: %v\n", err)
			} else {
				fmt.Printf("CustomerDashboard - Fallback: Product found: ID=%s, Name=%s, Price=%v\n",
					MyProduct.ID, MyProduct.Name, MyProduct.Price)
			}
		} else {
			fmt.Printf("CustomerDashboard - Fallback: No product_id found\n")

			// Additional fallback: Try to get customer with product relationship
			_, err = s.repository.GetCustomerWithProduct(myAccount.ID)
			if err == nil {
				fmt.Printf("CustomerDashboard - Customer found but no product relationship available\n")
			} else {
				fmt.Printf("CustomerDashboard - No customer found: %v\n", err)
			}
		}
	}

	MyInvoice, err := s.repository.MyInvoiceCustomerDashboard(myAccount.ID)
	if err != nil {
		fmt.Printf("CustomerDashboard - Error getting invoices: %v\n", err)
		return dashboard, err
	}

	fmt.Printf("CustomerDashboard - Found %d invoices\n", len(MyInvoice))

	dashboard.Customer = myAccount
	dashboard.Product = MyProduct
	dashboard.Invoice = MyInvoice
	// Add network device data for frontend to use for status checking
	dashboard.NetworkDevice = networkDevice

	fmt.Printf("CustomerDashboard - Final dashboard: Customer=%s, Product=%s, NetworkDevice=%s\n",
		dashboard.Customer.Name, dashboard.Product.Name, dashboard.NetworkDevice.ID)

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

func (s CustomerDashboardServiceStruct) CheckDeviceStatus(customerID string) (string, error) {
	// Get network device data for the customer
	networkDevice, err := s.repository.GetNetworkDeviceData(customerID)
	if err != nil {
		return "off", err
	}

	// If no network device found, return "off"
	if networkDevice.ID == "" {
		return "off", nil
	}

	// If no mac_address, return "off"
	if networkDevice.MacAddress == nil || *networkDevice.MacAddress == "" {
		return "off", nil
	}

	// Check device status via Mikrotik
	status, err := s.repository.CheckDeviceStatus(*networkDevice.MacAddress)
	if err != nil {
		return "unknown", err
	}

	return status, nil
}

func (s CustomerDashboardServiceStruct) GetAvailableProducts() ([]entities.Products, error) {
	fmt.Printf("CustomerDashboard - Getting available products\n")

	products, err := s.repository.GetAvailableProducts()
	if err != nil {
		fmt.Printf("CustomerDashboard - Error getting products: %v\n", err)
		return products, err
	}

	fmt.Printf("CustomerDashboard - Found %d available products\n", len(products))
	return products, nil
}
