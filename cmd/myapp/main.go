package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	inv "skripsi-be/internal/api/admin/invoice"
	"skripsi-be/internal/api/admin/recurring_invoice"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/routes"
	"skripsi-be/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	database.GetDB()
	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		BodyLimit:    50 * 1024 * 1024,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	})
	// CORS configuration
	app.Use(cors.New(cors.Config{
    	AllowOrigins: "http://localhost:3000, http://127.0.0.1:3000, http://192.168.1.7:3000, http://192.168.1.11:3000, http://localhost:3002, http://127.0.0.1:3002, http://192.168.1.11:3002",
    	AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    	AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
    	AllowCredentials: true,
    }))

	// Recovery and request logging for diagnostics
	app.Use(recovermw.New())
	app.Use(logger.New())

	// Initialize routes
	routes.RouteFiber(app)

	// Auto-connect to MikroTik on startup if credentials are available
	go func() {
		// Wait a bit for the server to fully start
		time.Sleep(2 * time.Second)

		host := os.Getenv("MIKROTIK_HOST")
		port := os.Getenv("MIKROTIK_PORT")
		username := os.Getenv("MIKROTIK_USERNAME")
		password := os.Getenv("MIKROTIK_PASSWORD")

		if host != "" && port != "" && username != "" && password != "" {
			log.Printf("[mikrotik] auto-connecting to %s:%s", host, port)

			// Create MikroTik service directly
			config := &services.MikroTikConfig{
				Host:     host,
				Port:     parsePort(port),
				Username: username,
				Password: password,
			}

			mtService := services.NewMikroTikService(config)
			if err := mtService.Connect(); err != nil {
				log.Printf("[mikrotik] auto-connect failed: %v", err)
			} else {
				services.SetSharedMikroTikService(mtService)
				log.Printf("[mikrotik] auto-connected successfully")
			}
		} else {
			log.Printf("[mikrotik] no credentials found in environment variables")
		}
	}()

	// Background scheduler to process due recurring invoices periodically
	go func() {
		// Track processed invoices to avoid re-enforcement
		processedInvoices := make(map[string]bool)

		run := func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[scheduler] panic: %v", r)
				}
			}()
			start := time.Now()
			log.Printf("[recurring] run startawdasdwadwawsda: %s", start.Format(time.RFC3339))
			db := database.GetDB()

			// First: Enforcement tick and precheck so we always see logs even if recurring blocks
			// Enforcement tick marker
			log.Printf("[enforce] tick")
			// Precheck: get newest invoice per customer and enforce based on status
			invRepo := inv.NewAdminInvoiceRepository(db)
			mt := services.GetSharedMikroTikService()
			mtConnected := mt != nil && mt.IsConnected()
			log.Printf("[enforce] precheck: mt_connected=%t", mtConnected)

			if mtConnected && mt != nil {
				// Get newest invoice for each customer
				newestInvoices, err := invRepo.FindNewestInvoicePerCustomer()
				if err != nil {
					log.Printf("[enforce] precheck: fetch newest invoices error: %v", err)
				} else {
					regularChanged := 0
					bypassedChanged := 0
					seen := map[string]bool{}

					for _, iv := range newestInvoices {
						// Skip if already processed this invoice
						if processedInvoices[iv.ID] {
							continue
						}

						var devices []struct{ MacAddress *string }
						if err := db.Table("network_devices").
							Select("mac_address").
							Where("customer_id = ? AND mac_address IS NOT NULL AND mac_address <> ''", iv.CustomerID).
							Scan(&devices).Error; err != nil {
							log.Printf("[enforce] devices fetch error customer=%s: %v", iv.CustomerID, err)
							continue
						}

						deviceCount := 0
						bindingType := ""

						// Determine binding type based on invoice status and due date
						if iv.Status == "unpaid" && iv.DueDate != nil && iv.DueDate.Before(time.Now()) {
							bindingType = "regular" // Restrict due unpaid
						} else if iv.Status == "paid" {
							bindingType = "bypassed" // Unlock paid
						} else if iv.Status == "unpaid" && (iv.DueDate == nil || iv.DueDate.After(time.Now())) {
							bindingType = "bypassed" // Allow unpaid but not due yet
						} else {
							continue // Skip if not unpaid or paid
						}

						for _, d := range devices {
							if d.MacAddress == nil {
								continue
							}
							mac := *d.MacAddress
							if seen[mac] {
								continue
							}
							if err := mt.SetHotspotIPBindingType(mac, bindingType); err != nil {
								log.Printf("[enforce] set %s failed mac=%s: %v", bindingType, mac, err)
							} else {
								deviceCount++
								seen[mac] = true
								if bindingType == "regular" {
									regularChanged++
								} else {
									bypassedChanged++
								}
							}
						}

						// Mark this invoice as processed if we found devices for it
						if deviceCount > 0 {
							processedInvoices[iv.ID] = true
							log.Printf("[enforce] processed invoice %s (%d devices) -> %s", iv.ID, deviceCount, bindingType)
						}
					}

					if regularChanged > 0 || bypassedChanged > 0 {
						log.Printf("[enforce] set type=regular for %d device(s), type=bypassed for %d device(s)", regularChanged, bypassedChanged)
					} else {
						log.Printf("[enforce] no devices updated (newest invoices: %d, already processed: %d)", len(newestInvoices), len(processedInvoices))
					}
				}
			} else {
				log.Printf("[enforce] MikroTik not connected; skipping")
			}

			// Then: process recurring invoices
			repo := recurring_invoice.NewAdminRecurringInvoiceRepository(db)
			n, err := repo.ProcessDueRecurringInvoices()
			dur := time.Since(start)
			if err != nil {
				log.Printf("[recurring] run error after %s: %v", dur, err)
			} else if n > 0 {
				log.Printf("[recurring] generated %d invoices in %s", n, dur)
			} else {
				log.Printf("[recurring] nothing due (duration %s)", dur)
			}
		}

		// immediate run for faster feedback
		run()

		ticker := time.NewTicker(1 * time.Minute) // adjust cadence as needed
		log.Printf("[recurring] scheduler started, interval=1m")
		defer ticker.Stop()
		for range ticker.C {
			run()
		}
	}()

	log.Fatal(app.Listen(":3001"))
}

func parsePort(portStr string) int {
	if port, err := strconv.Atoi(portStr); err == nil {
		return port
	}
	return 22 // default SSH port
}
