# Broadcast API Implementation

This document describes the Golang backend implementation for the broadcast feature.

## Files Created

1. **`database/broadcast_history.sql`** - SQL migration script for the broadcast_history table
2. **`internal/models/entities/broadcast_model.go`** - GORM model for BroadcastHistory
3. **`internal/api/broadcast/handlers.go`** - API handlers for broadcast endpoints
4. **`internal/api/broadcast/routes.go`** - Route setup function
5. **`examples/broadcast_main_example.go`** - Example main function showing how to integrate

## Setup Instructions

### 1. Run the SQL Migration

First, run the SQL migration in HeidiSQL or your MySQL client:

```sql
-- Run: database/broadcast_history.sql
```

This will create the `broadcast_history` table with all necessary columns and indexes.

### 2. Update Your GORM Models

The `broadcast_model.go` file includes the GORM model. Make sure it's in your project structure:

```
crm-be/
  internal/
    models/
      entities/
        broadcast_model.go  ← New file
```

### 3. Add Broadcast Routes

In your main application file, add the broadcast routes:

```go
import (
    "crm-be/internal/api/broadcast"
    // ... other imports
)

func main() {
    // ... database setup ...
    
    app := fiber.New()
    
    // Setup broadcast routes
    broadcast.SetupBroadcastRoutes(app, db)
    
    // ... other routes ...
    
    app.Listen(":8080")
}
```

### 4. Authentication Middleware

Make sure your authentication middleware sets `userID` in the context:

```go
app.Use(func(c *fiber.Ctx) error {
    // Extract user ID from JWT token or session
    userID := extractUserIDFromToken(c)
    c.Locals("userID", userID)
    return c.Next()
})
```

## API Endpoints

### POST /api/broadcast/send

Sends a broadcast message to customers or team.

**Request Body:**
```json
{
  "target": "customers",  // or "team"
  "phones": ["081234567890", "081234567891"],  // Required only if target is "customers"
  "message": "Your broadcast message here",
  "template_key": "template_outage"  // Optional
}
```

**Response:**
```json
{
  "status": "broadcast logged",
  "recipients": 2
}
```

**Logic:**
- If `target == "team"`: Queries all users from the `users` table with non-empty phone numbers
- If `target == "customers"`: Uses the `phones` array from the request
- Saves broadcast history to database
- Returns number of recipients

### GET /api/broadcast/history

Retrieves broadcast history (last 50 records).

**Response:**
```json
[
  {
    "message": "Broadcast message text",
    "target_group": "customers",
    "status": "Sent",
    "sent_at": "2024-01-15T10:30:00Z",
    "created_by": "user-id-here",
    "user_name": "John Doe"
  }
]
```

**Logic:**
- Queries `broadcast_history` table
- Preloads User relation to get sender name
- Orders by `sent_at DESC`
- Limits to 50 records
- Capitalizes status for frontend compatibility

## Database Schema

### broadcast_history Table

```sql
CREATE TABLE `broadcast_history` (
  `id` VARCHAR(191) PRIMARY KEY,
  `message` TEXT,
  `target_group` VARCHAR(50),
  `recipient_count` INT DEFAULT 0,
  `status` ENUM('sent', 'failed', 'pending') DEFAULT 'pending',
  `template_key` VARCHAR(50),
  `sent_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `created_by` VARCHAR(191),
  `created_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3),
  FOREIGN KEY (`created_by`) REFERENCES `users` (`id`)
);
```

## GORM Models

### BroadcastHistory Model

```go
type BroadcastHistory struct {
    ID             string
    Message        string
    TargetGroup    string
    RecipientCount int
    Status         string  // "sent", "failed", "pending"
    TemplateKey    string
    SentAt         time.Time
    CreatedBy      string
    CreatedAt      time.Time
    UpdatedAt      time.Time
    
    // Relations
    User User `gorm:"foreignKey:CreatedBy"`
}
```

## Integration with WhatsApp Gateway

The handler includes a TODO comment where you should integrate your WhatsApp gateway:

```go
// TODO: Iterate over phoneList and send to WhatsApp Gateway
// Example:
// for _, phone := range phoneList {
//     err := whatsappGateway.SendMessage(phone, request.Message)
//     if err != nil {
//         log.Printf("Failed to send to %s: %v", phone, err)
//         // Optionally update broadcast history status to 'failed'
//     }
// }
```

## Error Handling

The handlers return appropriate HTTP status codes:

- `400 Bad Request` - Invalid request body or missing required fields
- `500 Internal Server Error` - Database errors or query failures

## Notes

1. The status field uses lowercase values ('sent', 'failed', 'pending') in the database, but the GET endpoint converts them to capitalized format ('Sent', 'Failed') for frontend compatibility.

2. The `created_by` field is optional - if no user ID is available from middleware, it will be set to empty string.

3. The broadcast history is saved immediately after determining recipients, before actually sending messages. You should update the status if sending fails.

4. For team broadcasts, all users with non-empty phone numbers are included automatically.

