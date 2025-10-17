package dashboard

import (
	"fmt"
	"time"
)

type AdminDashboardServiceInterface interface {
	CardCustomer() (map[string]interface{}, error)
	CardPacketPopular() (map[string]interface{}, error)
	CardAreaPopular() (map[string]interface{}, error)
	CardReportCash() (map[string]interface{}, error)
	GetTotalIncome() (int64, error)
	GetTotalExpenses() (int64, error)
	GetTotalCustomer() (int64, error)
	GetDashboardStats(start *time.Time, end *time.Time) (map[string]interface{}, error)
	GetRecentInvoices() ([]map[string]interface{}, error)
	GetRecentTransactions() ([]map[string]interface{}, error)
	GetCustomerGrowth(start *time.Time, end *time.Time) (map[string]interface{}, error)
	GetRevenueChart(start *time.Time, end *time.Time) (map[string]interface{}, error)
	GetExpensesChart(start *time.Time, end *time.Time) (map[string]interface{}, error)
	GetUnpaidCustomersChart(start *time.Time, end *time.Time) (map[string]interface{}, error)
	GetUnpaidCustomersList() ([]map[string]interface{}, error)
}

type AdminDashboardServiceStruct struct {
	repository AdminDashboardRepositoryInterface
}

func NewAdminDashboardService(repository AdminDashboardRepositoryInterface) AdminDashboardServiceStruct {
	return AdminDashboardServiceStruct{repository}
}

func (s AdminDashboardServiceStruct) CardCustomer() (map[string]interface{}, error) {
	data, err := s.repository.CardCustomer()
	if err != nil {
		return data, err
	}

	return data, err
}

func (s AdminDashboardServiceStruct) CardPacketPopular() (map[string]interface{}, error) {
	data, err := s.repository.CardPacketPopular()
	if err != nil {
		return data, err
	}

	return data, err
}

func (s AdminDashboardServiceStruct) CardAreaPopular() (map[string]interface{}, error) {
	data, err := s.repository.CardAreaPopular()
	if err != nil {
		return data, err
	}

	return data, err
}

func (s AdminDashboardServiceStruct) CardReportCash() (map[string]interface{}, error) {
	data, err := s.repository.CardReportCash()
	if err != nil {
		return data, err
	}

	return data, err
}

func (s AdminDashboardServiceStruct) GetTotalIncome() (int64, error) {
	data, err := s.repository.GetTotalIncome()
	if err != nil {
		return 0, err
	}

	return data, nil
}

func (s AdminDashboardServiceStruct) GetTotalExpenses() (int64, error) {
	data, err := s.repository.GetTotalExpenses()
	if err != nil {
		return 0, err
	}

	return data, nil
}

func (s AdminDashboardServiceStruct) GetTotalCustomer() (int64, error) {
	data, err := s.repository.GetTotalCustomer()
	if err != nil {
		return 0, err
	}

	return data, nil
}

func (s AdminDashboardServiceStruct) GetDashboardStats(start *time.Time, end *time.Time) (map[string]interface{}, error) {
	data, err := s.repository.GetDashboardStats(start, end)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s AdminDashboardServiceStruct) GetRecentInvoices() ([]map[string]interface{}, error) {
	data, err := s.repository.GetRecentInvoices()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s AdminDashboardServiceStruct) GetRecentTransactions() ([]map[string]interface{}, error) {
	data, err := s.repository.GetRecentTransactions()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s AdminDashboardServiceStruct) GetCustomerGrowth(start *time.Time, end *time.Time) (map[string]interface{}, error) {
	data, err := s.repository.GetCustomerGrowth(start, end)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s AdminDashboardServiceStruct) GetRevenueChart(start *time.Time, end *time.Time) (map[string]interface{}, error) {
	data, err := s.repository.GetRevenueChart(start, end)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s AdminDashboardServiceStruct) GetExpensesChart(start *time.Time, end *time.Time) (map[string]interface{}, error) {
	data, err := s.repository.GetExpensesChart(start, end)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s AdminDashboardServiceStruct) GetUnpaidCustomersChart(start *time.Time, end *time.Time) (map[string]interface{}, error) {
	data, err := s.repository.GetUnpaidCustomersChart(start, end)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s AdminDashboardServiceStruct) GetUnpaidCustomersList() ([]map[string]interface{}, error) {
	fmt.Println("🔧 [DEBUG] GetUnpaidCustomersList: Service method called")

	data, err := s.repository.GetUnpaidCustomersList()
	if err != nil {
		fmt.Printf("🔧 [DEBUG] GetUnpaidCustomersList: Repository error - %v\n", err)
		return nil, err
	}

	fmt.Printf("🔧 [DEBUG] GetUnpaidCustomersList: Repository returned %d records\n", len(data))
	return data, nil
}
