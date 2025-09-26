package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
		User  struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
			Role  struct {
				Name string `json:"name"`
			} `json:"role"`
		} `json:"user"`
	} `json:"data"`
}

func main() {
	// Test different credentials
	testCredentials := []LoginRequest{
		{Email: "admin@email.com", Password: "password"},
		{Email: "admin@email.com", Password: "admin"},
		{Email: "admin@email.com", Password: "123456"},
		{Email: "admin", Password: "password"},
		{Email: "admin", Password: "admin"},
		{Email: "test@test.com", Password: "test"},
	}

	for i, creds := range testCredentials {
		fmt.Printf("\n=== Test %d: %s ===\n", i+1, creds.Email)
		
		// Create JSON payload
		jsonData, _ := json.Marshal(creds)
		
		// Make request
		resp, err := http.Post("http://localhost:3001/api/auth/login", 
			"application/json", 
			bytes.NewBuffer(jsonData))
		
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}
		
		// Read response
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		
		fmt.Printf("Status: %d\n", resp.StatusCode)
		fmt.Printf("Response: %s\n", string(body))
		
		// Try to parse response
		var loginResp LoginResponse
		if err := json.Unmarshal(body, &loginResp); err == nil {
			if loginResp.Success {
				fmt.Printf("✅ SUCCESS! Token: %s\n", loginResp.Data.Token)
				fmt.Printf("User: %s (%s)\n", loginResp.Data.User.Name, loginResp.Data.User.Email)
				fmt.Printf("Role: %s\n", loginResp.Data.User.Role.Name)
				break
			} else {
				fmt.Printf("❌ Failed: %s\n", loginResp.Message)
			}
		} else {
			fmt.Printf("❌ Could not parse response\n")
		}
	}
}
