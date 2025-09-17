package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type MikroTikService struct {
	sshClient   *ssh.Client
	config      *MikroTikConfig
	isConnected bool
}

type MikroTikConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MikroTikLog struct {
	Host      string    `json:"host"`
	Comment   string    `json:"comment"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Category  string    `json:"category"`
	Raw       string    `json:"raw,omitempty"`
}

type ConnectionStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
}

func NewMikroTikService(config *MikroTikConfig) *MikroTikService {
	return &MikroTikService{
		config: config,
	}
}

func (s *MikroTikService) Connect() error {
	sshConfig := &ssh.ClientConfig{
		User: s.config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         20 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port), sshConfig)
	if err != nil {
		s.isConnected = false
		return fmt.Errorf("failed to connect to MikroTik: %v", err)
	}

	s.sshClient = client
	s.isConnected = true
	log.Printf("Connected to MikroTik at %s:%d", s.config.Host, s.config.Port)
	return nil
}

func (s *MikroTikService) Disconnect() error {
	if s.sshClient != nil {
		err := s.sshClient.Close()
		s.sshClient = nil
		s.isConnected = false
		log.Println("Disconnected from MikroTik")
		return err
	}
	return nil
}

func (s *MikroTikService) IsConnected() bool {
	return s.isConnected && s.sshClient != nil
}

func (s *MikroTikService) GetConnectionStatus() ConnectionStatus {
	if s.IsConnected() {
		return ConnectionStatus{
			Status:  "connected",
			Message: "Connected to MikroTik",
			Host:    s.config.Host,
			Port:    s.config.Port,
		}
	}
	return ConnectionStatus{
		Status:  "disconnected",
		Message: "Not connected",
		Host:    s.config.Host,
		Port:    s.config.Port,
	}
}

func (s *MikroTikService) GetLogs(timeRange string) ([]MikroTikLog, error) {
	if !s.IsConnected() {
		return nil, fmt.Errorf("not connected to MikroTik")
	}

	// Convert timeRange to MikroTik time format
	var timeFilter string
	switch timeRange {
	case "1h":
		timeFilter = "1h"
	case "6h":
		timeFilter = "6h"
	case "1d":
		timeFilter = "1d"
	case "7d":
		timeFilter = "7d"
	case "30d":
		timeFilter = "30d"
	default:
		timeFilter = "1d"
	}

	// Get Netwatch logs specifically - focus on UP/DOWN status changes
	command := fmt.Sprintf("/log print where topics~\"netwatch\" and time>%s", timeFilter)

	session, err := s.sshClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer session.Close()

	output, err := session.Output(command)
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %v", err)
	}

	return s.parseLogs(string(output)), nil
}

func (s *MikroTikService) parseLogs(logData string) []MikroTikLog {
	lines := strings.Split(logData, "\n")
	var logs []MikroTikLog

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		log := s.parseLogLine(line)
		if log != nil {
			logs = append(logs, *log)
		}
	}

	// Reverse to show newest first
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}

	return logs
}

func (s *MikroTikService) parseLogLine(line string) *MikroTikLog {
	// Parse Netwatch logs - they typically contain UP/DOWN status
	// Example format: "netwatch host 192.168.1.100 DOWN"
	// or more detailed format with timestamp and message

	// Check if this is a Netwatch UP/DOWN log
	if strings.Contains(strings.ToLower(line), "netwatch") {
		// Extract host/IP address
		var host string
		var status string
		var message string = line

		// Try to extract IP address pattern
		if strings.Contains(line, ".") {
			// Look for IP address in the line
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.Count(part, ".") == 3 {
					host = part
					break
				}
			}
		}

		// Determine status from the log content
		lineLower := strings.ToLower(line)
		if strings.Contains(lineLower, "down") || strings.Contains(lineLower, "timeout") || strings.Contains(lineLower, "unreachable") {
			status = "DOWN"
		} else if strings.Contains(lineLower, "up") || strings.Contains(lineLower, "reachable") || strings.Contains(lineLower, "alive") {
			status = "UP"
		} else {
			// Default based on log level
			if strings.Contains(lineLower, "error") || strings.Contains(lineLower, "critical") {
				status = "DOWN"
			} else {
				status = "UP"
			}
		}

		if host == "" {
			host = "unknown"
		}

		return &MikroTikLog{
			Host:      host,
			Comment:   message,
			Status:    status,
			Timestamp: time.Now(), // In real implementation, parse actual timestamp
			Type:      s.getLogType(status),
			Category:  "netwatch",
			Raw:       line,
		}
	}

	// Handle pipe-separated format (custom format)
	if strings.Contains(line, "|") {
		parts := strings.Split(line, "|")
		if len(parts) >= 3 {
			host := strings.TrimSpace(parts[0])
			comment := strings.TrimSpace(parts[1])
			status := strings.TrimSpace(parts[2])

			return &MikroTikLog{
				Host:      host,
				Comment:   comment,
				Status:    strings.ToUpper(status),
				Timestamp: time.Now(),
				Type:      s.getLogType(status),
				Category:  "netwatch",
				Raw:       line,
			}
		}
	}

	// Fallback for other log formats - but prioritize UP/DOWN detection
	timestamp := time.Now() // In real implementation, parse actual timestamp

	// Extract status based on keywords
	lineLower := strings.ToLower(line)
	status := "INFO"
	if strings.Contains(lineLower, "down") || strings.Contains(lineLower, "error") || strings.Contains(lineLower, "critical") {
		status = "DOWN"
	} else if strings.Contains(lineLower, "up") {
		status = "UP"
	} else if strings.Contains(lineLower, "warning") {
		status = "WARNING"
	}

	// Extract host/topic if available
	topic := "system"
	if strings.Contains(line, "topic=") {
		if topicStart := strings.Index(line, "topic="); topicStart != -1 {
			topicEnd := strings.Index(line[topicStart:], " ")
			if topicEnd == -1 {
				topicEnd = len(line) - topicStart
			}
			topic = strings.TrimSpace(line[topicStart+6 : topicStart+topicEnd])
		}
	}

	// Get the message part
	message := line
	if strings.Contains(line, `message="`) {
		if msgStart := strings.Index(line, `message="`); msgStart != -1 {
			msgEnd := strings.Index(line[msgStart+9:], `"`)
			if msgEnd != -1 {
				message = line[msgStart+9 : msgStart+9+msgEnd]
			}
		}
	}

	return &MikroTikLog{
		Host:      topic,
		Comment:   message,
		Status:    status,
		Timestamp: timestamp,
		Type:      s.getLogType(status),
		Category:  "system",
		Raw:       line,
	}
}

func (s *MikroTikService) getLogType(status string) string {
	statusLower := strings.ToLower(status)
	switch statusLower {
	case "up", "info", "reachable", "alive":
		return "success"
	case "down", "error", "critical", "timeout", "unreachable":
		return "error"
	case "warning":
		return "warning"
	default:
		return "info"
	}
}

func (s *MikroTikService) ExecuteCommand(command string) (string, error) {
	if !s.IsConnected() {
		return "", fmt.Errorf("not connected to MikroTik")
	}

	session, err := s.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer session.Close()

	output, err := session.Output(command)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}

	return string(output), nil
}

func (s *MikroTikService) GetSystemInfo() (map[string]interface{}, error) {
	if !s.IsConnected() {
		return nil, fmt.Errorf("not connected to MikroTik")
	}

	info := make(map[string]interface{})

	// Get system resource info using correct MikroTik commands
	commands := map[string]string{
		"resource": "/system resource print",
		"identity": "/system identity print",
		"clock":    "/system clock print",
		"version":  "/system package print where name=\"system\"",
	}

	for key, cmd := range commands {
		output, err := s.ExecuteCommand(cmd)
		if err != nil {
			log.Printf("Failed to execute %s: %v", cmd, err)
			continue
		}
		info[key] = strings.TrimSpace(output)
	}

	// Parse resource info to extract specific values
	if resourceOutput, ok := info["resource"].(string); ok {
		// Parse CPU load, memory usage, uptime from resource output
		lines := strings.Split(resourceOutput, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "cpu-load:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					info["cpu_load"] = strings.TrimSpace(parts[1])
				}
			} else if strings.Contains(line, "free-memory:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					info["memory"] = strings.TrimSpace(parts[1])
				}
			} else if strings.Contains(line, "uptime:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					info["uptime"] = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return info, nil
}

// GetNetwatchDevices gets all Netwatch devices from MikroTik
func (s *MikroTikService) GetNetwatchDevices() ([]map[string]interface{}, error) {
	if !s.IsConnected() {
		return nil, fmt.Errorf("not connected to MikroTik")
	}

	command := "/tool netwatch print"
	output, err := s.ExecuteCommand(command)
	if err != nil {
		return nil, fmt.Errorf("failed to get netwatch devices: %v", err)
	}

	// Debug logging
	log.Printf("Netwatch devices raw output: %s", output)

	devices := s.parseNetwatchDevices(output)
	log.Printf("Parsed netwatch devices: %+v", devices)

	return devices, nil
}

// parseNetwatchDevices parses the netwatch devices output
func (s *MikroTikService) parseNetwatchDevices(output string) []map[string]interface{} {
	var devices []map[string]interface{}
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		device := make(map[string]interface{})

		// Parse netwatch device information
		// Example format: "0 host=192.168.1.1 interval=30s timeout=5s status=up"
		parts := strings.Fields(line)
		for _, part := range parts {
			if strings.Contains(part, "=") {
				kv := strings.SplitN(part, "=", 2)
				if len(kv) == 2 {
					key := strings.TrimSpace(kv[0])
					value := strings.TrimSpace(kv[1])

					// Convert status to uppercase for consistency
					if key == "status" {
						value = strings.ToUpper(value)
					}

					device[key] = value
				}
			}
		}

		// Only add if we have essential fields
		if _, hasHost := device["host"]; hasHost {
			devices = append(devices, device)
		}
	}

	return devices
}

// GetNetwatchLogs gets real-time netwatch logs
func (s *MikroTikService) GetNetwatchLogs(timeRange string) ([]MikroTikLog, error) {
	if !s.IsConnected() {
		return nil, fmt.Errorf("not connected to MikroTik")
	}

	// Convert timeRange to MikroTik time format
	var timeFilter string
	switch timeRange {
	case "1h":
		timeFilter = "1h"
	case "6h":
		timeFilter = "6h"
	case "1d":
		timeFilter = "1d"
	case "7d":
		timeFilter = "7d"
	case "30d":
		timeFilter = "30d"
	default:
		timeFilter = "1d"
	}

	// Get Netwatch logs specifically - focus on UP/DOWN status changes
	command := fmt.Sprintf("/log print where topics~\"netwatch\" and time>%s", timeFilter)

	session, err := s.sshClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer session.Close()

	output, err := session.Output(command)
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %v", err)
	}

	return s.parseNetwatchLogs(string(output)), nil
}

// parseNetwatchLogs parses netwatch-specific logs
func (s *MikroTikService) parseNetwatchLogs(logData string) []MikroTikLog {
	lines := strings.Split(logData, "\n")
	var logs []MikroTikLog

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		log := s.parseNetwatchLogLine(line)
		if log != nil {
			logs = append(logs, *log)
		}
	}

	// Reverse to show newest first
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}

	return logs
}

// parseNetwatchLogLine parses a single netwatch log line
func (s *MikroTikService) parseNetwatchLogLine(line string) *MikroTikLog {
	// Parse Netwatch logs - they typically contain UP/DOWN status
	// Example format: "netwatch host 192.168.1.100 DOWN"
	// or more detailed format with timestamp and message

	// Check if this is a Netwatch UP/DOWN log
	if strings.Contains(strings.ToLower(line), "netwatch") {
		// Extract host/IP address
		var host string
		var status string
		var message string = line

		// Try to extract IP address pattern
		if strings.Contains(line, ".") {
			// Look for IP address in the line
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.Count(part, ".") == 3 {
					host = part
					break
				}
			}
		}

		// Determine status from the log content
		lineLower := strings.ToLower(line)
		if strings.Contains(lineLower, "down") || strings.Contains(lineLower, "timeout") || strings.Contains(lineLower, "unreachable") {
			status = "DOWN"
		} else if strings.Contains(lineLower, "up") || strings.Contains(lineLower, "reachable") || strings.Contains(lineLower, "alive") {
			status = "UP"
		} else {
			// Default based on log level
			if strings.Contains(lineLower, "error") || strings.Contains(lineLower, "critical") {
				status = "DOWN"
			} else {
				status = "UP"
			}
		}

		if host == "" {
			host = "unknown"
		}

		return &MikroTikLog{
			Host:      host,
			Comment:   message,
			Status:    status,
			Timestamp: time.Now(), // In real implementation, parse actual timestamp
			Type:      s.getLogType(status),
			Category:  "netwatch",
			Raw:       line,
		}
	}

	return nil
}
