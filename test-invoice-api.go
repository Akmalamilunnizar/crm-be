package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Test invoice API
	fmt.Println("Testing Invoice API...")

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Test GET /api/admin/invoice/:id
	resp, err := http.Get("http://localhost:3001/api/admin/invoice/1")
	if err != nil {
		log.Printf("Error calling API: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("Error decoding response: %v", err)
			return
		}

		fmt.Printf("Response: %+v\n", result)

		if data, ok := result["data"].(map[string]interface{}); ok {
			if amount, ok := data["amount"].(float64); ok {
				fmt.Printf("Invoice Amount: %v\n", amount)
			}
			if items, ok := data["invoice_items"].([]interface{}); ok {
				fmt.Printf("Invoice Items Count: %d\n", len(items))
				for i, item := range items {
					if itemMap, ok := item.(map[string]interface{}); ok {
						fmt.Printf("Item %d: %+v\n", i+1, itemMap)
					}
				}
			}
		}
	} else {
		fmt.Printf("API returned status: %d\n", resp.StatusCode)
	}
}
