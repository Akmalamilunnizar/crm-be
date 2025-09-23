package main

import (
	"log"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/models/entities"
	"time"
)

func main() {
	db := database.GetDB()

	// Check if technician steps already exist
	var count int64
	db.Model(&entities.TechnicianStep{}).Count(&count)
	if count > 0 {
		log.Println("Technician steps already exist, skipping seed")
		return
	}

	// Predefined technician steps based on Panduan Teknis
	steps := []entities.TechnicianStep{
		{
			StepOrder:   1,
			Title:       "Pengecekan Power Listrik Utama",
			Description: "Memeriksa kondisi power listrik utama untuk memastikan suplai listrik stabil",
			Tools:       "Multimeter, Test Pen, Kabel Tester",
			SpareParts:  "Stop Kontak, Kabel Listrik, Fuse",
			Procedure:   "1. Periksa tegangan listrik dengan multimeter\n2. Cek kondisi stop kontak\n3. Test koneksi kabel listrik\n4. Verifikasi fuse tidak putus",
			Solution:    "Ganti komponen yang rusak, perbaiki koneksi yang longgar",
			IsActive:    true,
			CreatedAt:   &[]time.Time{time.Now()}[0],
			UpdatedAt:   &[]time.Time{time.Now()}[0],
		},
		{
			StepOrder:   2,
			Title:       "Pengecekan Adaptor Perangkat Media Converter",
			Description: "Memeriksa kondisi adaptor power untuk media converter",
			Tools:       "Multimeter, Test Pen",
			SpareParts:  "Adaptor, Kabel Power",
			Procedure:   "1. Periksa output voltage adaptor\n2. Cek kondisi kabel adaptor\n3. Test koneksi ke media converter",
			Solution:    "Ganti adaptor jika output tidak sesuai, perbaiki koneksi kabel",
			IsActive:    true,
			CreatedAt:   &[]time.Time{time.Now()}[0],
			UpdatedAt:   &[]time.Time{time.Now()}[0],
		},
		{
			StepOrder:   3,
			Title:       "Pengecekan Perangkat Media Converter (CVT)",
			Description: "Memeriksa kondisi dan fungsi media converter",
			Tools:       "Multimeter, Test Pen, Laptop",
			SpareParts:  "CVT, Fastcon, PatchCord",
			Procedure:   "1. Periksa LED indikator CVT\n2. Test koneksi input/output\n3. Cek konfigurasi CVT\n4. Verifikasi link status",
			Solution:    "Reset CVT, ganti jika rusak, perbaiki konfigurasi",
			IsActive:    true,
			CreatedAt:   &[]time.Time{time.Now()}[0],
			UpdatedAt:   &[]time.Time{time.Now()}[0],
		},
		{
			StepOrder:   4,
			Title:       "Pengecekan Status Link Upstream (Kedua Sisi Ujung)",
			Description: "Memeriksa status koneksi upstream dari kedua sisi",
			Tools:       "Laptop, Ping Test, Traceroute",
			SpareParts:  "PatchCord, Splicer",
			Procedure:   "1. Ping test ke gateway\n2. Traceroute ke server\n3. Cek latency dan packet loss\n4. Verifikasi routing",
			Solution:    "Perbaiki routing, ganti kabel jika rusak, konfigurasi ulang",
			IsActive:    true,
			CreatedAt:   &[]time.Time{time.Now()}[0],
			UpdatedAt:   &[]time.Time{time.Now()}[0],
		},
		{
			StepOrder:   5,
			Title:       "Pencarian & Perbaikan Titik Kesalahan/Putus/Kerusakan Kabel (Optik)",
			Description: "Mencari dan memperbaiki kerusakan pada kabel optik",
			Tools:       "OTDR, Power Meter, Splicer",
			SpareParts:  "Kabel Optik, Splicer, Fastcon",
			Procedure:   "1. Gunakan OTDR untuk deteksi kerusakan\n2. Cek power loss pada kabel\n3. Identifikasi titik putus\n4. Lakukan splicing jika diperlukan",
			Solution:    "Splice kabel yang putus, ganti kabel jika rusak parah",
			IsActive:    true,
			CreatedAt:   &[]time.Time{time.Now()}[0],
			UpdatedAt:   &[]time.Time{time.Now()}[0],
		},
		{
			StepOrder:   6,
			Title:       "Pengecekan Fungsi Kabel UTP",
			Description: "Memeriksa kondisi dan fungsi kabel UTP",
			Tools:       "Cable Tester, Multimeter",
			SpareParts:  "UTP, RJ45",
			Procedure:   "1. Test continuity kabel UTP\n2. Cek koneksi RJ45\n3. Verifikasi pinout kabel\n4. Test kecepatan kabel",
			Solution:    "Ganti kabel UTP jika rusak, crimping ulang RJ45",
			IsActive:    true,
			CreatedAt:   &[]time.Time{time.Now()}[0],
			UpdatedAt:   &[]time.Time{time.Now()}[0],
		},
		{
			StepOrder:   7,
			Title:       "Pengecekan Adaptor Perangkat Router Wireless",
			Description: "Memeriksa kondisi adaptor power untuk router wireless",
			Tools:       "Multimeter, Test Pen",
			SpareParts:  "Adaptor, Kabel Power",
			Procedure:   "1. Periksa output voltage adaptor\n2. Cek kondisi kabel adaptor\n3. Test koneksi ke router",
			Solution:    "Ganti adaptor jika output tidak sesuai",
			IsActive:    true,
			CreatedAt:   &[]time.Time{time.Now()}[0],
			UpdatedAt:   &[]time.Time{time.Now()}[0],
		},
		{
			StepOrder:   8,
			Title:       "Pengecekan Perangkat Router",
			Description: "Memeriksa kondisi dan fungsi router",
			Tools:       "Laptop, Multimeter, Test Pen",
			SpareParts:  "Router, Adaptor",
			Procedure:   "1. Periksa LED indikator router\n2. Test koneksi LAN/WAN\n3. Cek konfigurasi router\n4. Verifikasi DHCP server",
			Solution:    "Reset router, ganti jika rusak, update firmware",
			IsActive:    true,
			CreatedAt:   &[]time.Time{time.Now()}[0],
			UpdatedAt:   &[]time.Time{time.Now()}[0],
		},
		{
			StepOrder:   9,
			Title:       "Pengecekan Konfigurasi Router Wireless",
			Description: "Memeriksa dan memperbaiki konfigurasi wireless router",
			Tools:       "Laptop, Smartphone",
			SpareParts:  "Router",
			Procedure:   "1. Akses web interface router\n2. Cek konfigurasi SSID dan password\n3. Verifikasi channel wireless\n4. Test koneksi wireless",
			Solution:    "Update konfigurasi wireless, reset ke default jika perlu",
			IsActive:    true,
			CreatedAt:   &[]time.Time{time.Now()}[0],
			UpdatedAt:   &[]time.Time{time.Now()}[0],
		},
		{
			StepOrder:   10,
			Title:       "Pengecekan End Device (HP, Laptop, TV, CCTV)",
			Description: "Memeriksa kondisi dan koneksi perangkat end user",
			Tools:       "Laptop, Smartphone, Test Device",
			SpareParts:  "Kabel UTP, RJ45, Adaptor",
			Procedure:   "1. Test koneksi internet di device\n2. Cek konfigurasi IP address\n3. Verifikasi DNS settings\n4. Test aplikasi yang bermasalah",
			Solution:    "Update driver, reset network settings, ganti kabel jika perlu",
			IsActive:    true,
			CreatedAt:   &[]time.Time{time.Now()}[0],
			UpdatedAt:   &[]time.Time{time.Now()}[0],
		},
	}

	// Insert technician steps
	for _, step := range steps {
		if err := db.Create(&step).Error; err != nil {
			log.Printf("Error creating technician step %d: %v", step.StepOrder, err)
		} else {
			log.Printf("Created technician step: %s", step.Title)
		}
	}

	// Seed spare parts
	// Seed Spare Parts
	// Seed Spare Parts
	// Seed Spare Parts
	spareParts := []entities.SparePart{
		{Name: "Stop Kontak", Description: stringPtr("Power outlet untuk peralatan"), Category: "Electrical"},
		{Name: "Kabel Listrik", Description: stringPtr("Kabel power untuk peralatan"), Category: "Electrical"},
		{Name: "CVT", Description: stringPtr("Media Converter"), Category: "Network Device"},
		{Name: "Fastcon", Description: stringPtr("Konektor fiber optik"), Category: "Fiber Optic"},
		{Name: "PatchCord", Description: stringPtr("Kabel patch untuk koneksi"), Category: "Fiber Optic"},
		{Name: "Splicer", Description: stringPtr("Alat untuk menyambung kabel optik"), Category: "Fiber Optic"},
		{Name: "UTP", Description: stringPtr("Kabel UTP untuk jaringan"), Category: "Network Cable"},
		{Name: "RJ45", Description: stringPtr("Konektor kabel UTP"), Category: "Network Connector"},
		{Name: "Adaptor", Description: stringPtr("Power adapter untuk peralatan"), Category: "Electrical"},
		{Name: "Router", Description: stringPtr("Perangkat Router Wireless"), Category: "Network Device"},
	}

	for _, part := range spareParts {
		if err := db.Create(&part).Error; err != nil {
			log.Printf("Error creating spare part %s: %v", part.Name, err)
		} else {
			log.Printf("Created spare part: %s", part.Name)
		}
	}

	log.Println("Technician steps and spare parts seeded successfully!")
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
