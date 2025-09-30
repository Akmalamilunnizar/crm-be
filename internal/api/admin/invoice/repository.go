package invoice

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"skripsi-be/internal/models/entities"
	"skripsi-be/internal/services"
)

type AdminInvoiceRepositoryInterface interface {
	FindAdminInvoiceRepository() ([]entities.Invoice, error)
	CreateAdminInvoiceRepository(request CreateAdminInvoiceRequest) (entities.Invoice, error)
	FindByIdAdminInvoiceRepository(request IdAdminInvoiceRequest) (entities.Invoice, error)
	UpdateAdminInvoiceRepository(request UpdateAdminInvoiceRequest) (entities.Invoice, error)
	UpdateStatusAdminInvoiceRepository(request UpdateStatusAdminInvoiceRequest) (entities.Invoice, error)
	DeleteAdminInvoiceRepository(request IdAdminInvoiceRequest) (entities.Invoice, error)
	ProcessPartialPaymentRepository(request PartialPaymentRequest) (entities.Invoice, error)
	MarkPdfViewedRepository(request IdAdminInvoiceRequest) (entities.Invoice, error)
	FindDueUnpaidInvoices() ([]entities.Invoice, error)
	FindPaidInvoices() ([]entities.Invoice, error)
	FindNewestInvoicePerCustomer() ([]entities.Invoice, error)
	PrintAllUnpaidInvoicesRepository() (map[string]interface{}, error)
}

type AdminInvoiceRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminInvoiceRepository(db *gorm.DB) *AdminInvoiceRepositoryStruct {
	return &AdminInvoiceRepositoryStruct{db}
}

func (r AdminInvoiceRepositoryStruct) FindAdminInvoiceRepository() ([]entities.Invoice, error) {
	invoices := []entities.Invoice{}
	tx := r.db.Preload("Customer").Preload("InvoiceItems.Invoice").Order("createdAt DESC").Find(&invoices)

	if tx.Error != nil {
		return invoices, tx.Error
	}

	// For each invoice, calculate total paid from transactions
	for i := range invoices {
		var totalPaid int64
		r.db.Model(&entities.Transaction{}).Where("invoice_id = ?", invoices[i].ID).Select("COALESCE(SUM(amount), 0)").Scan(&totalPaid)

		// Create a virtual transaction object to hold total paid
		if totalPaid > 0 {
			invoices[i].Transaction = entities.Transaction{
				Amount: totalPaid,
			}
		}
	}

	return invoices, nil
}

func (r AdminInvoiceRepositoryStruct) FindByIdAdminInvoiceRepository(request IdAdminInvoiceRequest) (entities.Invoice, error) {
	invoice := entities.Invoice{}
	tx := r.db.Preload("Customer").Preload("InvoiceItems.Invoice").Preload("Transaction").First(&invoice, "id = ?", request.Id)

	if tx.Error != nil {
		return invoice, tx.Error
	}

	// Calculate total amount from invoice items if not set or incorrect
	if len(invoice.InvoiceItems) > 0 {
		var totalAmount int64 = 0
		for _, item := range invoice.InvoiceItems {
			totalAmount += item.Total
		}

		// Update invoice amount if it's different from calculated total
		if invoice.Amount != totalAmount {
			invoice.Amount = totalAmount
			// Update in database
			r.db.Model(&invoice).Update("amount", totalAmount)
		}
	}

	// Calculate total paid from all transactions for this invoice
	var totalPaid int64
	r.db.Model(&entities.Transaction{}).Where("invoice_id = ?", request.Id).Select("COALESCE(SUM(amount), 0)").Scan(&totalPaid)

	// Create a virtual transaction object to hold total paid
	if totalPaid > 0 {
		invoice.Transaction = entities.Transaction{
			Amount: totalPaid,
		}
	}

	return invoice, nil
}

func (r AdminInvoiceRepositoryStruct) CreateAdminInvoiceRepository(request CreateAdminInvoiceRequest) (entities.Invoice, error) {
	// Create invoice entity manually to ensure proper field mapping
	// Determine status from request; default unpaid
	status := entities.InvoiceStatusUnpaid
	if request.Status != nil {
		if *request.Status == string(entities.InvoiceStatusPaid) {
			status = entities.InvoiceStatusPaid
		} else if *request.Status == string(entities.InvoiceStatusPending) {
			status = entities.InvoiceStatusPending
		}
	}

	invoice := entities.Invoice{
		CustomerID: request.CustomerID,
		Amount:     request.Amount,
		Status:     status,
	}

	// Set invoice date to today and derive due date based on the first item's quantity interpreted as number of months
	// If qty is 0 or items are empty, default to 1 month
	computeClampedMonthlyDate := func(base time.Time, addMonths int, preferredDay int) time.Time {
		// Move to the first day of the target month, then clamp the day
		year, month, _ := base.Date()
		// Calculate target year/month with carry
		targetMonthIndex := int(month) - 1 + addMonths
		targetYear := year + targetMonthIndex/12
		targetMonth := time.Month(targetMonthIndex%12 + 1)
		// Last day of target month
		firstOfTarget := time.Date(targetYear, targetMonth, 1, base.Hour(), base.Minute(), base.Second(), base.Nanosecond(), base.Location())
		firstOfNext := firstOfTarget.AddDate(0, 1, 0)
		lastDay := firstOfNext.AddDate(0, 0, -1).Day()
		day := preferredDay
		if day > lastDay {
			day = lastDay
		}
		return time.Date(targetYear, targetMonth, day, base.Hour(), base.Minute(), base.Second(), base.Nanosecond(), base.Location())
	}

	now := time.Now()
	invoiceDate := now
	invoice.InvoiceDate = &invoiceDate

	if len(request.InvoiceItems) > 0 {
		months := int(request.InvoiceItems[0].Qty)
		if months <= 0 {
			months = 1
		}
		// Use invoice date as the base, keep the same day number when adding months
		preferredDay := now.Day()
		due := computeClampedMonthlyDate(now, months, preferredDay)
		invoice.DueDate = &due
	}

	// Convert request invoice items to entity invoice items
	var invoiceItems []entities.InvoiceItems
	for _, reqItem := range request.InvoiceItems {
		item := entities.InvoiceItems{
			// Don't set ID here - let it be generated fresh in the loop below
			Name:  reqItem.Name,
			Price: reqItem.Price,
			Qty:   reqItem.Qty,
			Total: reqItem.Total,
		}
		invoiceItems = append(invoiceItems, item)
	}
	invoice.InvoiceItems = invoiceItems

	tx := r.db.Begin()

	// Create the invoice first
	txInvoice := tx.Create(&invoice)
	if txInvoice.Error != nil {
		tx.Rollback()
		return entities.Invoice{}, txInvoice.Error
	}

	// Calculate total amount from all invoice items and save them
	var totalAmount int64 = 0
	for i := range invoice.InvoiceItems {
		// Generate a fresh UUID for each item to ensure uniqueness
		newID := uuid.New().String()
		log.Printf("Creating invoice item %d with ID: %s", i, newID)

		invoice.InvoiceItems[i].ID = newID
		invoice.InvoiceItems[i].InvoiceID = invoice.ID
		invoice.InvoiceItems[i].Total = invoice.InvoiceItems[i].Price * invoice.InvoiceItems[i].Qty
		totalAmount += invoice.InvoiceItems[i].Total

		// Save each invoice item individually
		txInvoiceItem := tx.Create(&invoice.InvoiceItems[i])
		if txInvoiceItem.Error != nil {
			log.Printf("Error creating invoice item %d: %v", i, txInvoiceItem.Error)
			tx.Rollback()
			return entities.Invoice{}, txInvoiceItem.Error
		}
		log.Printf("Successfully created invoice item %d with ID: %s", i, newID)
	}

	// Set total amount to sum of all items
	invoice.Amount = totalAmount
	txInvoice = tx.Model(&invoice).Updates(entities.Invoice{Amount: invoice.Amount})
	if txInvoice.Error != nil {
		tx.Rollback()
		return entities.Invoice{}, txInvoice.Error
	}

	tx.Commit()

	// If pending with reason, store it
	if invoice.Status == entities.InvoiceStatusPending && request.PendingReason != nil && *request.PendingReason != "" {
		pr := entities.InvoicePendingReason{ID: uuid.New().String(), InvoiceID: invoice.ID, Reason: *request.PendingReason}
		if err := r.db.Create(&pr).Error; err != nil {
			log.Printf("failed to store pending reason: %v", err)
		}
	}

	// Trigger MikroTik based on status (best-effort)
	go func(inv entities.Invoice) {
		defer func() { _ = recover() }()
		svc := services.GetSharedMikroTikService()
		if svc == nil || !svc.IsConnected() {
			log.Printf("[mikrotik] skip: not connected (invoice=%s status=%s)", inv.ID, string(inv.Status))
			return
		}
		// Fetch customer's network devices to get MACs
		var devices []entities.NetworkDevice
		if err := r.db.Where("customer_id = ? AND mac_address IS NOT NULL", inv.CustomerID).Find(&devices).Error; err != nil {
			log.Printf("[mikrotik] fetch devices error invoice=%s: %v", inv.ID, err)
			return
		}
		log.Printf("[mikrotik] applying status=%s to %d device(s) for customer=%s invoice=%s", string(inv.Status), len(devices), inv.CustomerID, inv.ID)
		for _, d := range devices {
			mac := strings.TrimSpace(*d.MacAddress)
			if mac == "" {
				continue
			}
			switch inv.Status {
			case entities.InvoiceStatusPaid:
				// bypass binding and mark script comment sudahbayar
				cmd1 := fmt.Sprintf("/ip hotspot ip-binding set [find mac-address=%s] type=byp", mac)
				log.Printf("[mikrotik] exec: %s", cmd1)
				if out, err := svc.ExecuteCommand(cmd1); err != nil {
					log.Printf("[mikrotik] error: %v", err)
				} else {
					log.Printf("[mikrotik] out: %s", strings.TrimSpace(out))
				}
				// optional: set script comment if script exists
			case entities.InvoiceStatusUnpaid:
				// disable binding and mark belumbayar
				cmd2 := fmt.Sprintf("/ip hotspot ip-binding set [find mac-address=%s] disabled=yes", mac)
				log.Printf("[mikrotik] exec: %s", cmd2)
				if out, err := svc.ExecuteCommand(cmd2); err != nil {
					log.Printf("[mikrotik] error: %v", err)
				} else {
					log.Printf("[mikrotik] out: %s", strings.TrimSpace(out))
				}
			case entities.InvoiceStatusPending:
				// keep bypassed; no change
				log.Printf("[mikrotik] pending: no change for mac=%s", mac)
			}
		}
	}(invoice)

	// Fetch the complete invoice with all relationships
	var completeInvoice entities.Invoice
	err := r.db.Preload("Customer").Preload("InvoiceItems").First(&completeInvoice, "id = ?", invoice.ID)
	if err.Error != nil {
		return entities.Invoice{}, err.Error
	}

	// Synchronous preflight for unpaid scheduler existence so errors can bubble to API
	svc := services.GetSharedMikroTikService()
	if completeInvoice.Status == entities.InvoiceStatusUnpaid && completeInvoice.DueDate != nil && svc != nil && svc.IsConnected() {
		name := r.generateSchedulerName(completeInvoice.CustomerID, completeInvoice.Customer.Name)
		checkCmd := fmt.Sprintf("/system scheduler print count-only where name=\"%s\"", name)
		if out, e := svc.ExecuteCommand(checkCmd); e == nil {
			if strings.TrimSpace(out) == "0" {
				msg := fmt.Sprintf("MikroTik scheduler not found: name=%s for invoice=%s", name, completeInvoice.ID)
				log.Printf("[mikrotik] %s", msg)
				_, _ = svc.ExecuteCommand(fmt.Sprintf(":log error \"CRM: %s\"", strings.ReplaceAll(msg, "\"", "'")))
				return entities.Invoice{}, fmt.Errorf(msg)
			}
		}
	}

	// Create a Winbox (MikroTik) scheduler once per invoice based on status
	go func(inv entities.Invoice) {
		defer func() { _ = recover() }()
		svc := services.GetSharedMikroTikService()
		if svc == nil || !svc.IsConnected() {
			log.Printf("[mikrotik] scheduler: skip (not connected) invoice=%s status=%s", inv.ID, string(inv.Status))
			return
		}

		// Get first MAC for this customer
		var dev []entities.NetworkDevice
		if err := r.db.Where("customer_id = ? AND mac_address IS NOT NULL AND mac_address <> ''", inv.CustomerID).Order("created_at ASC").Limit(1).Find(&dev).Error; err != nil {
			log.Printf("[mikrotik] scheduler: devices fetch error invoice=%s: %v", inv.ID, err)
			return
		}
		if len(dev) == 0 || dev[0].MacAddress == nil || strings.TrimSpace(*dev[0].MacAddress) == "" {
			log.Printf("[mikrotik] scheduler: no MAC found for customer=%s invoice=%s", inv.CustomerID, inv.ID)
			return
		}
		mac := strings.TrimSpace(*dev[0].MacAddress)

		// Generate scheduler name in format "Area - Customer Name"
		schedulerName := r.generateSchedulerName(inv.CustomerID, inv.Customer.Name)
		customerName := inv.Customer.Name
		if customerName == "" {
			customerName = inv.CustomerID
		}
		// Match script name with scheduler naming to keep things consistent on MikroTik
		// Result: open_<CodeName - Customer Name>
		scriptName := "open_" + schedulerName

		if inv.Status == entities.InvoiceStatusUnpaid {
			// For unpaid invoices, update the scheduler with due_date using the new naming format
			if inv.DueDate != nil {
				// Format the due date for MikroTik scheduler (e.g., "oct/28/2025")
				dueDateFormatted := inv.DueDate.Format("Jan/02/2006")
				// Convert to lowercase for MikroTik format
				dueDateFormatted = strings.ToLower(dueDateFormatted)
				// Pre-check: ensure the scheduler exists; if not, log error both in Go and RouterOS
				checkCmd := fmt.Sprintf("/system scheduler print count-only where name=\"%s\"", schedulerName)
				outCheck, errCheck := svc.ExecuteCommand(checkCmd)
				if errCheck != nil {
					log.Printf("[mikrotik] scheduler check error invoice=%s: %v", inv.ID, errCheck)
					// attempt to emit router log about the failure context
					_, _ = svc.ExecuteCommand(fmt.Sprintf(":log error \"CRM: scheduler check error name='%s' invoice=%s\"", schedulerName, inv.ID))
					return
				}
				if strings.TrimSpace(outCheck) == "0" {
					msg := fmt.Sprintf("MikroTik scheduler not found: name=%s for invoice=%s", schedulerName, inv.ID)
					log.Printf("[mikrotik] %s", msg)
					_, _ = svc.ExecuteCommand(fmt.Sprintf(":log error \"CRM: %s\"", strings.ReplaceAll(msg, "\"", "'")))
					// Surface this as hard error by updating DB status -> error and attaching message
					_ = r.db.Model(&entities.Invoice{}).Where("id = ?", inv.ID).Updates(map[string]interface{}{
						"status": string(entities.InvoiceStatusPending),
						"link":   msg,
					}).Error
					// Also return immediately to avoid executing the set
					return
				}
				// Append a RouterOS log so it is visible in Winbox logs
				cmd := fmt.Sprintf("/system scheduler set [find name=\"%s\"] start-date=%s; :log info \"CRM: scheduler set name='%s' start-date=%s for invoice=%s\"", schedulerName, dueDateFormatted, schedulerName, dueDateFormatted, inv.ID)
				log.Printf("[mikrotik] scheduler: set %s start-date=%s for unpaid invoice=%s", schedulerName, dueDateFormatted, inv.ID)
				if out, err := svc.ExecuteCommand(cmd); err != nil {
					log.Printf("[mikrotik] scheduler set error invoice=%s: %v", inv.ID, err)
				} else {
					log.Printf("[mikrotik] scheduler set out: %s", strings.TrimSpace(out))
				}
			} else {
				log.Printf("[mikrotik] scheduler: no due date for unpaid invoice=%s", inv.ID)
			}
		} else {
			// For paid and pending invoices, use the original scheduler logic with new naming format
			var onEvent string
			if inv.Status == entities.InvoiceStatusPaid {
				onEvent = fmt.Sprintf("ip/hotspot/ip-binding/set type=byp [find mac-address=%s]; sys/script/set comment=sudahbayar [find name=\"%s\"];", mac, scriptName)
			} else {
				// treat pending the same as unpaid for this scheduler request
				onEvent = fmt.Sprintf("ip/hotspot/ip-binding/set disabled=y [find mac-address=%s]; sys/script/set comment=belumbayar [find name=\"%s\"];", mac, scriptName)
			}

			// Append a RouterOS log so it is visible in Winbox logs
			cmd := fmt.Sprintf("/system scheduler add name=\"%s\" on-event=\"%s\" start-time=startup; :log info \"CRM: scheduler add name='%s' status=%s\"", schedulerName, onEvent, schedulerName, string(inv.Status))
			log.Printf("[mikrotik] scheduler: create name=%s invoice=%s status=%s mac=%s", schedulerName, inv.ID, string(inv.Status), mac)
			if out, err := svc.ExecuteCommand(cmd); err != nil {
				log.Printf("[mikrotik] scheduler create error invoice=%s: %v", inv.ID, err)
			} else {
				log.Printf("[mikrotik] scheduler create out: %s", strings.TrimSpace(out))
			}
		}
	}(completeInvoice)

	return completeInvoice, nil
}

// generateSchedulerName creates a scheduler name in format "CodeName - Customer Name"
func (r AdminInvoiceRepositoryStruct) generateSchedulerName(customerID, customerName string) string {
	// Fetch customer with area information
	var customer entities.Customer
	if err := r.db.Preload("Area").Where("id = ?", customerID).First(&customer).Error; err != nil {
		log.Printf("[scheduler] error fetching customer area for customer=%s: %v", customerID, err)
		// Fallback to customer name only if area fetch fails
		if customerName == "" {
			return "Unknown - Unknown"
		}
		return "Unknown - " + customerName
	}

	// Get area code name (use code_name as the area identifier)
	areaCode := "Unknown"
	if customer.Area != nil && customer.Area.CodeName != "" {
		areaCode = customer.Area.CodeName
	}

	// Use customer name; if missing, show Unknown (avoid exposing internal IDs)
	if customerName == "" {
		customerName = "Unknown"
	}

	// Create scheduler name in format "CodeName - Customer Name"
	schedulerName := fmt.Sprintf("%s - %s", areaCode, customerName)

	// Log the generated scheduler name for debugging
	log.Printf("[scheduler] generated scheduler name: %s for customer=%s", schedulerName, customerID)

	return schedulerName
}

func (r AdminInvoiceRepositoryStruct) UpdateAdminInvoiceRepository(request UpdateAdminInvoiceRequest) (entities.Invoice, error) {
	invoice := entities.Invoice{}
	tx := r.db.First(&invoice, "id = ?", request.Id)

	if tx.Error != nil {
		return invoice, tx.Error
	}
	copier.Copy(&invoice, &request)

	r.db.Save(&invoice)
	return invoice, nil
}

func (r AdminInvoiceRepositoryStruct) UpdateStatusAdminInvoiceRepository(request UpdateStatusAdminInvoiceRequest) (entities.Invoice, error) {
	invoice := entities.Invoice{}
	tx := r.db.First(&invoice, "id = ?", request.Id)

	if tx.Error != nil {
		return invoice, tx.Error
	}

	// Clear previous error link and update status
	tx = r.db.Model(&invoice).Updates(map[string]interface{}{
		"status": request.Status,
		"link":   "",
	})
	if tx.Error != nil {
		return invoice, tx.Error
	}

	// Reload the invoice to get updated data
	r.db.Preload("Customer").Preload("InvoiceItems.Invoice").Preload("Transaction").First(&invoice, "id = ?", request.Id)

	// Trigger MikroTik enforcement based on status (synchronous for unpaid to bubble errors)
	if request.Status == "paid" {
		go r.enforceMikroTikForPaidInvoice(invoice)
		// For PAID, still ensure scheduler start-date gets updated to the invoice due_date
		_, _ = services.EnqueueRouterJob(r.db, invoice.ID, services.RouterActionSetUnpaidScheduler, 0)
		// Also run customer-specific open_ script
		_, _ = services.EnqueueRouterJob(r.db, invoice.ID, services.RouterActionRunOpenScript, 0)
	} else if request.Status == "unpaid" {
		// Enqueue unpaid scheduler update instead of synchronous call
		_, _ = services.EnqueueRouterJob(r.db, invoice.ID, services.RouterActionSetUnpaidScheduler, 0)
		// Reload to capture any link/message written during scheduler checks
		r.db.Preload("Customer").First(&invoice, "id = ?", invoice.ID)
		if invoice.Link != "" && (strings.Contains(invoice.Link, "MikroTik") || strings.Contains(strings.ToLower(invoice.Link), "scheduler")) {
			return invoice, fmt.Errorf(invoice.Link)
		}
	} else if request.Status == "pending" {
		// For PENDING, keep scheduler aligned with due_date and run open script
		_, _ = services.EnqueueRouterJob(r.db, invoice.ID, services.RouterActionSetUnpaidScheduler, 0)
		_, _ = services.EnqueueRouterJob(r.db, invoice.ID, services.RouterActionRunOpenScript, 0)
	}

	return invoice, nil
}

// runOpenScriptForInvoice executes RouterOS script named "open_<Code - Customer>"
func (r AdminInvoiceRepositoryStruct) runOpenScriptForInvoice(invoice entities.Invoice) {
	defer func() { _ = recover() }()
	mt := services.GetSharedMikroTikService()
	if mt == nil || !mt.IsConnected() {
		log.Printf("[open-script] MikroTik not connected; skipping invoice %s", invoice.ID)
		return
	}
	// Ensure customer is loaded
	if invoice.Customer.ID == "" {
		r.db.Preload("Customer").First(&invoice, "id = ?", invoice.ID)
	}
	// Build script name from code_name - customer name
	schedulerName := r.generateSchedulerName(invoice.CustomerID, invoice.Customer.Name)
	scriptName := "open_" + schedulerName
	// Verify the script exists (exact name)
	check := fmt.Sprintf("/system script print count-only where name=\"%s\"", scriptName)
	out, err := mt.ExecuteCommand(check)
	if err != nil {
		log.Printf("[open-script] check failed invoice=%s: %v", invoice.ID, err)
		return
	}
	if strings.TrimSpace(out) == "0" {
		// Fallback: try pattern search by customer name regardless of code_name
		pattern := fmt.Sprintf("open_.* - %s", invoice.Customer.Name)
		check2 := fmt.Sprintf("/system script print count-only where name~\"%s\"", pattern)
		if out2, err2 := mt.ExecuteCommand(check2); err2 == nil && strings.TrimSpace(out2) != "0" {
			// Run by pattern find
			cmd := fmt.Sprintf("/system script run [find name~\"%s\"]; :log info \"CRM: run script by pattern '%s' for invoice=%s\"", pattern, pattern, invoice.ID)
			log.Printf("[open-script] run by pattern %s for invoice %s", pattern, invoice.ID)
			if out, err := mt.ExecuteCommand(cmd); err != nil {
				log.Printf("[open-script] error running pattern script for invoice %s: %v", invoice.ID, err)
				_ = r.db.Model(&entities.Invoice{}).Where("id = ?", invoice.ID).Update("link", err.Error()).Error
			} else {
				log.Printf("[open-script] run out: %s", strings.TrimSpace(out))
			}
			return
		}
		msg := fmt.Sprintf("MikroTik script not found: name=%s for invoice=%s", scriptName, invoice.ID)
		log.Printf("[open-script] %s", msg)
		_, _ = mt.ExecuteCommand(fmt.Sprintf(":log error \"CRM: %s\"", strings.ReplaceAll(msg, "\"", "'")))
		// Persist message to surface to API/UI
		_ = r.db.Model(&entities.Invoice{}).Where("id = ?", invoice.ID).Update("link", msg).Error
		return
	}
	// Run the script (exact)
	cmd := fmt.Sprintf("/system script run \"%s\"; :log info \"CRM: run script '%s' for invoice=%s\"", scriptName, scriptName, invoice.ID)
	log.Printf("[open-script] run %s for invoice %s", scriptName, invoice.ID)
	if out, err := mt.ExecuteCommand(cmd); err != nil {
		log.Printf("[open-script] error running script for invoice %s: %v", invoice.ID, err)
		_ = r.db.Model(&entities.Invoice{}).Where("id = ?", invoice.ID).Update("link", err.Error()).Error
	} else {
		log.Printf("[open-script] run out: %s", strings.TrimSpace(out))
	}
}

func (r AdminInvoiceRepositoryStruct) DeleteAdminInvoiceRepository(request IdAdminInvoiceRequest) (entities.Invoice, error) {

	invoice := entities.Invoice{}

	tx := r.db.First(&invoice, "id = ?", request.Id)
	if tx.Error != nil {
		return invoice, tx.Error
	}

	tx = r.db.Delete(&invoice)
	if tx.Error != nil {
		return invoice, tx.Error
	}
	return invoice, nil
}

// FindDueUnpaidInvoices returns unpaid invoices with due_date <= NOW().
// If your schema does not include due_date, adjust to your real condition.
func (r AdminInvoiceRepositoryStruct) FindDueUnpaidInvoices() ([]entities.Invoice, error) {
	var invoices []entities.Invoice
	tx := r.db.Where("status = ? AND due_date <= NOW()", entities.InvoiceStatusUnpaid).Find(&invoices)
	return invoices, tx.Error
}

// FindPaidInvoices returns the newest paid invoice for each customer
func (r AdminInvoiceRepositoryStruct) FindPaidInvoices() ([]entities.Invoice, error) {
	var invoices []entities.Invoice

	// Get the latest paid invoice for each customer
	subQuery := r.db.Table("invoices").
		Select("customer_id, MAX(createdAt) as max_created_at").
		Where("status = ?", entities.InvoiceStatusPaid).
		Group("customer_id")

	tx := r.db.Table("invoices i").
		Select("i.*").
		Joins("INNER JOIN (?) s ON i.customer_id = s.customer_id AND i.createdAt = s.max_created_at", subQuery).
		Where("i.status = ?", entities.InvoiceStatusPaid).
		Find(&invoices)

	return invoices, tx.Error
}

// FindNewestInvoicePerCustomer returns the newest invoice for each customer
func (r AdminInvoiceRepositoryStruct) FindNewestInvoicePerCustomer() ([]entities.Invoice, error) {
	var invoices []entities.Invoice

	// Get the latest invoice for each customer (regardless of status)
	subQuery := r.db.Table("invoices").
		Select("customer_id, MAX(createdAt) as max_created_at").
		Group("customer_id")

	tx := r.db.Table("invoices i").
		Select("i.*").
		Joins("INNER JOIN (?) s ON i.customer_id = s.customer_id AND i.createdAt = s.max_created_at", subQuery).
		Find(&invoices)

	return invoices, tx.Error
}

func (r AdminInvoiceRepositoryStruct) ProcessPartialPaymentRepository(request PartialPaymentRequest) (entities.Invoice, error) {
	// Start transaction
	tx := r.db.Begin()

	// Get invoice (no customer preload to avoid selecting legacy customer.product_id column)
	invoice := entities.Invoice{}
	err := tx.Preload("InvoiceItems").First(&invoice, "id = ?", request.Id).Error
	if err != nil {
		tx.Rollback()
		return invoice, err
	}

	// Validate invoice status (should be pending or unpaid)
	if invoice.Status == entities.InvoiceStatusPaid {
		tx.Rollback()
		return invoice, fmt.Errorf("invoice is already fully paid")
	}

	// Calculate current total paid from existing transactions
	var currentTotalPaid int64
	tx.Model(&entities.Transaction{}).Where("invoice_id = ?", request.Id).Select("COALESCE(SUM(amount), 0)").Scan(&currentTotalPaid)

	// Calculate new total paid
	newTotalPaid := currentTotalPaid + request.Amount

	// Validate payment amount doesn't exceed invoice amount
	if newTotalPaid > invoice.Amount {
		tx.Rollback()
		return invoice, fmt.Errorf("payment amount exceeds outstanding balance")
	}

	// Create payment transaction
	transaction := entities.Transaction{
		AccountID:   "b82074e7-7acb-40e2-ab33-014f4b09c1f8", // Default account ID
		InvoiceID:   request.Id,
		TypeCash:    entities.TransactionsTypeCashInternet,
		TypeInOut:   entities.TransactionsTypeInOutIn,
		Description: fmt.Sprintf("Partial payment for invoice %s", request.Id),
		Amount:      request.Amount,
		Category:    "partial_payment",
		Method:      "manual",
	}

	// If a reason/note was provided, append it to description
	if request.Reason != nil && *request.Reason != "" {
		transaction.Description = fmt.Sprintf("%s | Reason: %s", transaction.Description, *request.Reason)
	}

	err = tx.Create(&transaction).Error
	if err != nil {
		tx.Rollback()
		return invoice, err
	}

	// Update invoice status based on payment
	var newStatus entities.InvoiceStatus
	if newTotalPaid >= invoice.Amount {
		newStatus = entities.InvoiceStatusPaid
	} else {
		newStatus = entities.InvoiceStatusPending
	}

	err = tx.Model(&invoice).Update("status", newStatus).Error
	if err != nil {
		tx.Rollback()
		return invoice, err
	}

	// Commit transaction
	tx.Commit()

	// Reload invoice with updated data
	r.db.Preload("Customer").Preload("InvoiceItems").Preload("Transaction").First(&invoice, "id = ?", request.Id)

	return invoice, nil
}

func (r AdminInvoiceRepositoryStruct) MarkPdfViewedRepository(request IdAdminInvoiceRequest) (entities.Invoice, error) {
	var invoice entities.Invoice

	// Check if invoice exists
	err := r.db.Preload("Customer").Preload("InvoiceItems").Preload("Transaction").First(&invoice, "id = ?", request.Id).Error
	if err != nil {
		return invoice, err
	}

	// Check if PDF has already been viewed
	if invoice.PdfViewed {
		return invoice, fmt.Errorf("PDF has already been viewed and cannot be accessed again")
	}

	// Mark PDF as viewed with current timestamp
	now := time.Now()
	err = r.db.Model(&invoice).Updates(map[string]interface{}{
		"pdf_viewed":    true,
		"pdf_viewed_at": &now,
	}).Error

	if err != nil {
		return invoice, err
	}

	// Reload invoice with updated data
	r.db.Preload("Customer").Preload("InvoiceItems").Preload("Transaction").First(&invoice, "id = ?", request.Id)

	return invoice, nil
}

// setTLSKSchedulerForUnpaidInvoice sets the scheduler with due_date for unpaid invoices using Area - Customer Name format
func (r AdminInvoiceRepositoryStruct) setTLSKSchedulerForUnpaidInvoice(invoice entities.Invoice) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("[tlsk-scheduler] panic: %v", rec)
		}
	}()

	mt := services.GetSharedMikroTikService()
	if mt == nil || !mt.IsConnected() {
		log.Printf("[tlsk-scheduler] MikroTik not connected; skipping invoice %s", invoice.ID)
		return
	}

	if invoice.DueDate == nil {
		log.Printf("[tlsk-scheduler] no due date for unpaid invoice %s", invoice.ID)
		return
	}

	// Generate scheduler name in format "Area - Customer Name"
	schedulerName := r.generateSchedulerName(invoice.CustomerID, invoice.Customer.Name)

	// Pre-check existence of the target scheduler; RouterOS may accept set with empty find silently
	checkCmd := fmt.Sprintf("/system scheduler print count-only where name=\"%s\"", schedulerName)
	if out, err := mt.ExecuteCommand(checkCmd); err != nil {
		log.Printf("[tlsk-scheduler] scheduler check error for invoice %s: %v", invoice.ID, err)
		return
	} else if strings.TrimSpace(out) == "0" {
		msg := fmt.Sprintf("MikroTik scheduler not found: name=%s for invoice=%s", schedulerName, invoice.ID)
		log.Printf("[tlsk-scheduler] %s", msg)
		_, _ = mt.ExecuteCommand(fmt.Sprintf(":log error \"CRM: %s\"", strings.ReplaceAll(msg, "\"", "'")))
		// Persist message so the synchronous caller can return an error and the UI can show a toast
		_ = r.db.Model(&entities.Invoice{}).Where("id = ?", invoice.ID).Updates(map[string]interface{}{
			"link":   msg,
			"status": string(entities.InvoiceStatusPending),
		}).Error
		return
	}

	// Format the due date for MikroTik scheduler (e.g., "oct/28/2025")
	dueDateFormatted := invoice.DueDate.Format("Jan/02/2006")
	// Convert to lowercase for MikroTik format
	dueDateFormatted = strings.ToLower(dueDateFormatted)
	cmd := fmt.Sprintf("/system scheduler set [find name=\"%s\"] start-date=%s; :log info \"CRM: scheduler set name='%s' start-date=%s for invoice=%s\"", schedulerName, dueDateFormatted, schedulerName, dueDateFormatted, invoice.ID)

	log.Printf("[tlsk-scheduler] setting %s start-date=%s for unpaid invoice %s", schedulerName, dueDateFormatted, invoice.ID)
	if out, err := mt.ExecuteCommand(cmd); err != nil {
		log.Printf("[tlsk-scheduler] error setting scheduler for invoice %s: %v", invoice.ID, err)
	} else {
		log.Printf("[tlsk-scheduler] scheduler set successfully for invoice %s: %s", invoice.ID, strings.TrimSpace(out))
	}
}

// enforceMikroTikForPaidInvoice sets IP binding type to "bypassed" for paid invoices
func (r AdminInvoiceRepositoryStruct) enforceMikroTikForPaidInvoice(invoice entities.Invoice) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("[enforce-paid] panic: %v", rec)
		}
	}()

	mt := services.GetSharedMikroTikService()
	if mt == nil || !mt.IsConnected() {
		log.Printf("[enforce-paid] MikroTik not connected; skipping invoice %s", invoice.ID)
		return
	}

	log.Printf("[enforce-paid] processing paid invoice %s for customer %s", invoice.ID, invoice.CustomerID)

	// Get customer's network devices
	var devices []struct{ MacAddress *string }
	if err := r.db.Table("network_devices").
		Select("mac_address").
		Where("customer_id = ? AND mac_address IS NOT NULL AND mac_address <> ''", invoice.CustomerID).
		Scan(&devices).Error; err != nil {
		log.Printf("[enforce-paid] devices fetch error customer=%s: %v", invoice.CustomerID, err)
		return
	}

	changed := 0
	for _, d := range devices {
		if d.MacAddress == nil {
			continue
		}
		mac := *d.MacAddress

		if err := mt.SetHotspotIPBindingType(mac, "bypassed"); err != nil {
			log.Printf("[enforce-paid] set bypassed failed mac=%s: %v", mac, err)
		} else {
			changed++
			log.Printf("[enforce-paid] set bypassed for mac=%s", mac)
		}
	}

	if changed > 0 {
		log.Printf("[enforce-paid] set type=bypassed for %d device(s) (invoice %s)", changed, invoice.ID)
	} else {
		log.Printf("[enforce-paid] no devices found for customer %s (invoice %s)", invoice.CustomerID, invoice.ID)
	}
}

// PrintAllUnpaidInvoicesRepository generates thermal printer data for all unpaid invoices
func (r AdminInvoiceRepositoryStruct) PrintAllUnpaidInvoicesRepository() (map[string]interface{}, error) {
	var invoices []entities.Invoice

	// Get all unpaid invoices with customer and invoice items
	err := r.db.Preload("Customer").Preload("InvoiceItems").Where("status = ?", entities.InvoiceStatusUnpaid).Find(&invoices).Error
	if err != nil {
		return nil, err
	}

	// Generate thermal printer format
	thermalData := generateThermalPrinterData(invoices)

	result := map[string]interface{}{
		"total_invoices": len(invoices),
		"thermal_data":   thermalData,
		"print_ready":    true,
	}

	return result, nil
}

// generateThermalPrinterData creates thermal printer formatted data
func generateThermalPrinterData(invoices []entities.Invoice) string {
	var output strings.Builder

	// Compact header printed once
	output.WriteString(centerText("PT JR Nusa Menara Networks", 30) + "\n")
	output.WriteString(centerText("Jl Raya Talangsuko 373 Turen", 30) + "\n")
	output.WriteString(centerText("PH: 08123511147, (0341)8224357", 30) + "\n")
	output.WriteString(centerText("EMAIL: info@menara.net.id", 30) + "\n")
	output.WriteString(centerText("LINK: www.menara.net.id", 30) + "\n")
	output.WriteString(strings.Repeat("=", 30) + "\n")

	for i, invoice := range invoices {
		output.WriteString(centerText("*** Tanda Terima ***", 30) + "\n")
		output.WriteString("----- DITERBITKAN UNTUK -----\n")
		output.WriteString(fmt.Sprintf("%s\n", invoice.Customer.Name))
		output.WriteString("------------------------------\n")

		// Aligned key-values
		output.WriteString(formatKV("Nomor", invoice.ID, 30) + "\n")
		output.WriteString(formatKV("Tanggal", invoice.CreatedAt.Format("02-01-2006"), 30) + "\n")

		output.WriteString("------------------------------\n")
		output.WriteString(centerText("DETAIL", 30) + "\n")

		// Compact item list
		for _, item := range invoice.InvoiceItems {
			output.WriteString("- " + item.Name + "\n")
		}

		output.WriteString("------------------------------\n")
		output.WriteString(formatKV("Total", "Rp "+formatCurrency(invoice.Amount), 30) + "\n")
		output.WriteString(formatKV("Dibayar", "Rp 0", 30) + "\n")
		output.WriteString(formatKV("Saldo", "Rp "+formatCurrency(invoice.Amount), 30) + "\n")

		// Keep a single blank line between invoices, none at the end
		if i < len(invoices)-1 {
			output.WriteString("\n")
		}
	}

	// Trim trailing whitespace/newlines to reduce extra blank page risk
	return strings.TrimRight(output.String(), "\n\r\t ")
}

// formatCurrency formats currency in Indonesian Rupiah format
func formatCurrency(amount int64) string {
	return fmt.Sprintf("%.0f", float64(amount))
}

// padRight pads s with spaces to the right up to width

// padLeft pads s with spaces to the left up to width
func padLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}

func centerText(text string, width int) string {
	l := len(text)
	if l >= width || width <= 0 {
		return text
	}
	pad := (width - l) / 2
	if pad < 0 {
		pad = 0
	}
	return strings.Repeat(" ", pad) + text
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func formatKV(label, value string, width int) string {
	maxLabel := width - len(value) - 1
	if maxLabel < 1 {
		maxLabel = 1
	}
	if len(label) > maxLabel {
		label = label[:maxLabel]
	}
	return padRight(label, maxLabel) + " " + value
}
