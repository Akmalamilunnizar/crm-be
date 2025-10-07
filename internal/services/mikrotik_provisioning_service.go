package services

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
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
	startTime := time.Now()

	// Validate request
	if err := s.validateProvisioningRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate unique CODE-NAME
	codeName := s.generateCodeName(req)

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

	// Step 3: Check and create/update queue
	if err := s.createOrUpdateQueue(codeName, ipAddress, req.MaxLimit, req.Comment, result, req.DryRun); err != nil {
		return fmt.Errorf("failed to create/update queue: %w", err)
	}

	// Step 4: Check and create/update hotspot IP binding
	if err := s.createOrUpdateHotspotBinding(codeName, req.MACAddress, ipAddress, result, req.DryRun); err != nil {
		return fmt.Errorf("failed to create/update hotspot binding: %w", err)
	}

	// Step 5: Check and create/update netwatch
	if err := s.createOrUpdateNetwatch(codeName, req.MACAddress, ipAddress, result, req.DryRun); err != nil {
		return fmt.Errorf("failed to create/update netwatch: %w", err)
	}

	// Step 6: Check and create/update scheduler
	if err := s.createOrUpdateScheduler(codeName, req.MACAddress, req.StartDate, req.StartTime, result, req.DryRun); err != nil {
		return fmt.Errorf("failed to create/update scheduler: %w", err)
	}

	// Step 7: Check and create/update script
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
		cmd := fmt.Sprintf("/ip/hotspot/ip-binding/set [find mac-address=%s] address=%s to-addr=%s comment=\"%s\"",
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
	addCmd := fmt.Sprintf("/ip/hotspot/ip-binding/add mac-address=%s address=%s to-addr=%s comment=\"%s\"",
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
	// Simplify the onEvent string to avoid syntax errors
	onEvent := fmt.Sprintf("ip/hotspot/ip-binding/set type=reg [find mac-address=%s]", macAddress)

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
// Format: [area->name] - [customer->name] - R[number]
func (s *MikrotikProvisioningService) generateCodeName(req ProvisioningRequest) string {
	// Get the installation count for this customer in this area
	installationCount := s.getInstallationCount(req.CustomerID, req.AreaName)

	// Sanitize area name (remove special chars, keep only alphanumeric and spaces)
	areaReg := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	sanitizedArea := areaReg.ReplaceAllString(req.AreaName, "")
	sanitizedArea = strings.TrimSpace(sanitizedArea)

	// Sanitize customer name (remove special chars, keep only alphanumeric and spaces)
	customerReg := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	sanitizedCustomer := customerReg.ReplaceAllString(req.CustomerName, "")
	sanitizedCustomer = strings.TrimSpace(sanitizedCustomer)

	// Format: [area] - [customer] - R[number]
	return fmt.Sprintf("%s - %s - R%d", sanitizedArea, sanitizedCustomer, installationCount)
}

// getInstallationCount returns the number of installations for a customer in a specific area
func (s *MikrotikProvisioningService) getInstallationCount(customerID, areaName string) int {
	var count int64

	// Count installations for this customer in this area
	// Join through customer table to get area_id, then join with areas table
	err := s.db.Table("customer_installations").
		Joins("JOIN customers ON customer_installations.customer_id = customers.id").
		Joins("JOIN areas ON customers.area_id = areas.id").
		Where("customer_installations.customer_id = ? AND (areas.name_city = ? OR areas.name_subdistrict = ? OR areas.name_village = ?)",
			customerID, areaName, areaName, areaName).
		Count(&count).Error

	if err != nil {
		log.Printf("Error counting installations: %v", err)
		// Return 1 as default if there's an error
		return 1
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

	req := ProvisioningRequest{
		InstallationID: installation.ID,
		CustomerID:     *installation.CustomerID,
		CustomerName:   installation.Customer.Name,
		MACAddress:     *log.MACAddress,
		StartDate:      installation.PSBDate.Format("2006-01-02"),
		StartTime:      *installation.PSBTime,
		MaxLimit:       "10M/10M", // TODO: Get from product/package
		Comment:        "Retry provisioning",
		DryRun:         false,
		CreatedBy:      *log.CreatedBy,
	}

	return s.ProvisionInstallation(req)
}
