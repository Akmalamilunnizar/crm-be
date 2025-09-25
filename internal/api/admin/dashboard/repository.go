package dashboard

import (
	"time"

	"gorm.io/gorm"

	"skripsi-be/internal/models/entities"
)

type AdminDashboardRepositoryInterface interface {
	CardCustomer() (map[string]interface{}, error)
	CardPacketPopular() (map[string]interface{}, error)
	CardAreaPopular() (map[string]interface{}, error)
	CardReportCash() (map[string]interface{}, error)
	GetTotalIncome() (int64, error)
	GetTotalExpenses() (int64, error)
	GetTotalCustomer() (int64, error)
	GetDashboardStats() (map[string]interface{}, error)
	GetRecentInvoices() ([]map[string]interface{}, error)
	GetRecentTransactions() ([]map[string]interface{}, error)
	GetCustomerGrowth() (map[string]interface{}, error)
	GetRevenueChart() (map[string]interface{}, error)
}

type AdminDashboardRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminDashboradRepository(db *gorm.DB) AdminDashboardRepositoryStruct {
	return AdminDashboardRepositoryStruct{db}
}

type DailyCount struct {
	Date  time.Time
	Count int
}

type PacketCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type AreaCount struct {
	NameCity        string `json:"name_city"`
	NameSubdistrict string `json:"name_subdistrict"`
	NameVillage     string `json:"name_village"`
	Count           int    `json:"count"`
}

func (r AdminDashboardRepositoryStruct) CardCustomer() (map[string]interface{}, error) {
	var total int64
	if err := r.db.Model(&entities.Customer{}).Count(&total).Error; err != nil {
		return nil, err
	}

	var results []DailyCount
	err := r.db.
		Model(&entities.Customer{}).
		Select("DATE(createdAt) as date, COUNT(*) as count").
		Where("createdAt >= ?", time.Now().AddDate(0, 0, -6).Truncate(24*time.Hour)).
		Group("DATE(createdAt)").
		Order("date").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Map the results by date
	dateMap := make(map[string]int)
	for _, r := range results {
		dateMap[r.Date.Format("2006-01-02")] = r.Count
	}

	// Build 7-day graph data
	graphData := make([]map[string]interface{}, 0)
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		graphData = append(graphData, map[string]interface{}{
			"date":  date,
			"count": dateMap[date],
		})
	}

	// Final result
	data := map[string]interface{}{
		"total_customer": total,
		"graph_customer": graphData,
	}

	return data, nil
}

func (r AdminDashboardRepositoryStruct) CardPacketPopular() (map[string]interface{}, error) {
	var total int64
	if err := r.db.Model(&entities.Customer{}).Count(&total).Error; err != nil {
		return nil, err
	}

	var results []PacketCount
	// Join products through network_devices as product is linked to device, not directly on customer
	err := r.db.
		Model(&entities.NetworkDevice{}).
		Joins("JOIN products ON products.id = network_devices.product_id").
		Select("COUNT(*) as count, products.name").
		Group("products.id, products.name").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Final result
	data := map[string]interface{}{
		"total":                total,
		"graph_packet_popular": results,
	}

	return data, nil
}

func (r AdminDashboardRepositoryStruct) CardAreaPopular() (map[string]interface{}, error) {
	var total int64
	if err := r.db.Model(&entities.Customer{}).Count(&total).Error; err != nil {
		return nil, err
	}

	var results []AreaCount
	err := r.db.
		Model(&entities.Customer{}).
		Joins("JOIN areas ON areas.id = customer.area_id").
		Select("COUNT(*) as count, areas.name_city,areas.name_subdistrict,areas.name_village").
		Group("areas.id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Final result
	data := map[string]interface{}{
		"total":              total,
		"graph_area_popular": results,
	}

	return data, nil
}

func (r AdminDashboardRepositoryStruct) CardReportCash() (map[string]interface{}, error) {
	// card report cash by date

	type ReportCash struct {
		Date   time.Time `json:"date"`
		Amount int       `json:"amount"`
	}

	var Transactions []ReportCash
	if err := r.db.Model(&entities.Transaction{}).
		Group("DATE(date)").
		Select("DATE(date) as date, SUM(amount) as amount").
		Order("date DESC").
		Where("type_in_out = ?", entities.TransactionsTypeInOutIn).
		Find(&Transactions).Error; err != nil {
		return nil, err
	}
	// Map the results by date
	dateMap := make(map[string]int)
	for _, r := range Transactions {
		dateMap[r.Date.Format("2006-01-02")] = int(r.Amount)
	}

	// Group transactions by date
	graphData := make([]map[string]interface{}, 0)
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		graphData = append(graphData, map[string]interface{}{
			"date":  date,
			"count": dateMap[date],
		})
	}

	data := map[string]interface{}{
		"graph_report_cash": graphData,
		"total":             0,
	}
	return data, nil
}

func (r AdminDashboardRepositoryStruct) GetTotalIncome() (int64, error) {
	var totalIncome int64
	if err := r.db.Model(&entities.Transaction{}).
		Where("type_in_out = ?", entities.TransactionsTypeInOutIn).
		Select("SUM(amount)").
		Scan(&totalIncome).Error; err != nil {
		return 0, err
	}

	return totalIncome, nil
}

func (r AdminDashboardRepositoryStruct) GetTotalExpenses() (int64, error) {
	var totalIncome int64
	if err := r.db.Model(&entities.Transaction{}).
		Where("type_in_out = ?", entities.TransactionsTypeInOutOut).
		Select("SUM(amount)").
		Scan(&totalIncome).Error; err != nil {
		return 0, err
	}

	return totalIncome, nil
}

func (r AdminDashboardRepositoryStruct) GetTotalCustomer() (int64, error) {
	var total int64
	if err := r.db.Model(&entities.Customer{}).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r AdminDashboardRepositoryStruct) GetDashboardStats() (map[string]interface{}, error) {
	// Get total customers
	var totalCustomers int64
	if err := r.db.Model(&entities.Customer{}).Count(&totalCustomers).Error; err != nil {
		return nil, err
	}

	// Get total income
	var totalIncome int64
	if err := r.db.Model(&entities.Transaction{}).
		Where("type_in_out = ?", entities.TransactionsTypeInOutIn).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome).Error; err != nil {
		return nil, err
	}

	// Get total expenses
	var totalExpenses int64
	if err := r.db.Model(&entities.Transaction{}).
		Where("type_in_out = ?", entities.TransactionsTypeInOutOut).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpenses).Error; err != nil {
		return nil, err
	}

	// Get total invoices
	var totalInvoices int64
	if err := r.db.Model(&entities.Invoice{}).Count(&totalInvoices).Error; err != nil {
		return nil, err
	}

	// Get total areas
	var totalAreas int64
	if err := r.db.Model(&entities.Areas{}).Count(&totalAreas).Error; err != nil {
		return nil, err
	}

	// Get total products
	var totalProducts int64
	if err := r.db.Model(&entities.Products{}).Count(&totalProducts).Error; err != nil {
		return nil, err
	}

	// Get total tickets
	var totalTickets int64
	if err := r.db.Model(&entities.TroubleTicket{}).Count(&totalTickets).Error; err != nil {
		return nil, err
	}

	// Calculate net worth
	netWorth := totalIncome - totalExpenses

	data := map[string]interface{}{
		"total_customers": totalCustomers,
		"total_income":    totalIncome,
		"total_expenses":  totalExpenses,
		"net_worth":       netWorth,
		"total_invoices":  totalInvoices,
		"total_areas":     totalAreas,
		"total_products":  totalProducts,
		"total_tickets":   totalTickets,
	}

	return data, nil
}

func (r AdminDashboardRepositoryStruct) GetRecentInvoices() ([]map[string]interface{}, error) {
	// First, let's check if invoices table exists and has data
	var totalInvoices int64
	if err := r.db.Model(&entities.Invoice{}).Count(&totalInvoices).Error; err != nil {
		return nil, err
	}

	// If no invoices, return empty array instead of error
	if totalInvoices == 0 {
		return []map[string]interface{}{}, nil
	}

	var invoices []entities.Invoice
	if err := r.db.Preload("Customer").
		Order("createdAt DESC").
		Limit(10).
		Find(&invoices).Error; err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, invoice := range invoices {
		customerName := ""
		if invoice.Customer.Name != "" {
			customerName = invoice.Customer.Name
		}

		result = append(result, map[string]interface{}{
			"id":         invoice.ID,
			"invoice_no": invoice.ID, // Using ID as invoice number since InvoiceNo field doesn't exist
			"amount":     invoice.Amount,
			"status":     invoice.Status,
			"customer":   customerName,
			"created_at": invoice.CreatedAt,
		})
	}

	return result, nil
}

func (r AdminDashboardRepositoryStruct) GetRecentTransactions() ([]map[string]interface{}, error) {
	// First, let's check if transactions table exists and has data
	var totalTransactions int64
	if err := r.db.Model(&entities.Transaction{}).Count(&totalTransactions).Error; err != nil {
		return nil, err
	}

	// If no transactions, return empty array instead of error
	if totalTransactions == 0 {
		return []map[string]interface{}{}, nil
	}

	var transactions []entities.Transaction
	if err := r.db.Order("date DESC").
		Limit(10).
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, transaction := range transactions {
		result = append(result, map[string]interface{}{
			"id":          transaction.ID,
			"amount":      transaction.Amount,
			"type_in_out": transaction.TypeInOut,
			"description": transaction.Description,
			"date":        transaction.Date,
			"created_at":  transaction.CreatedAt,
		})
	}

	return result, nil
}

func (r AdminDashboardRepositoryStruct) GetCustomerGrowth() (map[string]interface{}, error) {
	// First, let's check if customers table exists and has data
	var totalCustomers int64
	if err := r.db.Model(&entities.Customer{}).Count(&totalCustomers).Error; err != nil {
		return nil, err
	}

	// If no customers, return empty data instead of error
	if totalCustomers == 0 {
		data := map[string]interface{}{
			"customer_growth": []map[string]interface{}{},
		}
		return data, nil
	}

	var results []DailyCount
	err := r.db.
		Model(&entities.Customer{}).
		Select("DATE(createdAt) as date, COUNT(*) as count").
		Where("createdAt >= ?", time.Now().AddDate(0, -3, 0).Truncate(24*time.Hour)). // Changed to 3 months to capture existing customers
		Group("DATE(createdAt)").
		Order("date").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Map the results by date
	dateMap := make(map[string]int)
	for _, r := range results {
		dateMap[r.Date.Format("2006-01-02")] = r.Count
	}

	// Build 90-day graph data (3 months)
	graphData := make([]map[string]interface{}, 0)
	for i := 89; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		graphData = append(graphData, map[string]interface{}{
			"date":  date,
			"count": dateMap[date],
		})
	}

	data := map[string]interface{}{
		"customer_growth": graphData,
	}

	return data, nil
}

func (r AdminDashboardRepositoryStruct) GetRevenueChart() (map[string]interface{}, error) {
	// First, check if there are any transactions
	var totalTransactions int64
	if err := r.db.Model(&entities.Transaction{}).Where("type_in_out = ?", entities.TransactionsTypeInOutIn).Count(&totalTransactions).Error; err != nil {
		return nil, err
	}

	// If no transactions, return empty data
	if totalTransactions == 0 {
		data := map[string]interface{}{
			"revenue_chart": []map[string]interface{}{},
		}
		return data, nil
	}

	var results []DailyCount
	err := r.db.
		Model(&entities.Transaction{}).
		Select("DATE(date) as date, SUM(amount) as count").
		Where("type_in_out = ? AND date >= ?", entities.TransactionsTypeInOutIn, time.Now().AddDate(0, -3, 0).Truncate(24*time.Hour)). // Changed to 3 months
		Group("DATE(date)").
		Order("date").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Map the results by date
	dateMap := make(map[string]int)
	for _, r := range results {
		dateMap[r.Date.Format("2006-01-02")] = r.Count
	}

	// Build 90-day graph data (3 months)
	graphData := make([]map[string]interface{}, 0)
	for i := 89; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		graphData = append(graphData, map[string]interface{}{
			"date":   date,
			"amount": dateMap[date],
		})
	}

	data := map[string]interface{}{
		"revenue_chart": graphData,
	}

	return data, nil
}
