package dashboard

import (
	"time"

	"gorm.io/gorm"

	"fmt"
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
	GetCustomerGrowth(start *time.Time, end *time.Time) (map[string]interface{}, error)
	GetRevenueChart(start *time.Time, end *time.Time) (map[string]interface{}, error)
	GetExpensesChart(start *time.Time, end *time.Time) (map[string]interface{}, error)
	GetUnpaidCustomersChart(start *time.Time, end *time.Time) (map[string]interface{}, error)
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
	Name  string `json:"name" gorm:"column:name"`
	Count int    `json:"count" gorm:"column:count"`
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

	// Build a stable local-day window to avoid timezone drift
	now := time.Now().In(time.Local)
	startDay := time.Date(now.Year(), now.Month(), now.Day()-6, 0, 0, 0, 0, time.Local)
	endDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)

	var results []DailyCount
	err := r.db.
		Model(&entities.Customer{}).
		Select("DATE(createdAt) as date, COUNT(*) as count").
		Where("DATE(createdAt) BETWEEN ? AND ?", startDay.Format("2006-01-02"), endDay.Format("2006-01-02")).
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
		date := time.Date(now.Year(), now.Month(), now.Day()-i, 0, 0, 0, 0, time.Local).Format("2006-01-02")
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
		// If no data found, return empty results instead of error
		results = []PacketCount{}
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

	// Limit to last 7 days with explicit local-day window to avoid timezone mismatches
	now := time.Now().In(time.Local)
	startDay := time.Date(now.Year(), now.Month(), now.Day()-6, 0, 0, 0, 0, time.Local)

	var Transactions []ReportCash
	if err := r.db.Model(&entities.Transaction{}).
		Group("DATE(date)").
		Select("DATE(date) as date, SUM(amount) as amount").
		Order("date DESC").
		Where("type_in_out = ? AND DATE(date) >= ?", entities.TransactionsTypeInOutIn, startDay.Format("2006-01-02")).
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
		date := time.Date(now.Year(), now.Month(), now.Day()-i, 0, 0, 0, 0, time.Local).Format("2006-01-02")
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

	// Get total trouble tickets (count of tickets)
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

func (r AdminDashboardRepositoryStruct) GetCustomerGrowth(start *time.Time, end *time.Time) (map[string]interface{}, error) {
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
	query := r.db.Model(&entities.Customer{}).
		Select("DATE(createdAt) as date, COUNT(*) as count")
	if start != nil && end != nil {
		// Expand to full local-day bounds
		s := time.Date(start.In(time.Local).Year(), start.In(time.Local).Month(), start.In(time.Local).Day(), 0, 0, 0, 0, time.Local)
		e := time.Date(end.In(time.Local).Year(), end.In(time.Local).Month(), end.In(time.Local).Day(), 23, 59, 59, 0, time.Local)
		query = query.Where("createdAt BETWEEN ? AND ?", s, e)
	}
	err := query.Group("DATE(createdAt)").Order("date").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Map the results by date
	dateMap := make(map[string]int)
	for _, r := range results {
		dateMap[r.Date.Format("2006-01-02")] = r.Count
	}

	// Build graph data across requested window (or full history by default)
	graphData := make([]map[string]interface{}, 0)
	if start != nil && end != nil {
		for d := start.Truncate(24 * time.Hour); !d.After(*end); d = d.AddDate(0, 0, 1) {
			ds := d.Format("2006-01-02")
			graphData = append(graphData, map[string]interface{}{
				"date":  ds,
				"count": dateMap[ds],
			})
		}
	} else {
		if len(results) > 0 {
			startDate := results[0].Date
			endDate := results[len(results)-1].Date
			for d := startDate.Truncate(24 * time.Hour); !d.After(endDate); d = d.AddDate(0, 0, 1) {
				ds := d.Format("2006-01-02")
				graphData = append(graphData, map[string]interface{}{
					"date":  ds,
					"count": dateMap[ds],
				})
			}
		}
	}

	data := map[string]interface{}{
		"customer_growth": graphData,
	}

	return data, nil
}

func (r AdminDashboardRepositoryStruct) GetRevenueChart(start *time.Time, end *time.Time) (map[string]interface{}, error) {
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
	query := r.db.Model(&entities.Transaction{}).
		Select("DATE(date) as date, SUM(amount) as count").
		Where("type_in_out = ?", entities.TransactionsTypeInOutIn)
	if start != nil && end != nil {
		query = query.Where("date BETWEEN ? AND ?", start, end)
	}
	err := query.Group("DATE(date)").Order("date").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Map the results by date
	dateMap := make(map[string]int)
	for _, r := range results {
		dateMap[r.Date.Format("2006-01-02")] = r.Count
	}

	// Build graph data across requested window (or full history by default)
	graphData := make([]map[string]interface{}, 0)
	if start != nil && end != nil {
		for d := start.Truncate(24 * time.Hour); !d.After(*end); d = d.AddDate(0, 0, 1) {
			ds := d.Format("2006-01-02")
			graphData = append(graphData, map[string]interface{}{
				"date":   ds,
				"amount": dateMap[ds],
			})
		}
	} else {
		if len(results) > 0 {
			// Normalize to local-day boundaries to avoid timezone off-by-one
			st := results[0].Date.In(time.Local)
			en := results[len(results)-1].Date.In(time.Local)
			startLocal := time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, time.Local)
			endLocal := time.Date(en.Year(), en.Month(), en.Day(), 0, 0, 0, 0, time.Local)
			for d := startLocal; !d.After(endLocal); d = d.AddDate(0, 0, 1) {
				ds := d.Format("2006-01-02")
				graphData = append(graphData, map[string]interface{}{
					"date":   ds,
					"amount": dateMap[ds],
				})
			}
		}
	}

	data := map[string]interface{}{
		"revenue_chart": graphData,
	}

	return data, nil
}

func (r AdminDashboardRepositoryStruct) GetExpensesChart(start *time.Time, end *time.Time) (map[string]interface{}, error) {
	// Check if there are any expense transactions
	var totalTransactions int64
	if err := r.db.Model(&entities.Transaction{}).Where("type_in_out = ?", entities.TransactionsTypeInOutOut).Count(&totalTransactions).Error; err != nil {
		return nil, err
	}
	if totalTransactions == 0 {
		data := map[string]interface{}{
			"expenses_chart": []map[string]interface{}{},
		}
		return data, nil
	}

	var results []DailyCount
	query := r.db.Model(&entities.Transaction{}).
		Select("DATE(date) as date, COALESCE(SUM(amount),0) as count").
		Where("type_in_out = ?", entities.TransactionsTypeInOutOut)
	if start != nil && end != nil {
		query = query.Where("date BETWEEN ? AND ?", start, end)
	}
	err := query.Group("DATE(date)").Order("date").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	dateMap := make(map[string]int)
	for _, r := range results {
		dateMap[r.Date.Format("2006-01-02")] = r.Count
	}

	graphData := make([]map[string]interface{}, 0)
	if start != nil && end != nil {
		for d := start.Truncate(24 * time.Hour); !d.After(*end); d = d.AddDate(0, 0, 1) {
			ds := d.Format("2006-01-02")
			graphData = append(graphData, map[string]interface{}{
				"date":   ds,
				"amount": dateMap[ds],
			})
		}
	} else if len(results) > 0 {
		st := results[0].Date.In(time.Local)
		en := results[len(results)-1].Date.In(time.Local)
		startLocal := time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, time.Local)
		endLocal := time.Date(en.Year(), en.Month(), en.Day(), 0, 0, 0, 0, time.Local)
		for d := startLocal; !d.After(endLocal); d = d.AddDate(0, 0, 1) {
			ds := d.Format("2006-01-02")
			graphData = append(graphData, map[string]interface{}{
				"date":   ds,
				"amount": dateMap[ds],
			})
		}
	}

	data := map[string]interface{}{
		"expenses_chart": graphData,
	}
	return data, nil
}

func (r AdminDashboardRepositoryStruct) GetUnpaidCustomersChart(start *time.Time, end *time.Time) (map[string]interface{}, error) {
	// If no invoices, return empty
	var totalInvoices int64
	if err := r.db.Model(&entities.Invoice{}).Count(&totalInvoices).Error; err != nil {
		return nil, err
	}
	if totalInvoices == 0 {
		return map[string]interface{}{"unpaid_customers_chart": []map[string]interface{}{}}, nil
	}

	type DailyDistinct struct {
		Date  time.Time
		Count int
	}

	// Normalize dates to local timezone before taking DATE() to avoid off-by-one
	offsetSeconds := time.Now().In(time.Local).Second() - time.Now().UTC().Second() // dummy to keep imports; replaced below
	_ = offsetSeconds
	_, tzOffset := time.Now().Zone() // seconds east of UTC
	offsetHours := tzOffset / 3600

	var results []DailyDistinct
	// Use due_date to define the day for unpaid invoices
	selectExpr := fmt.Sprintf("DATE(TIMESTAMPADD(HOUR, %d, due_date)) as date, COUNT(DISTINCT customer_id) as count", offsetHours)
	dateExpr := fmt.Sprintf("DATE(TIMESTAMPADD(HOUR, %d, due_date))", offsetHours)
	query := r.db.Model(&entities.Invoice{}).
		Select(selectExpr).
		Where("status = ?", entities.InvoiceStatusUnpaid)
	if start != nil && end != nil {
		// Compare against due_date window
		s := time.Date(start.In(time.Local).Year(), start.In(time.Local).Month(), start.In(time.Local).Day(), 0, 0, 0, 0, time.Local)
		e := time.Date(end.In(time.Local).Year(), end.In(time.Local).Month(), end.In(time.Local).Day(), 23, 59, 59, 0, time.Local)
		query = query.Where("due_date BETWEEN ? AND ?", s, e)
	}
	if err := query.Group(dateExpr).Order(dateExpr).Scan(&results).Error; err != nil {
		return nil, err
	}

	// Map
	dateMap := make(map[string]int)
	for _, r := range results {
		dateMap[r.Date.Format("2006-01-02")] = r.Count
	}

	graphData := make([]map[string]interface{}, 0)
	if start != nil && end != nil {
		for d := start.Truncate(24 * time.Hour); !d.After(*end); d = d.AddDate(0, 0, 1) {
			ds := d.Format("2006-01-02")
			graphData = append(graphData, map[string]interface{}{
				"date":  ds,
				"count": dateMap[ds],
			})
		}
	} else if len(results) > 0 {
		startDate := results[0].Date
		endDate := results[len(results)-1].Date
		for d := startDate.Truncate(24 * time.Hour); !d.After(endDate); d = d.AddDate(0, 0, 1) {
			ds := d.Format("2006-01-02")
			graphData = append(graphData, map[string]interface{}{
				"date":  ds,
				"count": dateMap[ds],
			})
		}
	}

	return map[string]interface{}{"unpaid_customers_chart": graphData}, nil
}
