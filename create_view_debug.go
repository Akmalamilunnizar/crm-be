package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "iqgncnzy_skripsi:XhYJOWlwNgsk@tcp(103.63.24.139:3306)/iqgncnzy_skripsi?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to database successfully")

	// Drop the view if exists
	log.Println("Dropping existing view...")
	err = db.Exec("DROP VIEW IF EXISTS installation_report_complete").Error
	if err != nil {
		log.Printf("Warning: Failed to drop view: %v", err)
	} else {
		log.Println("View dropped successfully")
	}

	// Create a simplified view first to test
	log.Println("Creating simplified view for testing...")
	simpleSql := `
	CREATE VIEW installation_report_complete AS
	SELECT
		ci.id AS installation_id,
		ci.customer_id AS customer_id,
		c.name AS customer_name
	FROM customer_installations ci
	LEFT JOIN customer c ON ci.customer_id = c.id
	WHERE ci.deleted_at IS NULL
	`

	err = db.Exec(simpleSql).Error
	if err != nil {
		log.Printf("❌ Failed to create simplified view: %v", err)
		log.Println("\nTrying to get more details about the error...")
		
		// Try to get SQL error details
		var result interface{}
		err2 := db.Raw("SELECT 1").Scan(&result).Error
		if err2 != nil {
			log.Printf("Database connection issue: %v", err2)
		}
		
		log.Fatal("Cannot proceed")
	}

	log.Println("✅ Simplified view created successfully!")
	
	// Now try to drop and create the full view
	log.Println("\nDropping simplified view...")
	db.Exec("DROP VIEW IF EXISTS installation_report_complete")
	
	log.Println("Creating full view...")
	fullSql := `
	CREATE VIEW installation_report_complete AS
	SELECT
		ci.id AS installation_id,
		ci.customer_id AS customer_id,
		c.name AS customer_name,
		c.address AS customer_address,
		c.phone AS customer_phone,
		c.service_request_date AS tgl_permintaan_psb,
		ci.technician_id AS technician_id,
		u.name AS technician_name,
		u.phone AS technician_phone,
		ci.status AS installation_status,
		ci.installation_type AS installation_type,
		ci.notes AS installation_notes,
		ci.on_air_date AS on_air_date,
		ci.trial_end_date AS trial_end_date,
		ci.service_ready_date AS service_ready_date,
		ci.installation_completed_at AS installation_completed_at,
		CASE
			WHEN c.service_request_date IS NOT NULL AND ci.installation_completed_at IS NOT NULL
			THEN TO_DAYS(ci.installation_completed_at) - TO_DAYS(c.service_request_date)
			ELSE NULL
		END AS durasi_psb,
		CASE
			WHEN c.service_request_date IS NOT NULL AND ci.installation_completed_at IS NOT NULL
			THEN CASE
				WHEN (TO_DAYS(ci.installation_completed_at) - TO_DAYS(c.service_request_date)) <= 3
				THEN 'Tepat Waktu'
				ELSE 'Terlambat'
			END
			ELSE NULL
		END AS status_psb,
		ci.document_type AS document_type,
		ci.document_photo AS document_photo,
		nd.id AS network_device_id,
		nd.switch_id AS switch_id,
		nd.port_number AS port_number,
		nd.remote_port AS remote_port,
		nd.eth_port AS eth_port,
		nd.mac_address AS mac_address,
		nd.ip_static AS ip_static,
		nd.kepemilikan_perangkat AS kepemilikan_perangkat,
		a.brand AS router_brand,
		a.type AS router_type,
		a.model AS router_model,
		a.serial_number AS router_serial,
		p.name AS product_name,
		p.description AS product_description,
		p.price AS product_price,
		p.download_speed_mbps AS download_speed_mbps,
		p.upload_speed_mbps AS upload_speed_mbps,
		cs.id AS customer_service_id,
		cs.user_login AS user_login,
		cs.password AS password,
		cs.user_status AS user_status,
		cs.installation_notes AS service_notes,
		cs.cable_type AS cable_type,
		cs.cable_length AS cable_length,
		cs.end_port_type AS end_port_type,
		ci.createdAt AS installation_created_at,
		ci.updatedAt AS installation_updated_at,
		p.id AS product_id
	FROM customer_installations ci
	LEFT JOIN customer c ON ci.customer_id = c.id
	LEFT JOIN users u ON ci.technician_id = u.id
	LEFT JOIN network_devices nd ON ci.id = nd.customer_installation_id
	LEFT JOIN assets a ON nd.assets_id = a.id
	LEFT JOIN products p ON nd.product_id = p.id
	LEFT JOIN customer_services cs ON ci.id = cs.customer_installation_id
	WHERE ci.deleted_at IS NULL
	`

	err = db.Exec(fullSql).Error
	if err != nil {
		log.Printf("❌ Failed to create full view: %v", err)
		log.Fatal("View creation failed")
	}

	log.Println("✅ Full view created successfully!")
	fmt.Println("\n========================================")
	fmt.Println("SUCCESS! The installation_report_complete view has been updated.")
	fmt.Println("You can now test the installation report detail page.")
	fmt.Println("========================================")
}
