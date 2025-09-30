# Test installation report creation with document photo
$apiUrl = "http://localhost:8080/api/admin/customer-installation/report/complete"

# You'll need to replace these with actual IDs from your database
$customerId = "your-customer-id-here"
$technicianId = "your-technician-id-here"
$documentPhotoPath = "uploads/documents/test/document_test.jpg"

Write-Host "Testing installation report creation..."
Write-Host "API URL: $apiUrl"

$requestData = @{
    customer_id = $customerId
    technician_id = $technicianId
    status = "completed"
    notes = "Test installation report with document photo"
    document_type = "KTP"
    document_photo = $documentPhotoPath
    installation_type = "new_installation"
    on_air_date = "2025-01-15"
    trial_end_date = "2025-02-15"
    service_ready_date = "2025-01-15"
    installation_completed_at = "2025-01-15T10:30:00"
    network_devices = @()
    customer_services = @()
    cables = @()
    image_ids = @()
} | ConvertTo-Json -Depth 3

Write-Host "Request data:"
Write-Host $requestData

try {
    $headers = @{
        "Content-Type" = "application/json"
        "Authorization" = "Bearer your-token-here"
    }
    
    $response = Invoke-RestMethod -Uri $apiUrl -Method Post -Body $requestData -Headers $headers
    Write-Host "Installation report created successfully!"
    Write-Host "Response:"
    $response | ConvertTo-Json -Depth 3
    
} catch {
    Write-Host "Installation report creation failed: $($_.Exception.Message)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response body: $responseBody"
    }
}

