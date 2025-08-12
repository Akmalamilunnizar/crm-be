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
	err := r.db.
		Model(&entities.Customer{}).
		Joins("JOIN products ON products.id = customer.product_id").
		Select("COUNT(*) as count, products.name").
		Group("products.id").
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
