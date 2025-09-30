package main

import (
	"fmt"
	"log"
	"os"
	"skripsi-be/internal/services"
	"strconv"
)

func main() {
	// Load environment variables
	host := os.Getenv("MIKROTIK_HOST")
	port := os.Getenv("MIKROTIK_PORT")
	username := os.Getenv("MIKROTIK_USERNAME")
	password := os.Getenv("MIKROTIK_PASSWORD")

	fmt.Printf("Testing MikroTik connection...\n")
	fmt.Printf("Host: %s\n", host)
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s\n", password)

	if host == "" || port == "" || username == "" || password == "" {
		log.Fatal("Missing MikroTik environment variables")
	}

	// Parse port
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	// Create MikroTik service
	config := &services.MikroTikConfig{
		Host:     host,
		Port:     portInt,
		Username: username,
		Password: password,
	}

	mtService := services.NewMikroTikService(config)

	// Test connection
	fmt.Printf("\nConnecting to MikroTik...\n")
	err = mtService.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer mtService.Disconnect()

	fmt.Printf("✅ Connected successfully!\n")

	// Test Netwatch devices
	fmt.Printf("\nGetting Netwatch devices...\n")
	devices, err := mtService.GetNetwatchDevices()
	if err != nil {
		log.Fatalf("Failed to get Netwatch devices: %v", err)
	}

	fmt.Printf("✅ Found %d Netwatch devices:\n", len(devices))
	for i, device := range devices {
		fmt.Printf("  %d. %+v\n", i+1, device)
	}

	// Test specific IP if provided
	if len(os.Args) > 1 {
		testIP := os.Args[1]
		fmt.Printf("\nTesting specific IP: %s\n", testIP)
		
		for _, device := range devices {
			if host, ok := device["host"].(string); ok && host == testIP {
				status, _ := device["status"].(string)
				fmt.Printf("✅ Found device %s with status: %s\n", testIP, status)
				return
			}
		}
		fmt.Printf("❌ Device %s not found in Netwatch\n", testIP)
	}

	fmt.Printf("\n✅ All tests passed!\n")
}
