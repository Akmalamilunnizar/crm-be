package dashboard

type AdminDashboardServiceInterface interface {
	CardCustomer() (map[string]interface{}, error)
	CardPacketPopular() (map[string]interface{}, error)
	CardAreaPopular() (map[string]interface{}, error)
	CardReportCash() (map[string]interface{}, error)
	GetTotalIncome() (int64, error)
	GetTotalExpenses() (int64, error)
	GetTotalCustomer() (int64, error)
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
