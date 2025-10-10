package services

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"skripsi-be/internal/models/entities"

	"gorm.io/gorm"
)

// MikrotikProvisioningService handles automated provisioning of customer installations on RouterOS
type MikrotikProvisioningService struct {
	db           *gorm.DB
	mikrotikConn *MikroTikService // Use the actual MikroTik service
}

// NewMikrotikProvisioningService creates a new provisioning service
func NewMikrotikProvisioningService(db *gorm.DB, mikrotikConn *MikroTikService) *MikrotikProvisioningService {
	return &MikrotikProvisioningService{
		db:           db,
		mikrotikConn: mikrotikConn,
	}
}

// getProductBandwidthInfo gets bandwidth information from customer's network device product
func (s *MikrotikProvisioningService) getProductBandwidthInfo(customerID, macAddress string) (maxLimit, comment string, err error) {
	var networkDevice entities.NetworkDevice
	var product entities.Products

	log.Printf("🔍 DEBUG: Looking for network device with customer_id=%s, mac_address=%s", customerID, macAddress)

	// Find the network device by customer ID and MAC address
	err = s.db.Preload("Product").Where("customer_id = ? AND mac_address = ?", customerID, macAddress).First(&networkDevice).Error
	if err != nil {
		log.Printf("❌ Network device not found for customer %s, MAC %s: %v", customerID, macAddress, err)
		// Return default values if not found
		return "10M/10M", "Default Package", nil
	}

	log.Printf("✅ Found network device: ID=%s, ProductID=%v", networkDevice.ID, networkDevice.ProductID)

	// Check if the network device has a product
	if networkDevice.Product == nil {
		log.Printf("❌ No product found for network device %s (ProductID: %v)", networkDevice.ID, networkDevice.ProductID)
		return "10M/10M", "No Product Package", nil
	}

	log.Printf("✅ Found product: ID=%s, Name=%s", networkDevice.Product.ID, networkDevice.Product.Name)

	product = *networkDevice.Product

	// Handle nullable bandwidth fields - use defaults if not set
	downloadMbps := 10 // Default 10 Mbps
	uploadMbps := 10   // Default 10 Mbps

	if product.DownloadSpeedMbps != nil {
		downloadMbps = *product.DownloadSpeedMbps
	}
	if product.UploadSpeedMbps != nil {
		uploadMbps = *product.UploadSpeedMbps
	}

	// Format as MikroTik expects (e.g., "100M/100M")
	maxLimit = fmt.Sprintf("%dM/%dM", downloadMbps, uploadMbps)

	// Use product description as comment, otherwise create one from product name
	if product.Description != "" {
		comment = product.Description
	} else {
		comment = fmt.Sprintf("%s - %dMbps/%dMbps", product.Name, downloadMbps, uploadMbps)
	}

	log.Printf("Product bandwidth info: %s, comment: %s", maxLimit, comment)
	return maxLimit, comment, nil
}

// ProvisioningRequest contains all data needed for provisioning
type ProvisioningRequest struct {
	InstallationID string
	CustomerID     string
	CustomerName   string
	AreaName       string // Area name for code generation
	MACAddress     string
	IPAddress      string // IP address from the form
	StartDate      string // YYYY-MM-DD
	StartTime      string // HH:MM:SS
	MaxLimit       string // e.g., "10M/10M"
	Comment        string // e.g., "150rb/10Mbps"
	DryRun         bool
	CreatedBy      string
}

// ProvisioningResult contains the result of provisioning operation
type ProvisioningResult struct {
	Success          bool     `json:"success"`
	InstallationID   string   `json:"installation_id"`
	CodeName         string   `json:"code_name"`
	MACAddress       string   `json:"mac_address"`
	IPAddress        string   `json:"ip_address,omitempty"`
	Commands         []string `json:"commands"`
	CommandsOutput   []string `json:"commands_output,omitempty"`
	Errors           []string `json:"errors,omitempty"`
	DryRun           bool     `json:"dry_run"`
	ExecutionTimeMs  int      `json:"execution_time_ms"`
	ResourcesCreated []string `json:"resources_created,omitempty"`
	ResourcesUpdated []string `json:"resources_updated,omitempty"`
	LogID            string   `json:"log_id"`
}

// ProvisionInstallation provisions a customer installation on MikroTik RouterOS
func (s *MikrotikProvisioningService) ProvisionInstallation(req ProvisioningRequest) (*ProvisioningResult, error) {
	return s.ProvisionInstallationWithCount(req, 0) // 0 means calculate count automatically
}

// ProvisionInstallationWithCount provisions with a pre-calculated installation count
func (s *MikrotikProvisioningService) ProvisionInstallationWithCount(req ProvisioningRequest, preCalculatedCount int) (*ProvisioningResult, error) {
	startTime := time.Now()

	// Validate request
	if err := s.validateProvisioningRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate unique CODE-NAME with pre-calculated count if provided
	var codeName string
	if preCalculatedCount > 0 {
		// Use pre-calculated count for R-number
		codeName = s.generateCodeNameWithCount(req, preCalculatedCount)
		log.Printf("🏷️ Generated code name with pre-calculated count R%d: %s", preCalculatedCount, codeName)
	} else {
		// Use automatic count calculation (legacy behavior)
		codeName = s.generateCodeName(req)
		log.Printf("🏷️ Generated code name with automatic count: %s", codeName)
	}

	// Initialize result
	result := &ProvisioningResult{
		InstallationID:   req.InstallationID,
		CodeName:         codeName,
		MACAddress:       req.MACAddress,
		Commands:         []string{},
		CommandsOutput:   []string{},
		Errors:           []string{},
		DryRun:           req.DryRun,
		ResourcesCreated: []string{},
		ResourcesUpdated: []string{},
	}

	// Create provisioning log entry
	// Handle CreatedBy - set to nil if empty to avoid foreign key constraint issues
	var createdBy *string
	if req.CreatedBy != "" {
		createdBy = &req.CreatedBy
	} else {
		createdBy = nil
	}

	provLog := &entities.InstallationProvisioningLog{
		CustomerInstallationID: req.InstallationID,
		CustomerID:             &req.CustomerID,
		MACAddress:             &req.MACAddress,
		CodeName:               &codeName,
		Status:                 entities.ProvisioningStatusQueued,
		ProvisioningType:       entities.ProvisioningTypeNew,
		DryRun:                 req.DryRun,
		CreatedBy:              createdBy,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	if req.DryRun {
		provLog.ProvisioningType = entities.ProvisioningTypeDryRun
	}

	if err := s.db.Create(provLog).Error; err != nil {
		return nil, fmt.Errorf("failed to create provisioning log: %w", err)
	}
	result.LogID = provLog.ID

	// Update status to running
	provLog.Status = entities.ProvisioningStatusRunning
	s.db.Save(provLog)

	// Execute provisioning sequence
	if err := s.executeProvisioningSequence(req, codeName, result); err != nil {
		// Mark as failed
		provLog.Status = entities.ProvisioningStatusFailed
		errMsg := err.Error()
		provLog.ErrorMessage = &errMsg
		result.Success = false
		result.Errors = append(result.Errors, err.Error())
	} else {
		// Mark as success
		provLog.Status = entities.ProvisioningStatusSuccess
		result.Success = true
	}

	// Calculate execution time
	executionTime := int(time.Since(startTime).Milliseconds())
	result.ExecutionTimeMs = executionTime
	provLog.ExecutionTimeMs = &executionTime

	// Update IP address if found
	if result.IPAddress != "" {
		provLog.IPAddress = &result.IPAddress
	}

	// Store commands and output as JSON
	commandsJSON, _ := json.Marshal(result.Commands)
	outputJSON, _ := json.Marshal(result.CommandsOutput)
	provLog.CommandsExecuted = entities.JSONArray{}
	provLog.CommandsOutput = entities.JSONArray{}
	json.Unmarshal(commandsJSON, &provLog.CommandsExecuted)
	json.Unmarshal(outputJSON, &provLog.CommandsOutput)

	// Save final log state
	s.db.Save(provLog)

	// Update installation record
	if !req.DryRun && result.Success {
		s.updateInstallationStatus(req.InstallationID, codeName, result.IPAddress)
	}

	return result, nil
}

// executeProvisioningSequence executes the RouterOS command sequence
func (s *MikrotikProvisioningService) executeProvisioningSequence(req ProvisioningRequest, codeName string, result *ProvisioningResult) error {
	// Step 1: Find DHCP lease and get IP address
	ipAddress, err := s.findDHCPLeaseIP(req.MACAddress, req.IPAddress, result)
	if err != nil {
		return fmt.Errorf("failed to find DHCP lease: %w", err)
	}
	result.IPAddress = ipAddress

	// Step 2: Make DHCP lease static
	if err := s.makeLeaseStatic(ipAddress, result, req.DryRun); err != nil {
		return fmt.Errorf("failed to make lease static: %w", err)
	}

	// Step 3: Get product bandwidth information
	maxLimit, comment, err := s.getProductBandwidthInfo(req.CustomerID, req.MACAddress)
	if err != nil {
		log.Printf("Warning: Failed to get product bandwidth info, using request values: %v", err)
		maxLimit = req.MaxLimit
		comment = req.Comment
	}

	// Step 4: Check and create/update queue
	if err := s.createOrUpdateQueue(codeName, ipAddress, maxLimit, comment, result, req.DryRun); err != nil {
		return fmt.Errorf("failed to create/update queue: %w", err)
	}

	// Step 5: Check and create/update hotspot IP binding
	if err := s.createOrUpdateHotspotBinding(codeName, req.MACAddress, ipAddress, result, req.DryRun); err != nil {
		return fmt.Errorf("failed to create/update hotspot binding: %w", err)
	}

	// Step 6: Check and create/update netwatch
	if err := s.createOrUpdateNetwatch(codeName, req.MACAddress, ipAddress, result, req.DryRun); err != nil {
		return fmt.Errorf("failed to create/update netwatch: %w", err)
	}

	// Step 7: Check and create/update scheduler
	if err := s.createOrUpdateScheduler(codeName, req.MACAddress, req.StartDate, req.StartTime, result, req.DryRun); err != nil {
		return fmt.Errorf("failed to create/update scheduler: %w", err)
	}

	// Step 8: Check and create/update script
	if err := s.createOrUpdateScript(codeName, req.MACAddress, result, req.DryRun); err != nil {
		return fmt.Errorf("failed to create/update script: %w", err)
	}

	return nil
}

// findDHCPLeaseIP finds the IP address from DHCP lease for the given MAC
func (s *MikrotikProvisioningService) findDHCPLeaseIP(macAddress string, providedIP string, result *ProvisioningResult) (string, error) {
	// If IP address is provided from the form, use it directly
	if providedIP != "" {
		cmd := fmt.Sprintf("/ip/dhcp-server/lease/print where mac-address=%s status=bound disabled=no", macAddress)
		result.Commands = append(result.Commands, cmd)
		result.CommandsOutput = append(result.CommandsOutput, "Using provided IP address: "+providedIP)
		return providedIP, nil
	}

	// Fallback: Try to find IP from DHCP lease
	cmd := fmt.Sprintf("/ip/dhcp-server/lease/print where mac-address=%s status=bound disabled=no", macAddress)
	result.Commands = append(result.Commands, cmd)

	if s.mikrotikConn == nil {
		// Dry run or no connection - use default IP
		result.CommandsOutput = append(result.CommandsOutput, "DRY RUN: Would find IP for MAC "+macAddress)
		return "192.168.1.100", nil // Default IP for dry run
	}

	// Execute actual command to find DHCP lease
	output, err := s.mikrotikConn.ExecuteCommand(cmd)
	if err != nil {
		return "", err
	}

	result.CommandsOutput = append(result.CommandsOutput, output)

	// Parse the output to extract the actual IP address
	// MikroTik output format: "address=10.10.21.125 mac-address=40:EE:15:7D:43:01 ..."
	ipAddress := s.parseDHCPLeaseIP(output)
	if ipAddress != "" {
		return ipAddress, nil
	}

	// Fallback to default IP if parsing fails
	return "192.168.1.100", nil
}

// parseDHCPLeaseIP parses MikroTik output to extract IP address
func (s *MikrotikProvisioningService) parseDHCPLeaseIP(output string) string {
	// MikroTik output format examples:
	// "address=10.10.21.125 mac-address=40:EE:15:7D:43:01 client-id=1:40:ee:15:7d:43:1 server=dhcp1 status=bound"
	// or with multiple lines, we need to find the line with the MAC address

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// Look for address= pattern in each line
		if strings.Contains(line, "address=") {
			// Extract IP address using regex
			re := regexp.MustCompile(`address=([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				return matches[1]
			}
		}
	}

	return ""
}

// makeLeaseStatic makes the DHCP lease static
func (s *MikrotikProvisioningService) makeLeaseStatic(ipAddress string, result *ProvisioningResult, dryRun bool) error {
	cmd := fmt.Sprintf("/ip/dhcp-server/lease/make-static [find address=%s]", ipAddress)
	result.Commands = append(result.Commands, cmd)

	if dryRun || s.mikrotikConn == nil {
		result.CommandsOutput = append(result.CommandsOutput, "DRY RUN: Would make lease static for "+ipAddress)
		return nil
	}

	output, err := s.mikrotikConn.ExecuteCommand(cmd)
	result.CommandsOutput = append(result.CommandsOutput, output)
	if err != nil {
		return err
	}

	result.ResourcesCreated = append(result.ResourcesCreated, "DHCP Static Lease: "+ipAddress)
	return nil
}

// CleanupInstallation removes all Mikrotik configurations related to an installation
func (s *MikrotikProvisioningService) CleanupInstallation(installation entities.CustomerInstallation, dryRun bool) error {
	if s.mikrotikConn == nil || !s.mikrotikConn.IsConnected() {
		log.Printf("⚠️ Mikrotik not connected, skipping cleanup for installation %s", installation.ID)
		return nil
	}

	log.Printf("🧹 Starting Mikrotik cleanup for installation %s (Code: %s)", installation.ID, *installation.CodeName)

	// Get all network devices for this installation
	var networkDevices []entities.NetworkDevice
	if err := s.db.Where("customer_installation_id = ?", installation.ID).Find(&networkDevices).Error; err != nil {
		log.Printf("❌ Failed to get network devices for cleanup: %v", err)
		return err
	}

	// Clean up configurations for each network device
	for _, device := range networkDevices {
		if device.MacAddress != nil && *device.MacAddress != "" {
			if err := s.cleanupDeviceConfigurations(*device.MacAddress, installation.CodeName, dryRun); err != nil {
				log.Printf("❌ Failed to cleanup device %s: %v", *device.MacAddress, err)
				// Continue with other devices even if one fails
			}
		}
	}

	log.Printf("✅ Mikrotik cleanup completed for installation %s", installation.ID)
	return nil
}

// cleanupDeviceConfigurations removes all Mikrotik configurations for a specific device
func (s *MikrotikProvisioningService) cleanupDeviceConfigurations(macAddress string, codeName *string, dryRun bool) error {
	if macAddress == "" {
		return nil
	}

	codeNameStr := ""
	if codeName != nil {
		codeNameStr = *codeName
	}

	// List of commands to remove all configurations
	cleanupCommands := []string{
		// Remove queue simple rules
		fmt.Sprintf("/queue/simple/remove [find name=\"%s\"]", codeNameStr),

		// Remove hotspot IP bindings
		fmt.Sprintf("/ip/hotspot/ip-binding/remove [find mac-address=%s]", macAddress),

		// Remove netwatch entries
		fmt.Sprintf("/tool/netwatch/remove [find comment=\"%s\"]", codeNameStr),

		// Remove schedulers
		fmt.Sprintf("/system/scheduler/remove [find name=\"%s\"]", codeNameStr),

		// Remove scripts
		fmt.Sprintf("/system/script/remove [find name=\"open_%s\"]", codeNameStr),

		// Remove DHCP leases
		fmt.Sprintf("/ip/dhcp-server/lease/remove [find mac-address=%s]", macAddress),
	}

	log.Printf("🗑️ Executing %d cleanup commands for MAC %s, Code %s", len(cleanupCommands), macAddress, codeNameStr)

	// Execute each cleanup command
	for _, cmd := range cleanupCommands {
		if dryRun {
			log.Printf("DRY RUN: Would execute: %s", cmd)
			continue
		}

		output, err := s.mikrotikConn.ExecuteCommand(cmd)
		if err != nil {
			log.Printf("⚠️ Warning: Failed to execute cleanup command '%s': %v", cmd, err)
			// Continue with other commands even if one fails
		} else {
			log.Printf("✅ Cleanup command executed: %s", cmd)
			if output != "" {
				log.Printf("   Output: %s", strings.TrimSpace(output))
			}
		}
	}

	return nil
}

// createOrUpdateQueue creates or updates queue simple rule
func (s *MikrotikProvisioningService) createOrUpdateQueue(codeName, ipAddress, maxLimit, comment string, result *ProvisioningResult, dryRun bool) error {
	// Try to create new queue - MikroTik will handle duplicates
	// If queue exists, it will return an error which we can handle gracefully
	existingQueue := false

	if existingQueue {
		// Update existing queue
		cmd := fmt.Sprintf("/queue/simple/set [find name=%s] target=%s max-limit=%s comment=\"%s\"",
			codeName, ipAddress, maxLimit, s.sanitizeComment(comment))
		result.Commands = append(result.Commands, cmd)

		if dryRun || s.mikrotikConn == nil {
			result.CommandsOutput = append(result.CommandsOutput, "DRY RUN: Would update queue "+codeName)
			result.ResourcesUpdated = append(result.ResourcesUpdated, "Queue: "+codeName)
			return nil
		}

		// Execute the queue update command
		output, err := s.mikrotikConn.ExecuteCommand(cmd)
		result.CommandsOutput = append(result.CommandsOutput, output)
		result.ResourcesUpdated = append(result.ResourcesUpdated, "Queue: "+codeName)
		return err
	}

	// Create new queue - use MikroTik's duplicate handling
	// First try to remove existing queue with same name, then add new one
	removeCmd := fmt.Sprintf("/queue/simple/remove [find name=\"%s\"]", codeName)
	addCmd := fmt.Sprintf("/queue/simple/add name=\"%s\" target=%s max-limit=%s comment=\"%s\"",
		codeName, ipAddress, maxLimit, s.sanitizeComment(comment))

	result.Commands = append(result.Commands, removeCmd)
	result.Commands = append(result.Commands, addCmd)

	if dryRun || s.mikrotikConn == nil {
		result.CommandsOutput = append(result.CommandsOutput, "DRY RUN: Would create queue "+codeName)
		result.ResourcesCreated = append(result.ResourcesCreated, "Queue: "+codeName)
		return nil
	}

	// Execute the queue removal command (ignore errors if queue doesn't exist)
	output1, _ := s.mikrotikConn.ExecuteCommand(removeCmd)
	result.CommandsOutput = append(result.CommandsOutput, output1)

	// Execute the queue creation command
	output2, err := s.mikrotikConn.ExecuteCommand(addCmd)
	result.CommandsOutput = append(result.CommandsOutput, output2)
	result.ResourcesCreated = append(result.ResourcesCreated, "Queue: "+codeName)
	return err
}

// createOrUpdateHotspotBinding creates or updates hotspot IP binding
func (s *MikrotikProvisioningService) createOrUpdateHotspotBinding(codeName, macAddress, ipAddress string, result *ProvisioningResult, dryRun bool) error {
	// Check if binding exists (simplified - always create new for now)
	// In a real implementation, you'd check if the binding exists first
	existingBinding := false

	if existingBinding {
		// Update existing binding
		cmd := fmt.Sprintf("/ip/hotspot/ip-binding/set [find mac-address=%s] address=%s to-addr=%s type=byp comment=\"%s\"",
			macAddress, ipAddress, ipAddress, codeName)
		result.Commands = append(result.Commands, cmd)

		if dryRun || s.mikrotikConn == nil {
			result.CommandsOutput = append(result.CommandsOutput, "DRY RUN: Would update hotspot binding for "+macAddress)
			result.ResourcesUpdated = append(result.ResourcesUpdated, "Hotspot Binding: "+macAddress)
			return nil
		}

		// Execute the hotspot binding update command
		output, err := s.mikrotikConn.ExecuteCommand(cmd)
		result.CommandsOutput = append(result.CommandsOutput, output)
		result.ResourcesUpdated = append(result.ResourcesUpdated, "Hotspot Binding: "+macAddress)
		return err
	}

	// Create new binding - use MikroTik's duplicate handling
	// First try to remove existing binding with same MAC, then add new one
	removeCmd := fmt.Sprintf("/ip/hotspot/ip-binding/remove [find mac-address=%s]", macAddress)
	addCmd := fmt.Sprintf("/ip/hotspot/ip-binding/add mac-address=%s address=%s to-addr=%s type=byp comment=\"%s\"",
		macAddress, ipAddress, ipAddress, codeName)

	result.Commands = append(result.Commands, removeCmd)
	result.Commands = append(result.Commands, addCmd)

	if dryRun || s.mikrotikConn == nil {
		result.CommandsOutput = append(result.CommandsOutput, "DRY RUN: Would create hotspot binding for "+macAddress)
		result.ResourcesCreated = append(result.ResourcesCreated, "Hotspot Binding: "+macAddress)
		return nil
	}

	// Execute the hotspot binding removal command (ignore errors if binding doesn't exist)
	output1, _ := s.mikrotikConn.ExecuteCommand(removeCmd)
	result.CommandsOutput = append(result.CommandsOutput, output1)

	// Execute the hotspot binding creation command
	output2, err := s.mikrotikConn.ExecuteCommand(addCmd)
	result.CommandsOutput = append(result.CommandsOutput, output2)
	result.ResourcesCreated = append(result.ResourcesCreated, "Hotspot Binding: "+macAddress)
	return err
}

// createOrUpdateNetwatch creates or updates netwatch entry
func (s *MikrotikProvisioningService) createOrUpdateNetwatch(codeName, macAddress, ipAddress string, result *ProvisioningResult, dryRun bool) error {
	downScript := fmt.Sprintf("/ip hotspot/host remove [find mac-address=%s];", macAddress)
	testScript := fmt.Sprintf("/ip hotspot/host remove [find mac-address=%s];", macAddress)
	upScript := ""

	// Check if netwatch exists (simplified - always create new for now)
	// In a real implementation, you'd check if the netwatch exists first
	existingNetwatch := false

	if existingNetwatch {
		// Update existing netwatch
		cmd := fmt.Sprintf("/tool/netwatch/set [find comment=%s] host=%s down-script=\"%s\" test-script=\"%s\" up-script=\"%s\"",
			codeName, ipAddress, downScript, testScript, upScript)
		result.Commands = append(result.Commands, cmd)

		if dryRun || s.mikrotikConn == nil {
			result.CommandsOutput = append(result.CommandsOutput, "DRY RUN: Would update netwatch "+codeName)
			result.ResourcesUpdated = append(result.ResourcesUpdated, "Netwatch: "+codeName)
			return nil
		}

		// Execute the netwatch update command
		output, err := s.mikrotikConn.ExecuteCommand(cmd)
		result.CommandsOutput = append(result.CommandsOutput, output)
		result.ResourcesUpdated = append(result.ResourcesUpdated, "Netwatch: "+codeName)
		return err
	}

	// Create new netwatch - use MikroTik's duplicate handling
	// First try to remove existing netwatch with same comment, then add new one
	removeCmd := fmt.Sprintf("/tool/netwatch/remove [find comment=\"%s\"]", codeName)
	addCmd := fmt.Sprintf("/tool/netwatch/add comment=\"%s\" host=%s disabled=no type=icmp down-script=\"%s\" test-script=\"%s\" up-script=\"%s\"",
		codeName, ipAddress, downScript, testScript, upScript)

	result.Commands = append(result.Commands, removeCmd)
	result.Commands = append(result.Commands, addCmd)

	if dryRun || s.mikrotikConn == nil {
		result.CommandsOutput = append(result.CommandsOutput, "DRY RUN: Would create netwatch "+codeName)
		result.ResourcesCreated = append(result.ResourcesCreated, "Netwatch: "+codeName)
		return nil
	}

	// Execute the netwatch removal command (ignore errors if netwatch doesn't exist)
	output1, _ := s.mikrotikConn.ExecuteCommand(removeCmd)
	result.CommandsOutput = append(result.CommandsOutput, output1)

	// Execute the netwatch creation command
	output2, err := s.mikrotikConn.ExecuteCommand(addCmd)
	result.CommandsOutput = append(result.CommandsOutput, output2)
	result.ResourcesCreated = append(result.ResourcesCreated, "Netwatch: "+codeName)
	return err
}

// createOrUpdateScheduler creates or updates scheduler entry
func (s *MikrotikProvisioningService) createOrUpdateScheduler(codeName, macAddress, startDate, startTime string, result *ProvisioningResult, dryRun bool) error {
	// Simplify the onEvent string to avoid syntax errors - use bypassed for installation reports
	onEvent := fmt.Sprintf("ip/hotspot/ip-binding/set type=byp [find mac-address=%s]", macAddress)

	// Check if scheduler exists (simplified - always create new for now)
	// In a real implementation, you'd check if the scheduler exists first
	existingScheduler := false

	if existingScheduler {
		// Update existing scheduler
		cmd := fmt.Sprintf("/system/scheduler/set [find name=%s] interval=4w2d start-date=%s start-time=%s on-event=\"%s\"",
			codeName, startDate, startTime, onEvent)
		result.Commands = append(result.Commands, cmd)

		if dryRun || s.mikrotikConn == nil {
			result.CommandsOutput = append(result.CommandsOutput, "DRY RUN: Would update scheduler "+codeName)
			result.ResourcesUpdated = append(result.ResourcesUpdated, "Scheduler: "+codeName)
			return nil
		}

		// Execute the scheduler update command
		output, err := s.mikrotikConn.ExecuteCommand(cmd)
		result.CommandsOutput = append(result.CommandsOutput, output)
		result.ResourcesUpdated = append(result.ResourcesUpdated, "Scheduler: "+codeName)
		return err
	}

	// Create new scheduler - use MikroTik's duplicate handling
	// First try to remove existing scheduler with same name, then add new one
	removeCmd := fmt.Sprintf("/system/scheduler/remove [find name=\"%s\"]", codeName)
	addCmd := fmt.Sprintf("/system/scheduler/add comment=isolir interval=4w2d name=\"%s\" on-event=\"%s\" start-date=%s start-time=%s",
		codeName, onEvent, startDate, startTime)

	result.Commands = append(result.Commands, removeCmd)
	result.Commands = append(result.Commands, addCmd)

	if dryRun || s.mikrotikConn == nil {
		result.CommandsOutput = append(result.CommandsOutput, "DRY RUN: Would create scheduler "+codeName)
		result.ResourcesCreated = append(result.ResourcesCreated, "Scheduler: "+codeName)
		return nil
	}

	// Execute the scheduler removal command (ignore errors if scheduler doesn't exist)
	output1, _ := s.mikrotikConn.ExecuteCommand(removeCmd)
	result.CommandsOutput = append(result.CommandsOutput, output1)

	// Execute the scheduler creation command
	output2, err := s.mikrotikConn.ExecuteCommand(addCmd)
	result.CommandsOutput = append(result.CommandsOutput, output2)
	result.ResourcesCreated = append(result.ResourcesCreated, "Scheduler: "+codeName)
	return err
}

// createOrUpdateScript creates or updates system script
func (s *MikrotikProvisioningService) createOrUpdateScript(codeName, macAddress string, result *ProvisioningResult, dryRun bool) error {
	scriptName := "open_" + codeName
	// Simplify the source string to avoid syntax errors
	source := fmt.Sprintf("ip/hotspot/ip-binding/set type=byp [find mac-address=%s]", macAddress)

	// Check if script exists (simplified - always create new for now)
	// In a real implementation, you'd check if the script exists first
	existingScript := false

	if existingScript {
		// Update existing script
		cmd := fmt.Sprintf("/system/script/set [find name=%s] source=\"%s\" comment=belumbayar",
			scriptName, source)
		result.Commands = append(result.Commands, cmd)

		if dryRun || s.mikrotikConn == nil {
			result.CommandsOutput = append(result.CommandsOutput, "DRY RUN: Would update script "+scriptName)
			result.ResourcesUpdated = append(result.ResourcesUpdated, "Script: "+scriptName)
			return nil
		}

		// Execute the script update command
		output, err := s.mikrotikConn.ExecuteCommand(cmd)
		result.CommandsOutput = append(result.CommandsOutput, output)
		result.ResourcesUpdated = append(result.ResourcesUpdated, "Script: "+scriptName)
		return err
	}

	// Create new script - use MikroTik's duplicate handling
	// First try to remove existing script with same name, then add new one
	removeCmd := fmt.Sprintf("/system/script/remove [find name=\"%s\"]", scriptName)
	addCmd := fmt.Sprintf("/system/script/add comment=belumbayar dont-require-permissions=no name=\"%s\" source=\"%s\"",
		scriptName, source)

	result.Commands = append(result.Commands, removeCmd)
	result.Commands = append(result.Commands, addCmd)

	if dryRun || s.mikrotikConn == nil {
		result.CommandsOutput = append(result.CommandsOutput, "DRY RUN: Would create script "+scriptName)
		result.ResourcesCreated = append(result.ResourcesCreated, "Script: "+scriptName)
		return nil
	}

	// Execute the script removal command (ignore errors if script doesn't exist)
	output1, _ := s.mikrotikConn.ExecuteCommand(removeCmd)
	result.CommandsOutput = append(result.CommandsOutput, output1)

	// Execute the script creation command
	output2, err := s.mikrotikConn.ExecuteCommand(addCmd)
	result.CommandsOutput = append(result.CommandsOutput, output2)
	result.ResourcesCreated = append(result.ResourcesCreated, "Script: "+scriptName)
	return err
}

// generateCodeName creates a unique CODE-NAME for RouterOS resources
// Format: [area_code]_EA[number]_[sales_rep_code]_[customer_name]_R[number]
func (s *MikrotikProvisioningService) generateCodeName(req ProvisioningRequest) string {
	// Get the installation count for this customer (increments for each new installation)
	installationCount := s.getInstallationCount(req.CustomerID, req.AreaName)

	// Get sales representative information
	salesRepCode, _ := s.getSalesRepresentativeInfo(req.CustomerID)

	// Generate EA code (EA0001, EA0002, etc.)
	eaCode := s.generateEACode()

	// Sanitize area name (remove special chars, keep only alphanumeric)
	areaReg := regexp.MustCompile("[^a-zA-Z0-9]+")
	sanitizedArea := areaReg.ReplaceAllString(req.AreaName, "")
	sanitizedArea = strings.TrimSpace(sanitizedArea)

	// Sanitize customer name (remove special chars, keep only alphanumeric and spaces)
	customerReg := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	sanitizedCustomer := customerReg.ReplaceAllString(req.CustomerName, "")
	sanitizedCustomer = strings.TrimSpace(sanitizedCustomer)

	// Sanitize sales rep code (remove special chars, keep only alphanumeric)
	salesRepReg := regexp.MustCompile("[^a-zA-Z0-9]+")
	sanitizedSalesRep := salesRepReg.ReplaceAllString(salesRepCode, "")
	sanitizedSalesRep = strings.TrimSpace(sanitizedSalesRep)

	// Format: [area_code]_EA[number]_[sales_rep_code]_[customer_name]_R[number]
	return fmt.Sprintf("%s_%s_%s_%s_R%d", sanitizedArea, eaCode, sanitizedSalesRep, sanitizedCustomer, installationCount)
}

// generateCodeNameWithCount creates a unique CODE-NAME with a pre-calculated installation count
// Format: [area_code]_EA[number]_[sales_rep_code]_[customer_name]_R[number]
func (s *MikrotikProvisioningService) generateCodeNameWithCount(req ProvisioningRequest, installationCount int) string {
	// Get sales representative information
	salesRepCode, _ := s.getSalesRepresentativeInfo(req.CustomerID)

	// Generate EA code (EA0001, EA0002, etc.)
	eaCode := s.generateEACode()

	// Sanitize area name (remove special chars, keep only alphanumeric)
	areaReg := regexp.MustCompile("[^a-zA-Z0-9]+")
	sanitizedArea := areaReg.ReplaceAllString(req.AreaName, "")
	sanitizedArea = strings.TrimSpace(sanitizedArea)

	// Sanitize customer name (remove special chars, keep only alphanumeric and spaces)
	customerReg := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	sanitizedCustomer := customerReg.ReplaceAllString(req.CustomerName, "")
	sanitizedCustomer = strings.TrimSpace(sanitizedCustomer)

	// Sanitize sales rep code (remove special chars, keep only alphanumeric)
	salesRepReg := regexp.MustCompile("[^a-zA-Z0-9]+")
	sanitizedSalesRep := salesRepReg.ReplaceAllString(salesRepCode, "")
	sanitizedSalesRep = strings.TrimSpace(sanitizedSalesRep)

	// Format: [area_code]_EA[number]_[sales_rep_code]_[customer_name]_R[number]
	return fmt.Sprintf("%s_%s_%s_%s_R%d", sanitizedArea, eaCode, sanitizedSalesRep, sanitizedCustomer, installationCount)
}

// getInstallationCount returns the number of installations for a customer
func (s *MikrotikProvisioningService) getInstallationCount(customerID, areaName string) int {
	var count int64

	// Count ALL installations for this customer (regardless of area)
	// This ensures R1, R2, R3... increment correctly for each new installation
	err := s.db.Table("customer_installations").
		Where("customer_installations.customer_id = ?", customerID).
		Count(&count).Error

	if err != nil {
		log.Printf("Error counting installations: %v", err)
		// Return 1 as default if there's an error
		return 1
	}

	log.Printf("🔍 Installation count for customer %s: found %d existing installations", customerID, count)
	log.Printf("🔍 Will return R%d for new installation", count+1)

	// Debug: Let's also query the actual records to see what's there
	var installations []struct {
		ID               string  `gorm:"column:id"`
		CodeName         *string `gorm:"column:code_name"`
		CreatedAt        string  `gorm:"column:created_at"`
		Status           string  `gorm:"column:status"`
		InstallationType string  `gorm:"column:installation_type"`
	}

	err = s.db.Table("customer_installations").
		Select("id, code_name, created_at, status, installation_type").
		Where("customer_installations.customer_id = ?", customerID).
		Find(&installations).Error

	if err != nil {
		log.Printf("Error querying installations: %v", err)
	} else {
		log.Printf("🔍 Found %d installation records for customer %s:", len(installations), customerID)
		for i, inst := range installations {
			codeName := "NULL"
			if inst.CodeName != nil {
				codeName = *inst.CodeName
			}
			log.Printf("🔍   [%d] ID: %s, CodeName: %s, Status: %s, Type: %s, CreatedAt: %s",
				i+1, inst.ID, codeName, inst.Status, inst.InstallationType, inst.CreatedAt)
		}

		// Additional debug: Check if there are any soft-deleted or hidden records
		var softDeletedCount int64
		s.db.Table("customer_installations").
			Where("customer_installations.customer_id = ? AND deleted_at IS NOT NULL", customerID).
			Count(&softDeletedCount)
		if softDeletedCount > 0 {
			log.Printf("⚠️  WARNING: Found %d soft-deleted records for this customer!", softDeletedCount)
		}
	}

	// Return count + 1 for the new installation
	return int(count) + 1
}

// sanitizeComment removes dangerous characters from comments
func (s *MikrotikProvisioningService) sanitizeComment(comment string) string {
	// Remove quotes and backslashes
	comment = strings.ReplaceAll(comment, "\"", "")
	comment = strings.ReplaceAll(comment, "\\", "")
	comment = strings.ReplaceAll(comment, ";", "")
	return comment
}

// validateProvisioningRequest validates the provisioning request
func (s *MikrotikProvisioningService) validateProvisioningRequest(req ProvisioningRequest) error {
	if req.InstallationID == "" {
		return fmt.Errorf("installation_id is required")
	}
	if req.CustomerID == "" {
		return fmt.Errorf("customer_id is required")
	}
	if req.CustomerName == "" {
		return fmt.Errorf("customer_name is required")
	}
	if req.MACAddress == "" {
		return fmt.Errorf("mac_address is required")
	}

	// Validate MAC address format
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	if !macRegex.MatchString(req.MACAddress) {
		return fmt.Errorf("invalid MAC address format: %s", req.MACAddress)
	}

	if req.StartDate == "" {
		return fmt.Errorf("start_date is required")
	}
	if req.StartTime == "" {
		return fmt.Errorf("start_time is required")
	}

	return nil
}

// updateInstallationStatus updates the installation record with provisioning results
func (s *MikrotikProvisioningService) updateInstallationStatus(installationID, codeName, ipAddress string) error {
	updates := map[string]interface{}{
		"provisioning_status":       "provisioned",
		"provisioning_completed_at": time.Now(),
		"code_name":                 codeName,
	}

	if ipAddress != "" {
		updates["ip_address"] = ipAddress
	}

	return s.db.Model(&entities.CustomerInstallation{}).
		Where("id = ?", installationID).
		Updates(updates).Error
}

// GetProvisioningLogs retrieves provisioning logs for an installation
func (s *MikrotikProvisioningService) GetProvisioningLogs(installationID string) ([]entities.InstallationProvisioningLog, error) {
	var logs []entities.InstallationProvisioningLog
	err := s.db.Where("customer_installation_id = ?", installationID).
		Order("createdAt DESC").
		Find(&logs).Error
	return logs, err
}

// RetryProvisioning retries failed provisioning
func (s *MikrotikProvisioningService) RetryProvisioning(logID string) (*ProvisioningResult, error) {
	var log entities.InstallationProvisioningLog
	if err := s.db.First(&log, "id = ?", logID).Error; err != nil {
		return nil, fmt.Errorf("log not found: %w", err)
	}

	// Increment retry count
	log.RetryCount++
	s.db.Save(&log)

	// Get installation details and retry
	var installation entities.CustomerInstallation
	if err := s.db.Preload("Customer").First(&installation, "id = ?", log.CustomerInstallationID).Error; err != nil {
		return nil, fmt.Errorf("installation not found: %w", err)
	}

	// Use installation completed date as start date for billing cycle
	startDate := time.Now().Format("2006-01-02") // Default to today
	startTime := "00:00:00"                      // Default to midnight

	if installation.InstallationCompletedAt != nil {
		startDate = installation.InstallationCompletedAt.Format("2006-01-02")
		startTime = installation.InstallationCompletedAt.Format("15:04:05")
	}

	req := ProvisioningRequest{
		InstallationID: installation.ID,
		CustomerID:     *installation.CustomerID,
		CustomerName:   installation.Customer.Name,
		MACAddress:     *log.MACAddress,
		StartDate:      startDate,
		StartTime:      startTime,
		MaxLimit:       "10M/10M", // TODO: Get from product/package
		Comment:        "Retry provisioning",
		DryRun:         false,
		CreatedBy:      *log.CreatedBy,
	}

	return s.ProvisionInstallation(req)
}

// getSalesRepresentativeInfo gets sales representative code for a customer
func (s *MikrotikProvisioningService) getSalesRepresentativeInfo(customerID string) (string, string) {
	var customer struct {
		SalesRepresentativeID *string `gorm:"column:sales_representative_id"`
		SalesRepCode          *string `gorm:"column:sales_rep_code"`
	}

	// Query customer with sales representative information
	err := s.db.Table("customer").
		Select("customer.sales_representative_id, users.code as sales_rep_code").
		Joins("LEFT JOIN users ON customer.sales_representative_id = users.id").
		Where("customer.id = ?", customerID).
		First(&customer).Error

	if err != nil {
		log.Printf("Error fetching sales representative info for customer %s: %v", customerID, err)
		return "UNKNOWN", ""
	}

	// Return sales rep code
	code := "UNKNOWN"

	if customer.SalesRepCode != nil && *customer.SalesRepCode != "" {
		code = *customer.SalesRepCode
	}

	return code, code
}

// generateEACode generates the next EA code (EA0001, EA0002, etc.)
func (s *MikrotikProvisioningService) generateEACode() string {
	// Get the highest EA number from existing installation codes
	var maxEANumber int

	// Query all existing code_names that contain EA followed by 4 digits
	var codes []string
	err := s.db.Table("customer_installations").
		Select("code_name").
		Where("code_name REGEXP '_EA[0-9]{4}_' AND code_name IS NOT NULL").
		Pluck("code_name", &codes).Error

	if err != nil {
		log.Printf("Error fetching existing EA codes: %v", err)
		return "EA0001" // Default to first EA code
	}

	// Extract EA numbers from existing codes
	for _, code := range codes {
		// Find EA pattern in the code
		eaPattern := regexp.MustCompile(`_EA(\d{4})_`)
		matches := eaPattern.FindStringSubmatch(code)
		if len(matches) > 1 {
			if number, err := strconv.Atoi(matches[1]); err == nil {
				if number > maxEANumber {
					maxEANumber = number
				}
			}
		}
	}

	// Generate next EA code
	nextEANumber := maxEANumber + 1
	return fmt.Sprintf("EA%04d", nextEANumber)
}
