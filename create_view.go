package main

import (
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

	// Select database
	err = db.Exec("USE iqgncnzy_skripsi").Error
	if err != nil {
		log.Fatal("Failed to select database:", err)
	}

	// Drop the view if exists
	err = db.Exec("DROP VIEW IF EXISTS installation_report_complete").Error
	if err != nil {
		log.Fatal("Failed to drop view:", err)
	}

	// Create the view
	sql := `
	CREATE VIEW installation_report_complete AS
	SELECT
		ci.id AS installation_id,
		ci.customer_id,
		c.name AS customer_name,
		c.address AS customer_address,
		c.phone AS customer_phone,
		c.service_request_date AS tgl_permintaan_psb,
		ci.technician_id,
		u.name AS technician_name,
		u.phone AS technician_phone,
		ci.status AS installation_status,
		ci.installation_type,
		ci.notes AS installation_notes,
		ci.on_air_date,
		ci.trial_end_date,
		ci.service_ready_date,
		ci.installation_completed_at,
		CASE
			WHEN c.service_request_date IS NOT NULL AND ci.installation_completed_at IS NOT NULL
			THEN DATEDIFF(ci.installation_completed_at, c.service_request_date)
			ELSE NULL
		END AS durasi_psb,
		CASE
			WHEN c.service_request_date IS NOT NULL AND ci.installation_completed_at IS NOT NULL
			THEN CASE
				WHEN DATEDIFF(ci.installation_completed_at, c.service_request_date) <= 3
				THEN 'Tepat Waktu'
				ELSE 'Terlambat'
			END
			ELSE NULL
		END AS status_psb,
		ci.document_type,
		ci.document_photo,
		'' AS network_device_id,
		'' AS switch_id,
		'' AS port_number,
		'' AS remote_port,
		'' AS eth_port,
		'' AS mac_address,
		'' AS ip_static,
		'' AS kepemilikan_perangkat,
		'' AS router_brand,
		'' AS router_type,
		'' AS router_model,
		'' AS router_serial,
		'' AS product_id,
		'' AS product_name,
		'' AS product_description,
		0 AS product_price,
		NULL AS download_speed_mbps,
		NULL AS upload_speed_mbps,
		cs.id AS customer_service_id,
		cs.user_login,
		cs.password,
		cs.user_status,
		cs.installation_notes AS service_notes,
		cs.cable_type,
		cs.cable_length,
		cs.end_port_type,
		ci.createdAt AS installation_created_at,
		ci.updatedAt AS installation_updated_at
	FROM customer_installations ci
	LEFT JOIN customer c ON ci.customer_id = c.id AND c.deleted_at IS NULL
	LEFT JOIN users u ON ci.technician_id = u.id
	LEFT JOIN customer_services cs ON ci.id = cs.customer_installation_id;
	`

	err = db.Exec(sql).Error
	if err != nil {
		log.Fatal("Failed to create view:", err)
	}

	log.Println("View created successfully")
}
