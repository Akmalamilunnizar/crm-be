package services

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"skripsi-be/internal/models/entities"

	"gorm.io/gorm"
)

type RouterAction string

const (
	RouterActionSetUnpaidScheduler RouterAction = "set_unpaid_scheduler"
	RouterActionRunOpenScript      RouterAction = "run_open_script"
)

// EnqueueRouterJob creates or upserts a job by unique key (invoiceID:action)
func EnqueueRouterJob(db *gorm.DB, invoiceID string, action RouterAction, delay time.Duration) (*entities.RouterJob, error) {
	if invoiceID == "" || action == "" {
		return nil, errors.New("invalid enqueue arguments")
	}
	unique := invoiceID + ":" + string(action)
	job := entities.RouterJob{}
	if err := db.Where("unique_key = ?", unique).First(&job).Error; err == nil {
		// If already succeeded, do nothing; if pending or error, reset to pending and schedule soon
		job.Status = entities.RouterJobStatusPending
		job.LastError = nil
		job.NextRunAt = time.Now().Add(delay)
		job.UpdatedAt = time.Now()
		if err := db.Save(&job).Error; err != nil {
			return nil, err
		}
		log.Printf("[routerjobs] reset job unique=%s invoice=%s action=%s retry_count=%d next_run_at=%s", unique, invoiceID, string(action), job.RetryCount, job.NextRunAt.Format(time.RFC3339))
		return &job, nil
	}
	// Create new
	job = entities.RouterJob{
		InvoiceID:  invoiceID,
		Action:     string(action),
		UniqueKey:  unique,
		Status:     entities.RouterJobStatusPending,
		RetryCount: 0,
		NextRunAt:  time.Now().Add(delay),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := db.Create(&job).Error; err != nil {
		return nil, err
	}
	log.Printf("[routerjobs] enqueued unique=%s invoice=%s action=%s next_run_at=%s", unique, invoiceID, string(action), job.NextRunAt.Format(time.RFC3339))
	return &job, nil
}

// StartRouterJobWorker runs a background loop to process jobs
func StartRouterJobWorker(db *gorm.DB) {
	log.Printf("[routerjobs] ✅ Router jobs worker started - monitoring for pending MikroTik operations")

	// Start with 2 second polling
	ticker := time.NewTicker(2 * time.Second)
	noJobsCount := 0
	lastStatusLog := time.Now()

	for range ticker.C {
		hasJobs := processDueJobs(db)

		if hasJobs {
			// Reset to fast polling when jobs are found
			noJobsCount = 0
			ticker.Reset(2 * time.Second)
		} else {
			noJobsCount++
			// Gradually increase polling interval when no jobs
			if noJobsCount > 5 { // After 10 seconds of no jobs
				ticker.Reset(10 * time.Second)
			} else if noJobsCount > 2 { // After 4 seconds of no jobs
				ticker.Reset(5 * time.Second)
			}

			// Log status every 5 minutes when idle
			if time.Since(lastStatusLog) > 5*time.Minute {
				log.Printf("[routerjobs] 💤 Worker idle - no pending jobs")
				lastStatusLog = time.Now()
			}
		}
	}
}

func processDueJobs(db *gorm.DB) bool {
	// Fetch one due job to avoid stampede; simple approach
	var job entities.RouterJob
	// Use Find instead of First to avoid "record not found" logs
	result := db.Where("status = ? AND next_run_at <= ?", entities.RouterJobStatusPending, time.Now()).Order("next_run_at ASC").Limit(1).Find(&job)
	if result.Error != nil {
		log.Printf("[routerjobs] query error: %v", result.Error)
		return false
	}
	if result.RowsAffected == 0 {
		// No jobs found - this is normal, don't log
		return false
	}

	// Mark as pending (already), ensure UpdatedAt bumps
	job.UpdatedAt = time.Now()
	_ = db.Save(&job)

	// Load invoice and relations needed
	var invoice entities.Invoice
	if err := db.Preload("Customer").Preload("Customer.Area").First(&invoice, "id = ?", job.InvoiceID).Error; err != nil {
		setJobError(db, &job, fmt.Errorf("invoice not found: %w", err))
		return true // Job was processed (even if failed)
	}

	// Acquire MikroTik connection
	mt := GetSharedMikroTikService()
	if mt == nil || !mt.IsConnected() {
		setJobRetry(db, &job, fmt.Errorf("mikrotik not connected"))
		return true // Job was processed (even if failed)
	}

	var err error
	switch RouterAction(job.Action) {
	case RouterActionSetUnpaidScheduler:
		err = executeSetUnpaidScheduler(mt, &invoice)
	case RouterActionRunOpenScript:
		err = executeRunOpenScript(mt, &invoice)
	default:
		err = fmt.Errorf("unknown action: %s", job.Action)
	}

	if err != nil {
		if isRetryableError(err) {
			setJobRetry(db, &job, err)
			log.Printf("[routerjobs] retry scheduled unique=%s invoice=%s action=%s retry_count=%d err=%v", job.UniqueKey, job.InvoiceID, job.Action, job.RetryCount, err)
		} else {
			setJobError(db, &job, err)
			log.Printf("[routerjobs] error unique=%s invoice=%s action=%s err=%v", job.UniqueKey, job.InvoiceID, job.Action, err)
		}
		return true // Job was processed (even if failed)
	}

	// Success
	job.Status = entities.RouterJobStatusSuccess
	job.LastError = nil
	job.UpdatedAt = time.Now()
	_ = db.Save(&job)
	log.Printf("[routerjobs] success unique=%s invoice=%s action=%s", job.UniqueKey, job.InvoiceID, job.Action)
	return true
}

func isRetryableError(err error) bool {
	msg := strings.ToLower(err.Error())
	// Treat transient connectivity or not-found as retryable for a few attempts
	return strings.Contains(msg, "timeout") || strings.Contains(msg, "temporar") || strings.Contains(msg, "connect")
}

func setJobRetry(db *gorm.DB, job *entities.RouterJob, err error) {
	msg := err.Error()
	job.RetryCount += 1
	// exponential backoff with cap 30m
	backoff := time.Duration(1<<min(job.RetryCount-1, 5)) * time.Minute // 1m,2m,4m,8m,16m,32m
	if backoff > 30*time.Minute {
		backoff = 30 * time.Minute
	}
	job.NextRunAt = time.Now().Add(backoff)
	job.Status = entities.RouterJobStatusPending
	job.LastError = &msg
	job.UpdatedAt = time.Now()
	_ = db.Save(job)
}

func setJobError(db *gorm.DB, job *entities.RouterJob, err error) {
	msg := err.Error()
	job.Status = entities.RouterJobStatusError
	job.LastError = &msg
	job.UpdatedAt = time.Now()
	_ = db.Save(job)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helpers to execute actions
func executeSetUnpaidScheduler(mt *MikroTikService, invoice *entities.Invoice) error {
	if invoice.DueDate == nil {
		return fmt.Errorf("invoice due date is nil")
	}
	name := buildSchedulerName(invoice)
	// Pre-check
	checkCmd := fmt.Sprintf("/system scheduler print count-only where name=\"%s\"", name)
	if out, err := mt.ExecuteCommand(checkCmd); err != nil {
		return fmt.Errorf("mikrotik check failed: %w", err)
	} else if strings.TrimSpace(out) == "0" {
		return fmt.Errorf("scheduler not found: %s", name)
	}
	startDate := invoice.DueDate.Format("Jan/02/2006")
	setCmd := fmt.Sprintf("/system scheduler set [find name=\"%s\"] start-date=\"%s\"", name, startDate)
	// Add log for visibility
	logCmd := fmt.Sprintf(":log info \"crm: set start-date for %s -> %s\"", name, startDate)
	if _, err := mt.ExecuteCommand(setCmd); err != nil {
		return err
	}
	if _, err := mt.ExecuteCommand(logCmd); err != nil {
		// non-fatal
		log.Printf("[routerjobs] log command error: %v", err)
	}
	return nil
}

func executeRunOpenScript(mt *MikroTikService, invoice *entities.Invoice) error {
	name := buildSchedulerName(invoice)
	scriptName := "open_" + name
	// Pre-check
	checkCmd := fmt.Sprintf("/system script print count-only where name=\"%s\"", scriptName)
	if out, err := mt.ExecuteCommand(checkCmd); err != nil {
		return fmt.Errorf("mikrotik check failed: %w", err)
	} else if strings.TrimSpace(out) == "0" {
		return fmt.Errorf("script not found: %s", scriptName)
	}
	runCmd := fmt.Sprintf("/system script run \"%s\"", scriptName)
	logCmd := fmt.Sprintf(":log info \"crm: run %s\"", scriptName)
	if _, err := mt.ExecuteCommand(runCmd); err != nil {
		return err
	}
	if _, err := mt.ExecuteCommand(logCmd); err != nil {
		log.Printf("[routerjobs] log command error: %v", err)
	}
	return nil
}

func buildSchedulerName(invoice *entities.Invoice) string {
	code := ""
	custName := ""
	if invoice.Customer.Name != "" {
		custName = invoice.Customer.Name
	}
	if invoice.Customer.Area.CodeName != "" {
		code = invoice.Customer.Area.CodeName
	}
	if code == "" {
		code = "N/A"
	}
	if custName == "" {
		custName = "Unknown"
	}
	return fmt.Sprintf("%s - %s", code, custName)
}
