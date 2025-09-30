package customerinstallation

import (
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"skripsi-be/internal/helpers"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ReportInstallationController struct {
	service ReportInstallationServiceInterface
}

func NewReportInstallationController(service ReportInstallationServiceInterface) *ReportInstallationController {
	return &ReportInstallationController{service: service}
}

// CreateReportInstallation - Create installation report with multipart form data
func (c *ReportInstallationController) CreateReportInstallation(ctx *fiber.Ctx) error {
	var request CreateReportInstallationRequest

	// Parse multipart form
	_, err := ctx.MultipartForm()
	if err != nil {
		return helpers.ResponseUtils(ctx, 400, false, "Invalid multipart form data", err.Error())
	}

	// Parse form values
	request.CustomerID = ctx.FormValue("customer_id")
	request.TechnicianID = ctx.FormValue("technician_id")
	request.Status = ctx.FormValue("status")
	request.Notes = ctx.FormValue("notes")
	request.InstallationType = ctx.FormValue("installation_type")
	request.OnAirDate = ctx.FormValue("on_air_date")
	request.TrialEndDate = ctx.FormValue("trial_end_date")
	request.ServiceReadyDate = ctx.FormValue("service_ready_date")
	request.InstallationCompletedAt = ctx.FormValue("installation_completed_at")
	request.DocumentType = ctx.FormValue("document_type")
	request.SwitchID = ctx.FormValue("switch_id")
	request.PortNumber = ctx.FormValue("port_number")
	request.RemotePort = ctx.FormValue("remote_port")
	request.EthPort = ctx.FormValue("eth_port")
	request.MacAddress = ctx.FormValue("mac_address")
	request.IPStatic = ctx.FormValue("ip_static")
	request.StatusPerangkat = ctx.FormValue("status_perangkat")
	request.KepemilikanPerangkat = ctx.FormValue("kepemilikan_perangkat")
	request.LastPingStatus = ctx.FormValue("last_ping_status")
	request.AssetsID = ctx.FormValue("assets_id")
	request.CableType = ctx.FormValue("cable_type")
	// Parse cable_length from string to float64
	if cableLengthStr := ctx.FormValue("cable_length"); cableLengthStr != "" {
		if parsed, err := strconv.ParseFloat(cableLengthStr, 64); err == nil {
			request.CableLength = parsed
		}
	}
	request.EndPortType = ctx.FormValue("end_port_type")
	request.UserLogin = ctx.FormValue("user_login")
	request.Password = ctx.FormValue("password")
	request.UserStatus = ctx.FormValue("user_status")
	request.InstallationNotes = ctx.FormValue("installation_notes")

	// Handle document photo upload
	var documentPhotoPath string

	// Log all form files for debugging
	form, formErr := ctx.MultipartForm()
	if formErr == nil {
		log.Printf("Multipart form files found: %v", len(form.File))
		for fieldName, files := range form.File {
			log.Printf("Form field '%s' has %d files", fieldName, len(files))
			for i, file := range files {
				log.Printf("  File %d: %s (size: %d)", i, file.Filename, file.Size)
			}
		}
	}

	if file, err := ctx.FormFile("document_photo"); err == nil {
		// Log file upload info
		log.Printf("Document photo upload started - filename: %s, size: %d", file.Filename, file.Size)

		// Validate file type
		if !isValidImageFile(file) {
			return helpers.ResponseUtils(ctx, 400, false, "Invalid file type. Only JPG and PNG are allowed", nil)
		}

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		filename := "document_" + time.Now().Format("20060102_150405") + ext

		// Create upload directory if not exists
		uploadDir := "uploads/installations/documents"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return helpers.ResponseUtils(ctx, 500, false, "Failed to create upload directory", err.Error())
		}

		// Save file
		filePath := filepath.Join(uploadDir, filename)
		if err := ctx.SaveFile(file, filePath); err != nil {
			return helpers.ResponseUtils(ctx, 500, false, "Failed to save document photo", err.Error())
		}

		documentPhotoPath = filePath

		// Log successful upload
		log.Printf("Document photo uploaded successfully - filePath: %s, filename: %s", filePath, filename)
	} else {
		// Log if no file was uploaded
		log.Printf("No document photo uploaded - error: %s", err.Error())

		// Check if document_photo field exists in form values
		if docPhotoValue := ctx.FormValue("document_photo"); docPhotoValue != "" {
			log.Printf("document_photo form value (not file): %s", docPhotoValue)
		} else {
			log.Printf("No document_photo field found in form values")
		}
	}

	// Set document photo path
	request.DocumentPhoto = documentPhotoPath

	// Log request data for debugging
	log.Printf("Installation report request data - customer_id: %s, technician_id: %s, document_type: %s, document_photo: %s, installation_type: %s, assets_id: %s, switch_id: %s, port_number: %s, mac_address: %s, ip_static: %s, cable_type: %s, cable_length: %f, user_login: %s, user_status: %s, installation_notes: %s",
		request.CustomerID, request.TechnicianID, request.DocumentType, request.DocumentPhoto, request.InstallationType, request.AssetsID, request.SwitchID, request.PortNumber, request.MacAddress, request.IPStatic, request.CableType, request.CableLength, request.UserLogin, request.UserStatus, request.InstallationNotes)

	// Validate required fields
	if request.CustomerID == "" {
		return helpers.ResponseUtils(ctx, 400, false, "Customer ID is required", nil)
	}
	if request.TechnicianID == "" {
		return helpers.ResponseUtils(ctx, 400, false, "Technician ID is required", nil)
	}
	if request.AssetsID == "" {
		return helpers.ResponseUtils(ctx, 400, false, "Assets ID is required", nil)
	}

	// Validate IP format if provided
	if request.IPStatic != "" && !isValidIP(request.IPStatic) {
		return helpers.ResponseUtils(ctx, 400, false, "Invalid IP address format", nil)
	}

	// Validate MAC address format if provided
	if request.MacAddress != "" && !isValidMAC(request.MacAddress) {
		return helpers.ResponseUtils(ctx, 400, false, "Invalid MAC address format", nil)
	}

	// Create installation report
	installation, err := c.service.CreateReportInstallationService(request)
	if err != nil {
		return helpers.ResponseUtils(ctx, 500, false, "Failed to create installation report", err.Error())
	}

	return helpers.ResponseUtils(ctx, 201, true, "Installation report created successfully", installation)
}

// isValidImageFile - Validate if uploaded file is a valid image
func isValidImageFile(file *multipart.FileHeader) bool {
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png"}
	contentType := file.Header.Get("Content-Type")

	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return true
		}
	}
	return false
}

// isValidIP - Basic IP address validation
func isValidIP(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}
		// Additional validation can be added here
	}
	return true
}

// isValidMAC - Basic MAC address validation
func isValidMAC(mac string) bool {
	// Remove common separators
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "-", "")
	mac = strings.ReplaceAll(mac, ".", "")

	// Check if it's 12 characters long and contains only hex characters
	if len(mac) != 12 {
		return false
	}

	for _, char := range mac {
		if !((char >= '0' && char <= '9') || (char >= 'A' && char <= 'F') || (char >= 'a' && char <= 'f')) {
			return false
		}
	}
	return true
}
